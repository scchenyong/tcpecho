package server

import (
	"testing"
	"net"
	"sync"
	"bufio"
)

func TestStart(t *testing.T) {
	//go Start()
	addr := "0.0.0.0:8888"
	//addr := "118.31.72.117:10008"
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		println("client>>连接服务出错", err.Error())
		return
	}
	println("client>>连接服务器成功：", conn.RemoteAddr().String())
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func(c net.Conn) {
		var b = make([]byte, 1024)
		n, err := c.Read(b)
		if err != nil {
			println("client>>读取消息出错", err.Error())
			return
		}
		bytes := b[:n]
		println("client>>接收消息：", string(bytes))
		wg.Done()
	}(conn)
	conn.Write([]byte("test\n"))
	wg.Wait()
}

func Test_ServerStart(t *testing.T) {
	go Start()
	select {}
}

func TestStartFor(t *testing.T) {
	x := 0
	for i := 0; i < 500000; i++ {
		//addr := "0.0.0.0:8888"
		addr := "118.31.72.117:8888"
		conn, err := net.Dial("tcp4", addr)
		if err != nil {
			t.Error("client>>连接服务出错", err.Error())
		}
		println("client>>连接服务器成功：", conn.RemoteAddr().String())
		x += 1
		conn.Write([]byte("test\n"))
		println("连接成功：", i+1)
		//time.Sleep(1 * time.Second)
	}
}

func TestStartForGo(t *testing.T) {
	x := 0
	for i := 0; i < 5000; i++ {
		//addr := "0.0.0.0:8888"
		go func() {
			addr := "118.31.72.117:8888"
			conn, err := net.Dial("tcp4", addr)
			if err != nil {
				t.Error("client>>连接服务出错", err.Error())
			}
			println("client>>连接服务器成功：", conn.RemoteAddr().String())
			x += 1
			conn.Write([]byte("test\n"))
			println("连接成功：", i+1)
		}()
		//time.Sleep(1 * time.Second)
	}
	select {}
}

func BenchmarkStart(b *testing.B) {
	x := 0
	for i := 0; i < b.N; i++ {
		addr := "118.31.72.117:10008"
		conn, err := net.Dial("tcp4", addr)
		if err != nil {
			b.Error("client>>连接服务出错", err.Error())
		}
		println("client>>连接服务器成功：", conn.RemoteAddr().String())
		x += 1
		conn.Write([]byte("test\n"))
	}
}

func BenchmarkStart1(b *testing.B) {
	//go Start()
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		conn, err := net.Dial("tcp4", "0.0.0.0:8888")
		if err != nil {
			println("client>>连接服务出错", err.Error())
			return
		}
		println("client>>连接服务器成功：", conn.RemoteAddr().String())
		wg.Add(1)
		go func(c net.Conn) {
			buf := bufio.NewReader(c)
			ine, err := buf.ReadString('\n')
			if err != nil {
				println("client>>读取消息出错", err.Error())
				return
			}
			println("client>>接收消息：", ine)
			wg.Done()
			conn.Close()
		}(conn)
		conn.Write([]byte("test\n"))
	}
	wg.Wait()
}
