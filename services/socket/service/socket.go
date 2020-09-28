package socket

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// Client Struct
type Client struct {
	hub     *Hub
	connect *websocket.Conn
	send    chan []byte
}

// Hub Struct
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
}

// NewHub for Initial
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// HubRun for Initial
func (h *Hub) HubRun() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("Client Connected")
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				log.Println("Client Disconnected")
				delete(h.clients, client)
				close(client.send)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// Handler Socket
func Handler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	connect, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	client := &Client{
		hub:     hub,
		connect: connect,
		send:    make(chan []byte, 256),
	}

	client.hub.register <- client

	chkStat(client)
}

func chkStat(client *Client) {
	for {
		_, _, err := client.connect.ReadMessage()
		if err != nil {
			client.hub.unregister <- client
			break
		}
	}
}

// SendMSG Socket
func SendMSG(hub *Hub, w http.ResponseWriter, r *http.Request) {

	var msg interface{}

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Println(err.Error())
	}

	message, err := json.Marshal(&msg)
	if err != nil {
		log.Println(err.Error())
	}

	for client := range hub.clients {
		err := client.connect.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			client.hub.unregister <- client
			log.Println(err.Error())
		}
	}

	log.Print(msg)

	json.NewEncoder(w).Encode(msg)
}
