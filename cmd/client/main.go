package main

import (
	"bufio"
	"io"
	"log"
	"nat-penetration/define"
	"nat-penetration/helper"
)

func main() {
	tcpConn, err := helper.CreateConn(define.ControlServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("[连接成功] %s\n", tcpConn.RemoteAddr().String())
	for {
		s, err := bufio.NewReader(tcpConn).ReadString('\n')
		if err != nil {
			log.Printf("Get Data Error: %v\n", err)
			continue
		}
		// New Connection
		if s == define.NewConnection {
			go messageForward()
		}
	}
}
func messageForward() {
	// 连接服务端的隧道
	tunnelConn, err := helper.CreateConn(define.TunnelServerAddr)
	if err != nil {
		panic(err)
	}
	// 连接客户端的服务
	localConn, err := helper.CreateConn(define.LocalServerAddr)
	if err != nil {
		panic(err)
	}
	// 消息转发
	go io.Copy(localConn, tunnelConn)
	go io.Copy(tunnelConn, localConn)
}
