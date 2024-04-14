package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int //当前客户端模式
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       99,
	}
	//链接server
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	//返回对象
	return client
}

// 处理server回应的消息，直接显示到标准输出
func (client *Client) DealResponse() {
	//一旦client.conn有数据，直接copy到stdout标准输出，永久阻塞监听
	io.Copy(os.Stdout, client.conn)
	//for {
	//	buf := make()
	//	client.conn.Read(buf)
	//	fmt.Println(string(buf))
	//}
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1..Broadcast mode")
	fmt.Println("2..Private mode")
	fmt.Println("3..update userName")
	fmt.Println("0..quit")
	fmt.Scanln(&flag)

	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println("Illegal input...")
		return false
	}
}

func (client *Client) UpdateName() bool {
	fmt.Println("Please enter your user name...")
	fmt.Scanln(&client.Name)
	sendMsg := "rename|" + client.Name + "\n"
	_, err := client.conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("client.conn.Write err:", err)
		return false
	}
	return true
}

func (client *Client) PublicChat() {
	//提示用户输入信息
	var chatMsg string
	fmt.Println("Please enter Chat content,enter 'exit' to exit... ")
	fmt.Scanln(&chatMsg)
	for chatMsg != "exit" {
		//发送给服务器
		if len(chatMsg) != 0 { //消息不为空
			sendMsg := chatMsg + "\n"
			_, err := client.conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("client.conn.Write err:", err)
				break
			}
		}
		chatMsg = ""
		fmt.Println("Please enter Chat content,enter 'exit' to exit... ")
		fmt.Scanln(&chatMsg)
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true { //非法输入时一直循环读取菜单选项
		}
		//根据不同模式处理不同业务
		switch client.flag {
		case 1:
			//广播
			//fmt.Println("Broadcast mode...")
			client.PublicChat()
			break
		case 2:
			//私聊
			fmt.Println("Private mode...")
			break
		case 3:
			//更新用户名
			fmt.Println("update userName...")
			client.UpdateName()
			break
		}
	}
}

var serverIp string
var serverPort int

// ./client -ip 127.0.0.1 -port 8888
func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "server ip address")
	flag.IntVar(&serverPort, "port", 8888, "server port")
}
func main() {
	//命令行解析
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>>>>>Failed to connect to the server..")
		return
	}
	//单独开启一个goroutine监听server返回的信息
	go client.DealResponse()

	fmt.Println(">>>>>>>>>Success to connect to the server..")
	client.Run()
}
