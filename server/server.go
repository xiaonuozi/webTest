package main

import (
	"fmt"
	"net"
	"sync"
)

func main() {
	Server()
}

func Server() {
	// 监听地址和端口
	host := "localhost"
	port := 8080

	//创建上线列表
	var m sync.Map

	// 创建UDP地址对象
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	// 创建UDP套接字
	serverConn, err := net.ListenUDP("udp", serverAddr)
	if err != nil {
		fmt.Println("Error listening:", err)
		return
	}
	defer serverConn.Close()

	fmt.Println("Server started, waiting for messages...")

	// 客户端消息通道
	clientMsgChan := make(chan string)

	// 启动并发 goroutine 处理客户端消息
	go handleClientMessagess(serverConn, clientMsgChan, &m)

	for {
		// 接收消息
		buffer := make([]byte, 1024)
		n, clientAddr, err := serverConn.ReadFromUDP(buffer)
		m.Store(clientAddr.String(), "value1")
		if err != nil {
			fmt.Println("Error reading message:", err)
			continue
		}

		// 提取消息内容
		message := string(buffer[:n])

		// 打印客户端消息
		fmt.Printf("Received message from %s: %s\n", clientAddr.String(), message)

		// 转发消息给其他客户端
		go forwardMessage(serverConn, clientMsgChan, clientAddr, message)
	}
}

func handleClientMessagess(serverConn *net.UDPConn, clientMsgChan <-chan string, m *sync.Map) {
	for {
		// 读取来自其他客户端的消息
		message := <-clientMsgChan

		// 遍历所有键值对
		m.Range(func(key, value interface{}) bool {
			fmt.Println(key, value)
			// 向所有客户端广播消息
			clientAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s", key))
			if err != nil {
				fmt.Println("client Error:", err.Error())
			}
			_, err = serverConn.WriteToUDP([]byte(message), clientAddr)
			if err != nil {
				fmt.Println("Error write UDP:", err.Error())
			}
			return true
		})
	}
}

func forwardMessage(serverConn *net.UDPConn, clientMsgChan chan<- string, senderAddr *net.UDPAddr, message string) {
	// 构建转发消息格式："[sender]: [message]"
	forwardMessage := fmt.Sprintf("[%s]: %s", senderAddr.String(), message)

	// 将消息放入客户端消息通道，以供广播处理
	clientMsgChan <- forwardMessage
}
