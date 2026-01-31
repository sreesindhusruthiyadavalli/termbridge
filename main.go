package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/creack/pty"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
var jwtSecret = []byte("termbridge-super-secret-2026-change-me")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateToken(username string) (string, error) {
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func validateToken(tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return nil, err
	}
	return token.Claims.(*Claims), nil
}

func main() {
	port := flag.String("port", "8080", "server port")
	flag.Parse()

	r := gin.Default()

	// PUBLIC: Login endpoint
	r.POST("/api/login", func(c *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request"})
			return
		}
		if req.Username == "admin" && req.Password == "password" {
			token, err := generateToken(req.Username)
			if err != nil {
				c.JSON(500, gin.H{"error": "Token generation failed"})
				return
			}
			c.JSON(200, gin.H{"token": token})
			return
		}
		c.JSON(401, gin.H{"error": "Invalid credentials"})
	})

	// PUBLIC: Static files (login page)
	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	// WS handler with MANUAL auth
	r.GET("/ws", func(c *gin.Context) {
		// Extract token from query param
		tokenStr := c.Query("token")
		if tokenStr == "" {
			http.Error(c.Writer, "Token required", http.StatusUnauthorized)
			return
		}

		// Validate JWT
		claims, err := validateToken(tokenStr)
		if err != nil {
			log.Printf("Token validation failed: %v", err)
			http.Error(c.Writer, "Invalid token", http.StatusUnauthorized)
			return
		}

		// UPGRADE to WebSocket AFTER auth
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WS upgrade failed: %v", err)
			return
		}

		// Welcome message with username
		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("\r\n🚀 Welcome %s to TermBridge!\r\n", claims.Username)))

		// Shell setup
		shellCmd := exec.Command("/bin/bash")
		shellCmd.Env = append(os.Environ(), "TERM=xterm-256color")

		ptyShell, err := pty.Start(shellCmd)
		if err != nil {
			log.Printf("PTY spawn failed: %v", err)
			conn.Close()
			return
		}

		log.Printf("Shell session started for %s", claims.Username)
		pty.Setsize(ptyShell, &pty.Winsize{Cols: 120, Rows: 40})

		// Shell → WS
		go func() {
			defer ptyShell.Close()
			defer conn.Close()
			buf := make([]byte, 4096)
			for {
				n, err := ptyShell.Read(buf)
				if err != nil {
					break
				}
				conn.WriteMessage(websocket.TextMessage, buf[:n])
			}
		}()

		// WS → Shell
		go func() {
			defer ptyShell.Close()
			defer conn.Close()
			for {
				msgType, msg, err := conn.ReadMessage()
				if err != nil {
					break
				}
				if msgType == websocket.TextMessage {
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
						continue
					}
				}
				ptyShell.Write(msg)
			}
		}()
	})

	fmt.Printf("TermBridge running on :%s\n", *port)
	r.Run(":" + *port)
}
