package main

import (
	"fmt"
	"io"
	"net"
	"sync"
	"time"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	server := &Server{Ip: ip, Port: port, OnlineMap: make(map[string]*User), Message: make(chan string)}
	return server
}

func (t *Server) ListenMessager() {
	for {
		msg := <-t.Message
		t.mapLock.Lock()
		for _, cli := range t.OnlineMap {
			cli.C <- msg
		}
		t.mapLock.Unlock()
	}
	select {}
}

func (t *Server) BroadCast(user *User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	t.Message <- sendMsg
}

func (t *Server) Handler(conn net.Conn) {
	user := NewUser(t, conn)

	user.Online()

	isLive := make(chan bool)

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				user.Offline()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("Conn Read err: ", err)
				return
			}
			msg := string(buf[:n-1])
			user.DoMessage(msg)
			isLive <- true
		}
	}()
	for {
		select {
		case <-isLive:
		case <-time.After(300 * time.Second):
			user.SendMsg("您已被踢出\n")
			close(user.C)
			user.conn.Close()
			return
		}
	}
}

func (t *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", t.Ip, t.Port))
	if err != nil {
		fmt.Println("net.Listen err:", err)
		return
	}
	defer listener.Close()
	go t.ListenMessager()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("listener accept err:", err)
			continue
		}
		go t.Handler(conn)
	}
}
