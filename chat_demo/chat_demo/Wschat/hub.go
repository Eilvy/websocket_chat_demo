package Wschat

// Hub 广播全体消息
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run(id string, username string) {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				err := client.Conn[username].Close()
				if err != nil {
					return
				}
			}
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send[id] <- message:
				default:
					//err := client.Conn[id].Close()
					//if err != nil {
					//	return
					//}
					delete(h.clients, client)
				}
			}
		}
	}
}
