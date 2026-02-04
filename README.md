<b>Method 1: Global Install (Recommended)</b>

``` 
# 1. Install globally (Go 1.22+ required)
go install github.com/sreesindhusruthiyadavalli/termbridge@latest

# 2. Run server
termbridge

# 3. Browser: http://localhost:8080
# Login: admin / password → Full terminal access!
```

<b>Method 2: Clone & Run (No Go install needed)</b>

```
# 1. Clone repo
git clone https://github.com/<YOURUSERNAME>/termbridge.git
cd termbridge

# 2. Run (auto-downloads dependencies)
go run main.go

# 3. http://localhost:8080 → admin/password
```

<b>Method 3: Docker (Production)</b>

```
# 1. Build & run
docker build -t termbridge .
docker run -p 8080:8080 termbridge

# 2. http://localhost:8080 → admin/password
```

<b>README Usage Section</b>

```
## 🚀 Quick Start

```bash
go install github.com/sreesindhusruthiyadavalli/termbridge@latest
termbridge

```

<b>Browser: http://localhost:8080</b>
<b>Login: admin / password</b>

<b>Flags</b>
```
termbridge -port=8000           # Custom port
termbridge -cmd="zsh"          # Custom shell
```

<b>## **What Users Get Instantly**</b>

✅ Single Go binary (5MB)
✅ JWT authentication (admin/password)
✅ Fullscreen xterm.js terminal
✅ ls, htop, vim, kubectl - ALL WORK
✅ Window resize support
✅ Zero configuration

<b>Production deployment</b>

```
## **Production Deployment**
```bash
# With Nginx reverse proxy
docker-compose up

# Or systemd service
sudo cp termbridge /usr/local/bin/
sudo systemctl enable termbridge

```

<b>Expected user experience</b>

```
$ go install github.com/sreesindhusruthiyadavalli/termbridge@latest  
$ termbridge                                           
# Auto-opens browser or shows: http://localhost:8080
# Login screen → admin/password → "🚀 Welcome to TermBridge!" → bash shell

```