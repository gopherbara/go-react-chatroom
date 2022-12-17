package socket

import "log"

// pool of users in chat
type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	Clients    map[*Client]bool
	Broadcast  chan Message // transfer for all clients
}

func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.Clients[client] = true
			log.Printf("new client added. Size of connection pool: %v", len(pool.Clients))
			for client, _ := range pool.Clients {
				//send all users message about new one
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "New user joined",
				})
			}
			break
		case client := <-pool.Unregister:
			// remove client when he end connection
			delete(pool.Clients, client)
			log.Printf("client was removed from pool. Size of connection pool: %v", len(pool.Clients))
			for client, _ := range pool.Clients {
				client.Conn.WriteJSON(Message{
					Type: 1,
					Body: "User disconnected",
				})
			}
			break
		case message := <-pool.Broadcast:
			log.Printf("sending message to all clients: %+v", message)
			for client, _ := range pool.Clients {
				if err := client.Conn.WriteJSON(message); err != nil {
					log.Fatalf("Can`t send message to all: %s", err.Error())
					return
				}
			}
			break
		}
	}
}
