package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
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

	//监听用户是否活跃的channel
	isLive := make(chan bool)

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

			//提取用户消息（去除最后的换行）
			msg := string(buf[:n-1])

			//用户处理消息
			user.DoMessage(msg)

			//用户活跃
			isLive <- true
		}
	}()

	//当前handler阻塞
	for {
		select {
		case <-isLive:
			//当前用户活跃，激活select重置定时器
		case <-time.After(time.Second * 300):
			//已经超时
			//将当前客户端强制关闭
			user.SendMessage("Forced offline...")

			//销毁用户资源
			close(user.C)

			//关闭连接
			conn.Close()

			//退出handler
			return //runtime.Goexit()
		}
	}

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
