package User

import (
	"IM-System/Server"
	"net"
)

type User struct {
	Name    string
	Address string
	C       chan string
	conn    net.Conn

	Server *Server.Server
}

// 创建用户Api
func NewUser(conn net.Conn, Server *Server.Server) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:    userAddr,
		Address: userAddr,
		C:       make(chan string),
		conn:    conn,
		Server:  Server,
	}
	//启动监听
	go user.ListenMessage()
	return user
}

// 用户上线
func (this *User) Online() {

	//用户上线，先加入OnlineMap
	this.Server.MapLock.Lock()
	this.Server.OnlineMap[this.Name] = this
	this.Server.MapLock.Unlock()

	//广播当前用户上线
	this.Server.Broad(this, "Online...")
}

// 用户下线
func (this *User) Offline() {

	//用户下线，从OnlineMap中移除
	this.Server.MapLock.Lock()
	delete(this.Server.OnlineMap, this.Name)
	this.Server.MapLock.Unlock()

	//广播当前用户下线
	this.Server.Broad(this, "Offline...")
}

// 处理消息
func (this *User) DoMessage(msg string) {

	this.Server.Broad(this, msg)
}

// 监听当前用户的channel,一旦有消息，发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
