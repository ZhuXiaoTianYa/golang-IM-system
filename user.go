package main

import (
	"net"
	"strings"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

func NewUser(server *Server, conn net.Conn) *User {
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()

	return user
}

func (t *User) ListenMessage() {
	for {
		msg := <-t.C
		t.conn.Write([]byte(msg + "\n"))
	}
}

func (this *User) Online() {
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "已上线")
}

func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "已下线")
}

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":在线...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		newName := strings.Split(msg, "|")[1]
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.SendMsg("当前用户名已被使用\n")
		} else {
			this.server.mapLock.Lock()
			delete(this.server.OnlineMap, this.Name)
			this.server.OnlineMap[newName] = this
			this.server.mapLock.Unlock()
			this.Name = newName
			this.SendMsg("您已更新用户名：" + newName + "\n")
		}
	} else if len(msg) > 4 && msg[:3] == "to|" {
		remoteName := strings.Split(msg, "|")[1]
		remoteContent := strings.Split(msg, "|")[2]
		if remoteName == "" {
			this.SendMsg("消息格式不正确。请使用\nto|张三|消息\n")
			return
		}
		remoteUser, ok := this.server.OnlineMap[remoteName]
		if !ok {
			this.SendMsg("用户不在线\n")
		} else {
			remoteUser.SendMsg(this.Name + "：" + remoteContent + "\n")
		}
	} else {
		this.server.BroadCast(this, msg)
	}
}
