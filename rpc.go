package main

import (
	"fmt"
	"gii/glog"
	"gii/srpc"
	"log"
	"net"
	"sync"
	"time"
)

type HelloServer struct {
}

func (h *HelloServer) Hello(req string, reply *string) error {
	//time.Sleep(time.Second * 11)
	*reply = "hellosada:" + req
	return nil
}

func startServer(addr chan string) {
	srpc.Register(new(HelloServer))
	// pick a free port
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("network error:", err)
	}
	glog.Info("start rpc server on", l.Addr())
	addr <- "localhost:1234"
	srpc.Accept(l)
}

func main() {
	addr := make(chan string)
	go startServer(addr)

	client := srpc.DefaultDial("tcp", <-addr)
	defer func() { _ = client.Close() }()

	time.Sleep(time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := fmt.Sprintf("world %d", i)
			var reply string
			if err := client.Call("HelloServer.Hello", args, &reply); err != nil {
				glog.Error("call Foo.Sum error:", err)
			}
			glog.Info("reply:", reply)
		}(i)
	}
	wg.Wait()
}
