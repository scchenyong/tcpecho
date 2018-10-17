package server

import (
	"net"
	"bufio"
	"io"
	"sync/atomic"
	"time"
)

var (
	cc = int64(0)
	t  = time.Second * 5
)

func Start() {
	go func() {
		t1 := time.NewTicker(t)
		for _ = range t1.C {
			println("cc:", atomic.LoadInt64(&cc))
		}
	}()
	l, e := net.Listen("tcp4", "0.0.0.0:20023")
	if e != nil {
		println("server>>开启TCP监控失败", e.Error())
		return
	}
	println("server>>TCPEcho服务启动成功!", l.Addr().String())
	for {
		c0, err := l.Accept()
		if err != nil {
			println("server>>接受连接失败", err.Error())
			break
		}
		atomic.AddInt64(&cc, 1)
		go func(c1 net.Conn) {
			defer func(c2 net.Conn) {
				c2.Close()
				atomic.AddInt64(&cc, -1)
				c2 = nil
			}(c1)
			buf := bufio.NewReaderSize(c1, 1024)
			for {
				c1.SetReadDeadline(time.Now().Add(60 * time.Second))
				ine, err := buf.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						return
					}
					break
				}
				c1.SetWriteDeadline(time.Now().Add(60 * time.Second))
				c1.Write(ine)
			}
		}(c0)
	}
}
