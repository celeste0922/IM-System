package Server

import (
	"IM-System/User"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip   string
	Port int

	//在线用户列表
	OnlineMap map[string]*User.User
	//全局Map加锁
	mapLock sync.RWMutex

	//消息广播的channel
	Message chan string
}

// 创建server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User.User),
		Message:   make(chan string),
	}
	return server
}

// 监听message广播消息channel的goroutine,一旦有消息就发送给全部的在线user
func (this *Server) ListenMessager() {
	for {
		msg := <-this.Message

		//将msg发送给全部在线user
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()
	}
}

// 广播消息的方法
func (this *Server) Broad(user *User.User, msg string) {
	sendMsg := "[" + user.Address + "]" + user.Name + ":" + msg
	this.Message <- sendMsg
}

func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务....
	fmt.Println("链接建立成功....")

	user := User.NewUser(conn)

	//用户上线，先加入OnlineMap
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()

	//广播当前用户上线消息
	this.Broad(user, "已上线")

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
