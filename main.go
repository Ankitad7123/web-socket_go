package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type Client struct {
	Conn *websocket.Conn
	Send chan []byte
}

func handelWebsocket(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	client := &Client{Conn: ws, Send: make(chan []byte)}

	go client.readMessage()
	go client.writeMessage()

	fmt.Println("connected")
}

func (c *Client) readMessage() {
	defer c.Conn.Close()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		log.Println("Received:", string(message))

	}

}

func (c *Client) writeMessage() {
	defer c.Conn.Close()

	for {
		message, ok := <-c.Send
		if !ok {
			log.Println("Error sending message:", ok)
			break
		}
		err := c.Conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error writing message:", err)
			break
		}
	}
}
