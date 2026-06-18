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
	Conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	client := &Client{ServerIp: serverIp, ServerPort: serverPort, flag: 99}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net Dial err: ", err)
		return nil
	}
	client.Conn = conn
	return client
}

func (this *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更改名字")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		this.flag = flag
		return true
	} else {
		fmt.Println("请输入合法数字")
		return false
	}
}

func (this *Client) updateName() bool {
	fmt.Println(">>>>>请输入名字：")
	fmt.Scanln(&this.Name)
	sendMsg := "rename|" + this.Name + "\n"
	_, err := this.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send msg err: ", err)
		return false
	}
	return true
}

func (this *Client) PublicChat() {
	fmt.Println(">>>>>请输入消息，exit退出：")
	var msg string
	fmt.Scanln(&msg)
	for msg != "exit" {
		sendMsg := msg + "\n"
		_, err := this.Conn.Write([]byte(sendMsg))
		if err != nil {
			fmt.Println("send msg err: ", err)
			break
		}
		msg = ""
		fmt.Println(">>>>>请输入消息，exit退出：")
		fmt.Scanln(&msg)
	}

}

func (this *Client) SelectUser() {
	sendMsg := "who\n"
	_, err := this.Conn.Write([]byte(sendMsg))
	if err != nil {
		fmt.Println("send msg err: ", err)
		return
	}
}

func (this *Client) PrivateChat() {
	var name string
	var chatmsg string
	this.SelectUser()
	fmt.Println(">>>>>请输入用户名，exit退出：")
	fmt.Scanln(&name)
	for name != "exit" {
		fmt.Println(">>>>>请输入消息，exit退出：")
		fmt.Scanln(&chatmsg)
		for chatmsg != "exit" {
			sendMsg := "to|" + name + "|" + chatmsg + "\n"
			_, err := this.Conn.Write([]byte(sendMsg))
			if err != nil {
				fmt.Println("send msg err: ", err)
				break
			}
			chatmsg = ""
			fmt.Println(">>>>>请输入消息，exit退出：")
			fmt.Scanln(&chatmsg)
		}
		this.SelectUser()
		chatmsg = ""
		name = ""
		fmt.Println(">>>>>请输入用户名，exit退出：")
		fmt.Scanln(&name)
	}

}
func (this *Client) Run() {
	for this.flag != 0 {
		for this.menu() != true {
		}

		switch this.flag {
		case 1:
			this.PublicChat()
			break
		case 2:
			this.PrivateChat()
			break
		case 3:
			this.updateName()
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "ip", "127.0.0.1", "服务器Ip 默认127.0.0.1")
	flag.IntVar(&serverPort, "port", 8080, "服务器Port 默认8080")
}

func (this *Client) DealResponse() {
	io.Copy(os.Stdout, this.Conn)
	// for{
	// 	buf := make([]byte,4096)
	// 	this.Conn.Read(buf)
	// 	fmt.Println(buf)
	// }
}

func main() {
	flag.Parse()
	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println("链接服务器失败")
		return
	}
	go client.DealResponse()
	fmt.Println("链接服务器成功")
	client.Run()
}
