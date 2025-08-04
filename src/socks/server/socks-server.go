package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:12345")
	if err != nil {
		panic(err)
	}
	fmt.Println("TCP server listening on :12345")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return // 客户端断开
		}
		fmt.Printf("recv: %s", msg)
		conn.Write([]byte("echo: " + msg))
	}
}
