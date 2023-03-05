package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

type Hub struct {
	clients []*websocket.Conn
}

func (h *Hub) subscribe(client *websocket.Conn) {
	h.clients = append(h.clients, client)
	go h.update(client)
}

func (h *Hub) update(client *websocket.Conn) {
	for {
		_, message, err := client.ReadMessage()
		if err != nil {
			log.Println("read failed: ", err)
			break
		}

		h.publish(message)
	}
}

func (h *Hub) publish(message []byte) {
	output := []byte(message)

	for _, client := range h.clients {
		if err := client.WriteMessage(2, output); err != nil {
			log.Println("write failed:", err)
		}
	}
}

func main() {
	h := &Hub{}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Print("upgrade failed: ", err)
			return
		}

		h.subscribe(conn)
	})

	http.ListenAndServe(":3000", nil)
}
