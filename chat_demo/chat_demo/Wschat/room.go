package Wschat

import (
	"fmt"
	"log"
	"time"
)

// Room 房间广播
type Room struct {
	Id        string
	Members   map[*Client]bool
	Broadcast chan []byte
}

var rooms map[string]*Room

func (c *Client) joinRoom(roomID string) (room *Room) {
	room, ok := rooms[roomID]
	if !ok {
		room = &Room{
			Id:        roomID,
			Members:   make(map[*Client]bool),
			Broadcast: make(chan []byte),
		}
		rooms = make(map[string]*Room)
		rooms[roomID] = room
	}
	room.Members[c] = true
	go c.broadCastMessages(room, c.username)
	return room
}

func (c *Client) broadCastMessages(room *Room, username string) {
	defer func() {
		log.Printf("user %s exit chat room %s", username, room.Id)
		err := c.Conn[c.username].Close()
		if err != nil {
			log.Println("chan close error : ", err)
			return
		}
		leaveRoom(c, room.Id)
		return
	}()
	for {
		select {
		case msg := <-room.Broadcast:
			msg = []byte(fmt.Sprintf("%s %s 说 :%s", time.Now().Format("01/02 06:05"), username, msg))
			for members := range room.Members {
				if members != c {
					select {
					case members.send[room.Id] <- msg:
					default:
						delete(room.Members, members)
					}
				}
			}
		}
	}

}

func leaveRoom(conn *Client, roomID string) {
	room, ok := rooms[roomID]
	if ok {
		delete(room.Members, conn)
		if len(room.Members) == 0 {
			delete(rooms, roomID)
		}
	}
}
