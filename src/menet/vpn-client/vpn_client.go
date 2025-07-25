package main

import (
	"bufio"
	flag "cmd_linux2win/src/lib/github.com/spf13/pflag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

func main() {
	// 使用 pflag 定义命令行参数
	serverHost := flag.StringP("server", "s", "localhost", "代理服务器地址")
	serverPort := flag.IntP("server-port", "p", 10808, "代理服务器端口")
	localPort := flag.IntP("local-port", "l", 1080, "本地监听端口")
	// 解析命令行参数
	flag.Parse()

	// 监听本地端口
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", *localPort))
	if err != nil {
		log.Fatalf("无法监听本地端口: %v", err)
	}
	defer listener.Close()

	log.Printf("客户端启动，监听本地端口 %d，代理服务器: %s:%d", *localPort, *serverHost, *serverPort)

	// 循环接受本地连接
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("接受本地连接失败: %v", err)
			continue
		}

		// 为每个连接启动一个 goroutine 处理
		go handleLocalConnection(conn, *serverHost, *serverPort)
	}
}

// 处理本地应用程序的连接
func handleLocalConnection(localConn net.Conn, serverHost string, serverPort int) {
	defer localConn.Close()
	log.Printf("新的本地连接: %s", localConn.RemoteAddr())

	// 读取本地应用发送的数据
	reader := bufio.NewReader(localConn)
	request, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Printf("读取本地请求失败: %v", err)
		return
	}

	// 解析目标地址 (简化处理，假设第一行为目标地址)
	targetAddr := strings.TrimSpace(request)
	if targetAddr == "" {
		log.Println("目标地址不能为空")
		return
	}

	// 读取实际请求内容
	requestData, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		log.Printf("读取请求数据失败: %v", err)
		return
	}

	// 构建发送到服务端的数据 (格式: "host:port|request")
	fullData := fmt.Sprintf("%s|%s", targetAddr, requestData)

	// 连接代理服务器
	serverAddr := fmt.Sprintf("%s:%d", serverHost, serverPort)
	serverConn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		log.Printf("连接代理服务器 %s 失败: %v", serverAddr, err)
		return
	}
	defer serverConn.Close()

	// 直接发送明文数据到代理服务器
	_, err = serverConn.Write([]byte(fullData))
	if err != nil {
		log.Printf("发送数据到代理服务器失败: %v", err)
		return
	}

	// 接收代理服务器的响应
	response := make([]byte, 4096)
	n, err := serverConn.Read(response)
	if err != nil {
		log.Printf("读取代理服务器响应失败: %v", err)
		return
	}

	// 直接将响应发送回本地应用
	_, err = localConn.Write(response[:n])
	if err != nil {
		log.Printf("发送响应到本地应用失败: %v", err)
		return
	}

	log.Printf("成功处理到 %s 的请求", targetAddr)
}

// 转发数据的函数，从src读取数据并写入到dst
func proxy(src, dst net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer src.Close()
	defer dst.Close()

	// 从src读取数据并写入到dst
	_, err := io.Copy(dst, src)
	if err != nil && err != io.EOF {
		log.Printf("转发数据错误: %v", err)
	}
}

// 启动代理服务器
func startProxy(listenAddr, targetAddr string) error {
	// 监听指定地址
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("TCP代理已启动，监听地址: %s，目标地址: %s", listenAddr, targetAddr)

	for {
		// 接受客户端连接
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("接受连接错误: %v", err)
			continue
		}

		log.Printf("新的客户端连接: %s", clientConn.RemoteAddr())

		// 连接到目标服务器
		targetConn, err := net.Dial("tcp", targetAddr)
		if err != nil {
			log.Printf("连接目标服务器错误: %v", err)
			clientConn.Close()
			continue
		}

		log.Printf("已连接到目标服务器: %s", targetAddr)

		// 使用WaitGroup等待双向数据转发完成
		var wg sync.WaitGroup
		wg.Add(2)

		// 启动双向数据转发
		go proxy(clientConn, targetConn, &wg)
		go proxy(targetConn, clientConn, &wg)

		// 等待转发完成
		go func() {
			wg.Wait()
			log.Printf("连接已关闭: %s", clientConn.RemoteAddr())
		}()
	}
}
