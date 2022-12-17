package socket

import (
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Client struct {
	mu   sync.Mutex
	Id   string
	Pool *Pool
	Conn *websocket.Conn
}

func NewClient(pool *Pool, con *websocket.Conn) *Client {
	return &Client{
		Pool: pool,
		Conn: con,
	}
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func (cl *Client) Read() {
	// unregister client in the end
	defer func() {
		cl.Pool.Unregister <- cl
		cl.Conn.Close()
	}()

	for {
		// read message from websocket
		msgType, p, err := cl.Conn.ReadMessage()
		if err != nil {
			log.Fatalf("error on read from client: %s", err.Error())
			return
		}

		message := Message{
			Type: msgType,
			Body: string(p),
		}
		// send message for all users in chat
		cl.Pool.Broadcast <- message
		log.Printf("message received: %+v", message)
	}
}
