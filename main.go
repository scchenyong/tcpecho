package main

import (
	"sctele.com/tcpecho/server"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	server.Start()
}
