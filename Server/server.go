package Server

import (
	"fmt"
	"net"
)

type Server struct {
	Ip   string
	Port int
}

// 创建server接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn) {
	//当前链接的业务....
	fmt.Println("链接建立成功....")
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
