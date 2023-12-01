package srpc

import (
	"gii/glog"
	"gii/srpc/codec"
	"go/ast"
	"io"
	"log"
	"net"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
)

// 服务方法
type methodType struct {
	method    reflect.Method
	argType   reflect.Type
	replyType reflect.Type
	numCalls  uint64
}

// NumCalls 服务方法调用次数
func (t *methodType) NumCalls() uint64 {
	return atomic.LoadUint64(&t.numCalls)
}

// NewArg 创建第一个参数实例
func (t *methodType) NewArg() reflect.Value {
	if t.argType.Kind() == reflect.Ptr {
		return reflect.New(t.argType.Elem())
	}
	return reflect.New(t.argType).Elem()
}

// NewReply 创建第二个参数实例
func (t *methodType) NewReply() reflect.Value {
	reply := reflect.New(t.replyType.Elem())
	switch t.replyType.Elem().Kind() {
	case reflect.Map:
		reply.Elem().Set(reflect.MakeMap(t.replyType.Elem()))
	case reflect.Slice:
		reply.Elem().Set(reflect.MakeSlice(t.replyType.Elem(), 0, 0))
	}
	return reply
}

// 服务
type service struct {
	name    string
	typ     reflect.Type
	rcvr    reflect.Value
	methods map[string]*methodType
}

// 判断参数是否被导出
func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}

// 创建新的服务实例
func newService(rcvr any) *service {
	s := new(service)
	s.rcvr = reflect.ValueOf(rcvr)
	s.name = reflect.Indirect(s.rcvr).Type().Name()
	if !ast.IsExported(s.name) {
		glog.ErrorF("rpc server: method %s is not exported", s.name)
	}
	s.typ = reflect.TypeOf(rcvr)
	s.methods = make(map[string]*methodType)
	s.registerMethod()
	return s
}

// 注册服务对应的方法
func (s *service) registerMethod() {
	numMethod := s.typ.NumMethod()
	for i := 0; i < numMethod; i++ {
		method := s.typ.Method(i)
		mType := method.Type
		// 参数
		// 入参 0->self 1->arg 2->*reply
		// 出参 0->error
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		// 校验出参是不是error
		if mType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
			continue
		}
		argType, replyType := mType.In(1), mType.In(2)
		if !isExportedOrBuiltinType(argType) || !isExportedOrBuiltinType(replyType) {
			continue
		}
		s.methods[method.Name] = &methodType{
			method:    method,
			argType:   argType,
			replyType: replyType,
		}
		glog.InfoF("rpc server: register server %s.%s", s.name, method.Name)
	}
}

// 服务方法调用
func (s *service) call(mt *methodType, arg, reply reflect.Value) error {
	// update call num
	atomic.AddUint64(&mt.numCalls, 1)
	fn := mt.method.Func
	returnValue := fn.Call([]reflect.Value{s.rcvr, arg, reply})
	if errInterface := returnValue[0].Interface(); errInterface != nil {
		return errInterface.(error)
	}
	return nil
}

var invalidRequest = struct{}{}

type Request struct {
	header     *codec.Header
	arg, reply reflect.Value
	mt         *methodType
	svc        *service
}

// Server RPC Server
type Server struct {
	ServerMap sync.Map // 注册的服务表
}

func NewServer() *Server {
	return &Server{}
}

var DefaultServer = NewServer()

func Accept(lis net.Listener) { DefaultServer.Accept(lis) }

func Register(rcvr any) { DefaultServer.Register(rcvr) }

// Register 注册服务
func (s *Server) Register(rcvr any) {
	ss := newService(rcvr)
	if _, ok := s.ServerMap.LoadOrStore(ss.name, ss); ok {
		glog.ErrorF("rpc server: can not register service:%s", ss.name)
	}
}

// 查询服务
func (s *Server) findService(serviceMethod string) (svc *service, mt *methodType) {
	// serviceMethod = "T.method"
	idx := strings.LastIndex(serviceMethod, ".")
	if idx == 0 {
		glog.ErrorF("rpc server: cannot find service %s", serviceMethod)
	}
	serviceName, methodName := serviceMethod[:idx], serviceMethod[idx+1:]
	svcc, ok := s.ServerMap.Load(serviceName)
	if !ok {
		glog.ErrorF("rpc server: cannot get service by name %s", serviceName)
	}
	svc = svcc.(*service)
	mt = svc.methods[methodName]
	if mt == nil {
		glog.ErrorF("rpc server: cannot get methodType by name %s", methodName)
	}
	return
}

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
	defer func() { _ = conn.Close() }()
	// 解析protocol
	var p RpcProto
	_, err := conn.Read(p[0:])
	if err != nil {
		glog.Error("rpc server: decode option error: ", err)
	}

	fn := codec.TypeMaps[CheckEnc(p)]
	if fn == nil {
		glog.Error("rpc server: invalid codec")
	}
	s.ServerCodec(fn(conn))
}

func (s *Server) ServerCodec(c codec.Codec) {
	defer func() { _ = c.Close() }()
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
}

func (s *Server) readRequestHeader(c codec.Codec) (*codec.Header, error) {
	var h codec.Header
	if err := c.ReadHeader(&h); err != nil {
		if err != io.EOF && err != io.ErrUnexpectedEOF {
			glog.Error("rpc server: read header error: ", err)
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
	req.svc, req.mt = s.findService(h.ServerMethod)

	req.arg = req.mt.NewArg()
	req.reply = req.mt.NewReply()

	argInterface := req.arg.Interface()
	if req.arg.Type().Kind() != reflect.Ptr {
		argInterface = req.arg.Addr().Interface()
	}

	if err = c.ReadBody(argInterface); err != nil {
		log.Println("rpc server: read argv err:", err)
	}
	return req, nil
}

func (s *Server) sendResponse(c codec.Codec, header *codec.Header, body any, sending *sync.Mutex) {
	sending.Lock()
	defer sending.Unlock()
	if err := c.Write(header, body); err != nil {
		glog.Error("rpc server: write error: ", err)
	}
}

func (s *Server) handleRequest(c codec.Codec, req *Request, wg *sync.WaitGroup, sending *sync.Mutex) {
	defer wg.Done()
	glog.Info(req.header, req.arg)
	err := req.svc.call(req.mt, req.arg, req.reply)
	if err != nil {
		req.header.Error = err
		s.sendResponse(c, req.header, invalidRequest, sending)
	}
	s.sendResponse(c, req.header, req.reply.Interface(), sending)
}
