package srpc

import (
	"context"
	"errors"
	"fmt"
	"gii/glog"
	"gii/srpc/codec"
	"io"
	"net"
	"sync"
	"time"
)

// Call 一次请求RPC的信息
type Call struct {
	Seq          uint64 // 唯一请求编号
	ServerMethod string // 服务方法 "T.method"
	Args         any
	reply        any
	Error        error
	Done         chan *Call
}

// Client 客户端结构信息
type Client struct {
	codec    codec.Codec  // 编解码器
	seq      uint64       // 唯一请求编号，每次请求依次递增
	header   codec.Header // header信息
	proto    *RProto      // 协议
	closing  bool
	shutdown bool
	sending  sync.Mutex
	mu       sync.Mutex
	pending  map[uint64]*Call // 待处理的call请求
}

var ErrShutdown = errors.New("connection is shut down")
var _ io.Closer = (*Client)(nil)

// 表示这次RPC调用完成
func (c *Call) done() {
	c.Done <- c
}

func (c *Client) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing {
		return ErrShutdown
	}
	c.closing = true
	return c.codec.Close()
}

func (c *Client) IsAvailable() bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return !c.closing && !c.shutdown
}

// 将call添加到client
func (c *Client) registerCall(call *Call) (seq uint64, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closing || c.shutdown {
		return 0, ErrShutdown
	}
	call.Seq = c.seq
	c.pending[c.seq] = call
	c.seq++
	return call.Seq, nil
}

// 根据 seq，从 client.pending 中移除对应的 call
func (c *Client) removeCall(seq uint64) *Call {
	c.mu.Lock()
	defer c.mu.Unlock()
	call := c.pending[seq]
	delete(c.pending, seq)
	return call
}

// 服务端或客户端发生错误时调用，将 shutdown 设置为 true，且将错误信息通知所有 pending 状态的 call
func (c *Client) terminateCall(err error) {
	c.mu.Lock()
	c.sending.Lock()
	defer c.mu.Unlock()
	defer c.sending.Unlock()
	c.shutdown = true
	for _, call := range c.pending {
		call.Error = err
		call.done()
	}
}

// 接受响应
func (c *Client) receive() {
	var err error
	for err == nil {
		var h codec.Header
		if err = c.codec.ReadHeader(&h); err != nil {
			break
		}
		call := c.removeCall(h.Seq)
		switch {
		case call == nil:
			err = c.codec.ReadBody(nil)
		case h.Error != "":
			call.Error = fmt.Errorf(h.Error)
			err = c.codec.ReadBody(nil)
			call.done()
		default:
			if err = c.codec.ReadBody(call.reply); err != nil {
				call.Error = err
			}
			call.done()
		}
	}
	// err != nil
	c.terminateCall(err)
}

// 发送请求
func (c *Client) send(call *Call) {
	c.sending.Lock()
	defer c.sending.Unlock()
	// 注册call
	seq, err := c.registerCall(call)
	if err != nil {
		call.Error = err
		call.done()
		return
	}
	// 处理header
	c.header.ServerMethod = call.ServerMethod
	c.header.Seq = seq
	// 编码并发送
	if err = c.codec.Write(&c.header, call.Args); err != nil {
		call := c.removeCall(seq)
		if call != nil {
			call.Error = err
			call.done()
		}
	}
}

// CallTimeout example
// ctx, _ := context.WithTimeout(context.Background(), time.Second)
// var reply int
// err := client.CallTimeout(ctx, "Foo.Sum", &Args{1, 2}, &reply)
func (c *Client) CallTimeout(ctx context.Context, serviceMethod string, args, reply any) error {
	call := c.Go(serviceMethod, args, reply, make(chan *Call, 1))
	select {
	case <-ctx.Done():
		c.removeCall(call.Seq)
		return errors.New("rpc client: call failed: " + ctx.Err().Error())
	case call := <-call.Done:
		return call.Error
	}
}

// Call 同步阻塞调用
func (c *Client) Call(serviceMethod string, args, reply any) error {
	call := <-c.Go(serviceMethod, args, reply, make(chan *Call, 1)).Done
	return call.Error
}

// Go 异步调用
func (c *Client) Go(serviceMethod string, args, reply any, done chan *Call) *Call {
	call := new(Call)
	call.ServerMethod = serviceMethod
	call.Args = args
	call.reply = reply
	if done == nil {
		done = make(chan *Call, 10)
	} else {
		if cap(done) == 0 {
			glog.Error("rpc client: done channel is unbuffered")
		}
	}
	call.Done = done
	c.send(call)
	return call
}

// DefaultDial 默认gob编码,默认10s超时
func DefaultDial(network, addr string) (client *Client) {
	return Dial(network, addr, DefaultProtocol())
}

// Dial 创建client实例
// proto是协议类型
func Dial(network, addr string, proto *RProto) (client *Client) {
	// 判断是否需要设置连接超时
	var conn net.Conn
	var err error
	if proto.deProto.ConnectTimeout > 0 {
		conn, err = net.DialTimeout(network, addr, proto.deProto.ConnectTimeout)
	} else {
		conn, err = net.Dial(network, addr)
	}

	if err != nil {
		glog.Error("rpc client: Dial error: ", err)
	}

	ch := make(chan *Client)
	go func() {
		ch <- NewClient(conn, proto)
	}()

	if proto.deProto.ConnectTimeout == 0 {
		return <-ch
	}

	select {
	case <-time.After(proto.deProto.ConnectTimeout):
		glog.ErrorF("rpc client: connect timeout: expect within %s", proto.deProto.ConnectTimeout)
	case client = <-ch:
		return
	}
	defer func() {
		if client == nil {
			if err = conn.Close(); err != nil {
				glog.Error("rpc client: Dial close error: ", err)
			}
		}
	}()
	return
}

func NewClient(conn net.Conn, proto *RProto) *Client {
	fn := codec.TypeMaps[proto.deProto.EncType]
	if fn == nil {
		glog.Error("rpc server: invalid codec")
	}
	_, err := conn.Write(proto.proto[0:])
	if err != nil {
		glog.Error("rpc client: send protocol error: ", err)
	}
	return newClientCodec(fn(conn), proto)
}

func newClientCodec(c codec.Codec, proto *RProto) *Client {
	client := &Client{
		seq:     1,
		codec:   c,
		proto:   proto,
		pending: make(map[uint64]*Call),
	}
	// 创建协程，接受响应
	go client.receive()
	return client
}
