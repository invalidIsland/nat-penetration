package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"nat-penetration/define"
	"nat-penetration/helper"
	"nat-penetration/service"
	"net"
	"sync"
)

var wg sync.WaitGroup
var controlConn *net.TCPConn
var userConn *net.TCPConn

func main() {
	wg.Add(1)
	// 控制中心监听
	go controlListen()
	// 用户请求监听
	go userRequestListen()
	// 隧道监听
	go tunnelListen()
	// 启动Web服务
	go runGin()
	wg.Wait()
}

func controlListen() {
	listener, err := helper.CreateListen(define.ControlServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("[控制中心] 监听中 %s\n", listener.Addr().String())
	for {
		controlConn, err = listener.AcceptTCP()
		if err != nil {
			log.Printf("ControlListen AcceptTCP Error: %v\n", err)
			return
		}
		go helper.KeepAlive(controlConn)
	}
}

func userRequestListen() {
	listener, err := helper.CreateListen(define.UserRequestAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("[用户请求] 监听中 %s\n", listener.Addr().String())
	for {
		userConn, err = listener.AcceptTCP()
		if err != nil {
			log.Printf("UserRequestListen AcceptTCP Error: %v\n", err)
			return
		}
		// 发送消息，告诉客户端有新的连接
		_, err := controlConn.Write([]byte(define.NewConnection))
		if err != nil {
			log.Printf("发送失败 %v\n", err)
		}
	}
}

func tunnelListen() {
	listener, err := helper.CreateListen(define.TunnelServerAddr)
	if err != nil {
		panic(err)
	}
	log.Printf("[隧道] 监听中 %s\n", listener.Addr().String())
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Printf("TunnelListen AcceptTCP Error: %v\n", err)
			return
		}
		// 数据转发
		go io.Copy(userConn, conn)
		go io.Copy(conn, userConn)
	}
}

func runGin() {
	server := gin.Default()
	serverConf, err := helper.GetServerConf()
	if err != nil {
		return
	}
	// 用户登录
	server.POST("/login", service.Login)
	server.Run(serverConf.Web.Port)
}
