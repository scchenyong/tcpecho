package main

import (
	"net"
	"time"
	"sync/atomic"
	"bufio"
	"os"
	"strconv"
	"fmt"
)

var (
	cc = int64(0)
	t  = time.Second * 5
	l  = 1024
)

func main() {
	args := os.Args
	al := len(args)
	if al < 2 {
		panic("无服务端地址参数!")
	}
	addr := args[1]
	c := 1
	if al >= 3 {
		x, err := strconv.Atoi(args[2])
		if err == nil {
			c = x
		} else {
			println("并发参数无效，默认：1")
		}
	}
	if al >= 4 {
		x, err := strconv.Atoi(args[3])
		if err == nil {
			l = x
		} else {
			println("并发参数无效，默认：1")
		}
	}
	go func() {
		t1 := time.NewTicker(t)
		for _ = range t1.C {
			println("cc:", atomic.LoadInt64(&cc))
		}
	}()

	s := []byte(fmt.Sprintf(fmt.Sprintf("%s0%dd\n", "%", l-2), 1))
	for {
		for i := int64(0); i < int64(c)-atomic.LoadInt64(&cc); i++ {
			go newConnect(addr, s)
		}
		time.Sleep(1 * time.Minute)
	}
	select {}
}

func newConnect(addr string, msg []byte) {
	conn, err := net.Dial("tcp4", addr)
	if err != nil {
		return
	}
	atomic.AddInt64(&cc, 1)
	go func(c1 net.Conn) {
		defer func(c2 net.Conn) {
			c2.Close()
			atomic.AddInt64(&cc, -1)
			c2 = nil
		}(c1)
		buf := bufio.NewReader(c1)
		for {
			c1.SetReadDeadline(time.Now().Add(60 * time.Second))
			_, err := buf.ReadBytes('\n')
			if err != nil {
				return
			}
		}
	}(conn)
	go func(c1 net.Conn) {
		for {
			c1.SetWriteDeadline(time.Now().Add(60 * time.Second))
			c1.Write(msg)
			time.Sleep(30 * time.Second)
		}
	}(conn)
	conn.Write(msg)
}
