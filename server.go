package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User
	//全局Map加锁
	MapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
	return server
}

// 监听message广播消息channel的goroutine,一旦有消息就发送给全部的在线user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部在线user
		this.MapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.MapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) Broad(user *User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务....
	fmt.Println("链接建立成功....")

	user := NewUser(conn, this)

	//用户上线，先加入OnlineMap并广播当前用户上线消息
	user.Online()

	//接受客户端消息
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			//客户端关闭,用户下线
			if n == 0 {
				user.Offline()
				return
			}
			//有错误提示
			if err != nil && err != io.EOF {
				fmt.Println("conn read err:", err)
				return
			}

			//提取用户消息
			msg := string(buf[:n-1])

			//用户处理消息
			user.DoMessage(msg)
		}
	}()

	//当前handler阻塞
	select {}
}

// 启动服务器接口
func (this *Server) Start() {

	//socket listen
	//127.0.0.1:8888
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Printf("net.Listen err:%v\n", err)
		return
	}

	//close listen socket
	defer listener.Close()

	//启动监听message的goroutine
	go this.ListenMessager()
	//curl.exe --http0.9 127.0.0.1:8888
	for {
		//accept
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener.Accept err:%v\n", err)
			continue
		}

		//do handler
		go this.Handler(conn)
	}

}
