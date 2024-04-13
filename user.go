package main

import (
	"net"
	"strings"
)

type User struct {
	Name    string
	Address string
	C       chan string
	conn    net.Conn

	Server *Server
}

// 创建用户Api
func NewUser(conn net.Conn, Server *Server) *User {
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

// 给对应客户端发送消息
func (this *User) SendMessage(msg string) {
	this.conn.Write([]byte(msg))
}

// 处理消息
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		//查询当前所有在线用户
		this.Server.MapLock.Lock()
		for _, user := range this.Server.OnlineMap {
			onlineMsg := "[" + this.Address + "]" + user.Name + " is online...\n"
			this.SendMessage(onlineMsg)
		}
		this.Server.MapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" { //rename后需要+名字所以长度判断
		//消息格式：rename|张三
		newName := strings.Split(msg, "|")[1]

		//判断nane是否存在
		_, ok := this.Server.OnlineMap[newName]
		if ok {
			this.SendMessage("newName is exist...")
		} else {
			this.Server.MapLock.Lock()
			delete(this.Server.OnlineMap, this.Name)
			this.Server.OnlineMap[newName] = this
			this.Server.MapLock.Unlock()
			this.Name = newName
			this.SendMessage("rename success.." + this.Name + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		//私聊：消息格式：to|zhangsan|消息内容
		//获取对方用户名
		remoteName := strings.Split(msg, "|")[1]
		if remoteName == "" {
			this.SendMessage("incorrect message format...->to|userName|message\n")
			return
		}
		//根据用户名找到对方suer对象
		remoteUser, ok := this.Server.OnlineMap[remoteName]
		if !ok {
			this.SendMessage("user is unExist...")
		}

		//获取消息内容，通过对方user对象发送
		content := strings.Split(msg, "|")[2]
		if content == "" {
			this.SendMessage("message is empty\n")
			return
		}
		remoteUser.SendMessage(this.Name + " Say to you:" + content)
	} else {
		this.Server.Broad(this, msg)
	}

}

// 监听当前用户的channel,一旦有消息，发送给对应客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\n"))
	}
}
