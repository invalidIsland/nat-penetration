package test

import (
	"bufio"
	"io"
	"log"
	"nat-penetration/helper"
	"net"
	"sync"
	"testing"
)

const (
	ControlServerAddr = "0.0.0.0:8080"
	RequestServerAddr = "0.0.0.0:8081"
	KeepAliveStr      = "KeepAlive\n"
)

var wg sync.WaitGroup
var clientConn *net.TCPConn

// 服务端
func TestUserServer(t *testing.T) {
	wg.Add(1)
	// 监听控制中心
	go ControlServer()
	// 监听用户的请求
	go RequestServer()
	wg.Wait()
}

func ControlServer() {
	listener, err := helper.CreateListen(ControlServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("ControlServer is listening on %s\n", ControlServerAddr)
	for {
		clientConn, err = listener.AcceptTCP()
		if err != nil {
			return
		}
		go helper.KeepAlive(clientConn)
	}
}

func RequestServer() {
	listener, err := helper.CreateListen(RequestServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("RequestServer is listening on %s\n", RequestServerAddr)
	for {
		tcpConn, err := listener.AcceptTCP()
		if err != nil {
			return
		}
		go io.Copy(clientConn, tcpConn)
		go io.Copy(tcpConn, clientConn)
	}
}

// 客户端
func TestUserClient(t *testing.T) {
	tcpConn, err := helper.CreateConn(ControlServerAddr)
	if err != nil {
		log.Printf("[连接失败] %s", err)
	}
	for {
		s, err := bufio.NewReader(tcpConn).ReadString('\n')
		if err != nil {
			log.Printf("Get Data Error: %v", err)
			continue
		}
		log.Printf("Get Data: %v", s)
	}
}

// 用户端
func TestUserRequestClient(t *testing.T) {
	tcpConn, err := helper.CreateConn(RequestServerAddr)
	if err != nil {
		log.Printf("[连接失败] %s", err)
	}
	_, err = tcpConn.Write([]byte("hello world\n"))
	if err != nil {
		log.Printf("[发送失败] %s", err)
	}
}
