package helper

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"gopkg.in/yaml.v3"
	"log"
	"nat-penetration/conf"
	"nat-penetration/define"
	"net"
	"os"
	"time"
)

type UserClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// CreateListen 监听
func CreateListen(serverAddr string) (*net.TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP("tcp", tcpAddr)
	return listener, err
}

// CreateConn 创建连接
func CreateConn(serverAddr string) (*net.TCPConn, error) {
	tcpAddr, err := net.ResolveTCPAddr("tcp", serverAddr)
	if err != nil {
		return nil, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	return tcpConn, err
}

// KeepAlive 设置连接保活
func KeepAlive(conn *net.TCPConn) {
	for {
		_, err := conn.Write([]byte(define.KeepAliveStr))
		if err != nil {
			log.Printf("[KeepAlive] Error %s", err)
			return
		}
		time.Sleep(time.Second * 3)
	}
}

// GetDataFromConnection 获取Connection中的数据
func GetDataFromConnection(bufSize int, conn *net.TCPConn) ([]byte, error) {
	b := make([]byte, 0)
	for {
		buf := make([]byte, bufSize)
		n, err := conn.Read(buf)
		if err != nil {
			return nil, err
		}
		b = append(b, buf[:n]...)
		if n < bufSize {
			break
		}
	}
	return b, nil
}

// GetServerConf 解析 server.yaml
func GetServerConf() (*conf.Server, error) {
	s := new(conf.Server)
	b, err := os.ReadFile("./conf/server.yaml")
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(b, s)
	return s, err
}

var myKey = []byte("nat-penetration-key")

// GenerateToken 生成token
func GenerateToken(name string) (string, error) {
	userClaims := &UserClaims{
		Username: name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 15).Unix(),
		},
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, userClaims)
	token, err := claims.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

// AnalyseToken 解析token
func AnalyseToken(token string) (*UserClaims, error) {
	userClaims := new(UserClaims)
	claims, err := jwt.ParseWithClaims(token, userClaims, func(token *jwt.Token) (interface{}, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	}
	if _, ok := claims.Claims.(*UserClaims); ok && claims.Valid {
		return userClaims, nil
	} else {
		return nil, fmt.Errorf("analyse Token Error: %v", err)
	}
}
