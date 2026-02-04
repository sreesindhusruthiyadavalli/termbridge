# termbridge
Go CLI turning browser WebSockets into remote shell access

<div style="border: 2px solid #8a2be2; padding: 20px; border-radius: 8px; background: #1a1a2e; color: #000000; font-family: monospace;">

# 🚀 TermBridge

**Go CLI → Browser Terminal**  
Single binary • JWT Auth • static UI • PTY Shell

```bash
go install github.com/YOURUSER/termbridge@latest
termbridge
# → http://localhost:8080 (admin/password)
```

[![Demo](./demo.gif)](./demo.gif)

</div>



# Test1: TermBridge
- go run main.go -port=8080 -cmd=bash
- # New terminal: curl localhost:8080/health → {"status":"ok","cmd":"bash"}
- wscat -c ws://localhost:8080/ws → sends/receives echo

# Demo: TermBridge
  $ ls                    # Shows main.go, go.mod
  $ whoami               # Shows your username  
  $ pwd                  # Shows project dir
  $ go version          # Shows Go version
  $ PS1="🚀 TermBridge> "   # Custom prompt
  🚀 TermBridge> htop     # Install if needed: sudo apt install htop
# Resize browser window → watch terminal resize live!
🚀 TermBridge> exit

# Docker run Nginx to serve static files
Terminal1: go run main.go -port=8000

Terminal2: docker run --rm -p 80:80 \
-v $(pwd)/nginx.conf:/etc/nginx/conf.d/default.conf \
-v $(pwd)/static:/usr/share/nginx/html/static:ro \
nginx:alpine

# Access http://localhost/static in browser
http://localhost/           ← Nginx port 80 → Go port 8000
POST localhost/api/login    ← Nginx 80 → Go 8000/api/login  
ws://localhost/ws?token=…   ← Nginx 80 → Go 8000/ws
