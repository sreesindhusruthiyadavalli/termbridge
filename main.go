package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Dev only
}

func main() {
	port := flag.String("port", "8080", "server port")
	cmd := flag.String("cmd", "bash", "shell command")
	flag.Parse()

	r := gin.Default()

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "cmd": *cmd})
	})

	// WS endpoint (Phase 2 stub)
	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("WS upgrade failed: %v\n", err)
			return
		}
		defer conn.Close()
		fmt.Println("WS connected (stub)")
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
			// Echo stub for testing
			conn.WriteMessage(websocket.TextMessage, []byte("echo: pong"))
		}
	})

	fmt.Printf("TermBridge running on :%s (cmd: %s)\n", *port, *cmd)
	r.Run(":" + *port)
}
