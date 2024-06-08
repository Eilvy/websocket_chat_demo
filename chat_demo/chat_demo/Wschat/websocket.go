package Wschat

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go_code/chat_demo/chat_demo/resps"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type Client struct {
	hub      *Hub
	Conn     map[string]*websocket.Conn //连接池
	username string
	send     map[string]chan []byte
	//online   map[string]map[string]chan []byte
}

var (
	mux sync.Mutex
)

// Upgrade 升级为websocket的结构体
var Upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var Hub1 = *NewHub()

// WsChat 聊天室主程序
func WsChat(ctx *gin.Context) {
	client1 := Client{
		Conn:     make(map[string]*websocket.Conn),
		username: "talker" + strconv.Itoa(rand.Intn(10000)+1000),
		send:     make(map[string]chan []byte),

		//online:   make(map[string]map[string]chan []byte),
	}
	id := ctx.Param("id")

	go Hub1.Run(id, client1.username)

	//心跳计时器
	pingTicker := time.NewTicker(time.Second * 10)

	//client.online[id][client.username] <- []byte(client.username)

	conn, err := Upgrade.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("websocket upgrade error : ", err)
		resps.InternalErr(ctx)
		return
	}

	Hub1.register <- &client1

	//将打进来的链接接入连接池并绑定ID
	client1.addClient(client1.username, conn)

	//获取客户端的消息通道
	m, exist := client1.getNewsChannel(id)
	if !exist {
		//若不存在消息通道则创建
		m = make(chan []byte)
		client1.addSendClient(id, m)
	}

	conn.SetCloseHandler(func(code int, text string) error {
		client1.deleteClient(id)
		log.Println(code)
		return nil
	})

	mux.Lock()
	room := client1.joinRoom(id)
	mux.Unlock()
	//心跳检测
	//go client1.Heart(id, pingTicker)
	go client1.Read(id, pingTicker)

	go client1.Write(room)
}

//func (c *Client) Heart(id string, pingTicker *time.Ticker) {
//	defer func() {
//		pingTicker.Stop()
//		log.Printf("connect room %s time out ", id)
//		Hub1.unregister <- c
//		//err := c.Conn[c.username].Close()
//		//if err != nil {
//		//	log.Println("chan close error : ", err)
//		//	return
//		//}
//		//c.deleteNewsChannel(id)
//		//c.deleteClient(id)
//	}()
//	for {
//		select {
//		case <-pingTicker.C: //服务端向客户端心跳检测
//			err := c.Conn[c.username].SetWriteDeadline(time.Now().Add(time.Second * 20))
//			if err != nil {
//				log.Println("set writer deadline error : ", err)
//				return
//			}
//			//若超时WriteMessage会返回一个错误
//			err = c.Conn[c.username].WriteMessage(websocket.PingMessage, []byte{})
//			if err != nil {
//				log.Println("server get pong time out ")
//				log.Printf("connect room %s time out ", id)
//				Hub1.unregister <- c
//				err = c.Conn[c.username].Close()
//				if err != nil {
//					log.Println("chan close error : ", err)
//					return
//				}
//				c.deleteNewsChannel(id)
//				c.deleteClient(c.username) //删除心跳消失的客户端链接
//				continue                   //结束循环并回到defer
//
//			}
//		}
//	}
//}

func (c *Client) Read(id string, pingTicker *time.Ticker) {
	defer func() {
		log.Printf("user %s exit chat room %s", c.username, id)
		err := c.Conn[c.username].Close()
		if err != nil {
			log.Println("chan close error : ", err)
			return
		}
	}()

	for {
		select {
		case msg := <-c.send[id]:
			err := c.Conn[c.username].WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write message error : ", err)
				return
			}
		case <-pingTicker.C:
			err := c.Conn[c.username].SetWriteDeadline(time.Now().Add(time.Second * 20))
			if err != nil {
				log.Println("set writer deadline error : ", err)
				return
			}
			//若超时WriteMessage会返回一个错误
			err = c.Conn[c.username].WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("server get pong time out ")
				log.Printf("connect room %s time out ", id)
				err = c.Conn[c.username].Close()
				if err != nil {
					log.Println("chan close error : ", err)
					return
				}
				c.deleteNewsChannel(id)
				continue //结束循环并回到defer
			}
		}
	}
}
func (c *Client) Write(room *Room) {
	defer func() {
		err := c.Conn[c.username].Close()
		if err != nil {
			log.Println("chan close error : ", err)
			return
		}
		//c.deleteNewsChannel(id)
		//c.deleteClient(id)
	}()

	for {
		msgType, msgByte, err := c.Conn[c.username].ReadMessage()
		if err != nil {
			log.Println("get message error : ", err)
			break
		}
		switch msgType {
		case websocket.TextMessage:
			msg := []byte(fmt.Sprintf("%s %s 说 :%s", time.Now().Format("01/02 06:05"), c.username, msgByte))
			//Hub1.broadcast <- msg
			room.Broadcast <- msg
		}
	}
}
