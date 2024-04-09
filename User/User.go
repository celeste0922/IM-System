package User

import (
	"net"
)

type User struct {
	Name    string
	Address string
	C       chan string
	conn    net.Conn
}

// 创建用户Api
func NewUser(conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Address: userAddr,
		C:       make(chan string),
		conn:    conn,
	}
	//启动监听
	go user.ListenMessage()
	return user
}

// 监听当前用户的channel,一旦有消息，发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
