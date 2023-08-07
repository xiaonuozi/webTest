package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	Client()
}

func Client() {
	// 服务器地址和端口
	serverHost := "localhost"
	serverPort := 8080

	// 创建UDP地址对象
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", serverHost, serverPort))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 创建UDP套接字
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	// 启动并发 goroutine 接收服务器消息
	go receiveMessages(conn)

	// 读取用户输入并发送消息给服务器
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, _ := reader.ReadString('\n')

		// 发送消息
		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func receiveMessages(conn *net.UDPConn) {
	for {
		// 接收消息
		buffer := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		// 提取消息内容
		message := string(buffer[:n])

		// 打印接收到的消息
		fmt.Println("Received message:", message)
	}
}
