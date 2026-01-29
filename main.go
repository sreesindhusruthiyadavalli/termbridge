package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"os/exec"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func main() {
	port := flag.String("port", "8080", "server port")
	cmd := flag.String("cmd", "/bin/bash", "shell command") // Fixed: full path
	flag.Parse()

	r := gin.Default()

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "cmd": *cmd})
	})

	r.GET("/ws", func(c *gin.Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WS upgrade failed: %v", err)
			return
		}

		// FIXED: Proper bash command with args
		shellCmd := exec.Command("/bin/bash")
		shellCmd.Env = append(os.Environ(), "TERM=xterm-256color")

		ptyShell, err := pty.Start(shellCmd)
		if err != nil {
			log.Printf("PTY spawn failed: %v", err)
			conn.WriteMessage(websocket.TextMessage, []byte("Shell spawn error"))
			conn.Close()
			return
		}

		log.Println("Shell session started: bash")
		pty.Setsize(ptyShell, &pty.Winsize{Cols: 120, Rows: 40})

		// Shell -> WS (stdout)
		go func() {
			defer ptyShell.Close()
			defer conn.Close()
			buf := make([]byte, 4096)
			for {
				n, err := ptyShell.Read(buf)
				if err != nil {
					log.Printf("Shell read error: %v", err)
					break
				}
				// Send as text for better terminal handling
				conn.WriteMessage(websocket.TextMessage, buf[:n])
			}
		}()

		// WS -> Shell (stdin) - FIXED: Add \r\n for bash
		go func() {
			defer ptyShell.Close()
			defer conn.Close()
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					break
				}

				if msgType == websocket.TextMessage {
					// Check for resize command
					var resize struct {
						Resize bool `json:"resize"`
						Cols   int  `json:"cols"`
						Rows   int  `json:"rows"`
					}
					if json.Unmarshal(msg, &resize) == nil && resize.Resize {
						pty.Setsize(ptyShell, &pty.Winsize{
							Cols: uint16(resize.Cols),
							Rows: uint16(resize.Rows),
						})
						log.Printf("Resized to %dx%d", resize.Cols, resize.Rows)
						continue
					}
				}

				// Convert \n to \r\n for proper bash line handling
				// msg = append(msg, '\r', '\n')
				ptyShell.Write(msg)
			}
		}()
	})

	fmt.Printf("TermBridge running on :%s\n", *port)
	r.Run(":" + *port)
}
