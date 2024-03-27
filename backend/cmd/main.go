package main

import (
	"fmt"
	"github.com/gopherbara/go-react-chatroom/backend/pkg/socket"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Start server")
	setupRoutes()
	log.Fatalln(http.ListenAndServe(":4000", nil))
}

func setupRoutes() {
	pool := socket.NewPool()
	go pool.Start()
	http.HandleFunc("/ws", func(writer http.ResponseWriter, request *http.Request) {
		serveWS(pool, writer, request)
	})
}

func serveWS(pool *socket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := socket.Upgrade(w, r)

	if err != nil {
		log.Printf("error on socket connection: %s", err.Error())
		// write to ui
		fmt.Fprintf(w, "error on socket connection: %s", err.Error())
		//return
	}
	client := socket.NewClient(pool, conn)
	pool.Register <- client
	go client.Read()
}
