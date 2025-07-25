package main

import (
	"io"
	"log"
	"net"
	"os"
)

// 代理服务器配置
type ProxyConfig struct {
	ListenAddr string // 本地监听地址
	TargetAddr string // 目标服务器地址
}

// 启动代理服务器
func StartProxy(config ProxyConfig) error {
	listener, err := net.Listen("tcp", config.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("代理服务器启动，监听地址: %s，目标地址: %s",
		config.ListenAddr, config.TargetAddr)

	for {
		// 接受客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接失败: %v", err)
			continue
		}

		// 异步处理连接
		go handleConnection(clientConn, config.TargetAddr)
	}
}

// 处理客户端连接
func handleConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	// 连接目标服务器
	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("连接目标服务器失败: %v", err)
		return
	}
	defer targetConn.Close()

	log.Printf("新连接: %s -> %s",
		clientConn.RemoteAddr(), targetAddr)

	// 双向转发数据
	done := make(chan struct{})

	// 客户端到目标服务器
	go func() {
		io.Copy(targetConn, clientConn)
		done <- struct{}{}
	}()

	// 目标服务器到客户端
	go func() {
		io.Copy(clientConn, targetConn)
		done <- struct{}{}
	}()

	<-done
	log.Printf("连接关闭: %s", clientConn.RemoteAddr())
}

func main() {
	// 简单配置示例
	config := ProxyConfig{
		ListenAddr: ":8080",           // 本地监听端口
		TargetAddr: "localhost:10807", // 目标服务器
	}

	// 启动代理
	if err := StartProxy(config); err != nil {
		log.Fatalf("代理服务器启动失败: %v", err)
		os.Exit(1)
	}
}
