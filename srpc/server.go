package srpc

import (
	"fmt"
	"gii/glog"
	"gii/srpc/codec"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
)

var invalidRequest = struct{}{}

type Request struct {
	header       *codec.Header
	argv, replyv reflect.Value
}

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func (s *Server) Accept(lis net.Listener) {
	for {
		conn, err := lis.Accept()
		if err != nil {
			glog.Error("rpc server: accept err:", err)
		}
		go s.ServerConn(conn)
	}
}

func (s *Server) ServerConn(conn io.ReadWriteCloser) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Println(err)
			glog.Error("rpc: conn close error: ", err)
		}
	}()
	// 解析protocol
	var p [2]byte
	_, err := conn.Read(p[0:])
	if err != nil {
		glog.Error("rpc server: decode option error: ", err)
	}

	fn := codec.TypeMaps[CheckEnc(p[0:])]
	if fn == nil {
		glog.Error("rpc server: invalid codec")
	}
	s.ServerCodec(fn(conn))
}

func (s *Server) ServerCodec(c codec.Codec) {
	sending := new(sync.Mutex)
	wg := new(sync.WaitGroup)
	// 读取数据
	for {
		req, err := s.readRequest(c)
		if err != nil {
			if req == nil {
				break
			}
			req.header.Error = err
			s.sendResponse(c, req.header, invalidRequest, sending)
			continue
		}
		wg.Add(1)
		go s.handleRequest(c, req, wg, sending)
	}
	wg.Wait()
	err := c.Close()
	if err != nil {
		glog.Error(err)
	}
}

func (s *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := c.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			glog.Error("rpc server: read header error: ", err)
			//log.Println("rpc server: read header error:", err)
		}
		return nil, err
	}
	return &h, nil
}

func (s *Server) readRequest(c codec.Codec) (*Request, error) {
	h, err := s.readRequestHeader(c)
	if err != nil {
		return nil, err
	}
	req := &Request{header: h}
	req.argv = reflect.New(reflect.TypeOf(""))
	if err = c.ReadBody(req.argv.Interface()); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (s *Server) sendResponse(c codec.Codec, header *codec.Header, body interface{}, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := c.Write(header, body); err != nil {
		glog.Error("rpc server: write error: ", err)
	}
}

func (s *Server) handleRequest(c codec.Codec, req *Request, wg *sync.WaitGroup, sending *sync.Mutex) {
	defer wg.Done()
	log.Println(req.header, req.argv.Elem())
	req.replyv = reflect.ValueOf(fmt.Sprintf("geerpc resp %d", req.header.Seq))
	s.sendResponse(c, req.header, req.replyv.Interface(), sending)
}
