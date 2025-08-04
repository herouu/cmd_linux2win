package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:12345")
	if err != nil {
		panic(err)
	}
	tcpConn := conn.(*net.TCPConn)
	_ = tcpConn.SetKeepAlive(true)
	defer conn.Close()

	stdin := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("send > ")
		text, _ := stdin.ReadString('\n')
		text = strings.TrimSpace(text)
		if strings.ToLower(text) == "q" {
			return
		}
		if _, err := conn.Write([]byte(text + "\n")); err != nil {
			fmt.Println("write error:", err)
			return
		}

		resp, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Print("recv < " + resp)
	}
}
