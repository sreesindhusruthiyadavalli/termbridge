# termbridge
Go CLI turning browser WebSockets into remote shell access

# Test1: TermBridge
- go run main.go -port=8080 -cmd=bash
- # New terminal: curl localhost:8080/health → {"status":"ok","cmd":"bash"}
- wscat -c ws://localhost:8080/ws → sends/receives echo