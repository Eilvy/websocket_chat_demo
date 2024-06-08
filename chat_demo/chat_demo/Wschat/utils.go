package Wschat

import (
	"github.com/gorilla/websocket"
)

// 加入连接池
func (c *Client) addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	c.Conn[id] = conn
	mux.Unlock()
}

// 获取当前选择的聊天室id和client
func (c *Client) getClient(id string, conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = c.Conn[id]
	mux.Unlock()
	return
}

// 删除客户端连接
func (c *Client) deleteClient(id string) {
	mux.Lock()
	delete(c.Conn, id)
	mux.Unlock()
}

// 添加用户指定消息通道
func (c *Client) addSendClient(id string, m chan []byte) {
	mux.Lock()
	c.send[id] = m
	mux.Unlock()
}

// 获取指定用户消息通道
func (c *Client) getNewsChannel(id string) (m chan []byte, exist bool) {
	mux.Lock()
	m, exist = c.send[id]
	mux.Unlock()
	return
}

// 删除指定消息通道
func (c *Client) deleteNewsChannel(id string) {
	mux.Lock()
	if m, ok := c.send[id]; ok {
		close(m)
		delete(c.send, id)
	}
	mux.Unlock()
}
