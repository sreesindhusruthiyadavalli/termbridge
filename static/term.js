class TermBridge {
    constructor() {
        this.term = new Terminal({
            cursorBlink: true,
            theme: { background: '#1e1e1e', foreground: '#fff', cursor: '#0ff' }
        });
        this.ws = null;
        this.init();
    }

    init() {
        this.term.open(document.getElementById('terminal'));
        document.getElementById('status').textContent = 'Connected! Type commands:';

        this.ws = new WebSocket(`ws://${location.host}/ws`);

        this.ws.onopen = () => {
            console.log('WebSocket connected');
            this.term.write('\r\n$ Welcome to TermBridge shell\r\n');
        };

        this.ws.onmessage = (event) => {
            this.term.write(event.data);
        };

        this.ws.onclose = () => {
            this.term.write('\r\n[Disconnected]\r\n');
        };

        // Handle typing → WebSocket
        this.term.onData((data) => {
            // Only send Enter key (\r) or normal chars - NO auto \n
            if (data === '\r' || data === '\n') {
                this.ws.send('\r');
            } else {
                this.ws.send(data);
            }
        });

        // Handle window resize
        window.addEventListener('resize', () => this.resize());
        this.resize();
    }

    resize() {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            this.term.fit();
            this.ws.send(JSON.stringify({
                resize: true,
                cols: this.term.cols,
                rows: this.term.rows
            }));
        }
    }
}

// Load FitAddon for proper resize
const fitScript = document.createElement('script');
fitScript.src = 'https://cdn.jsdelivr.net/npm/xterm-addon-fit@0.8.0/lib/xterm-addon-fit.js';
fitScript.onload = () => new TermBridge();
document.head.appendChild(fitScript);
