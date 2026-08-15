<div align="center">

# 🚇 LocalHoist

**Expose your local server to the internet — over a WebSocket tunnel, written in Go.**

[![Go Version](https://img.shields.io/badge/Go-1.25%2B-00ADD8?style=flat&logo=go)](https://go.dev)
[![Gin](https://img.shields.io/badge/Framework-Gin-008ECF?style=flat)](https://github.com/gin-gonic/gin)
[![WebSocket](https://img.shields.io/badge/Protocol-WebSocket-4B32C3?style=flat)](https://github.com/gorilla/websocket)
[![Status](https://img.shields.io/badge/Status-Active%20Development-yellow?style=flat)]()

*A minimal, self-hostable alternative to ngrok — built to understand tunneling internals from the ground up.*

</div>

---

## 📖 Overview

**LocalHoist** lets you take an app running on your machine (say, `localhost:8080`) and give it a public URL that anyone on the internet can hit. It's the same idea behind tools like `ngrok` or `localtunnel`, implemented from scratch in Go as a two-part system:

- **Server** — a Gin HTTP server that accepts WebSocket connections from clients and issues each one a unique public tunnel URL.
- **Client (CLI)** — a Cobra-based command line agent that connects to the server, listens for incoming requests over the socket, replays them against a local port you specify, and streams the response back.

Server and client are independent Go modules (`LocalHoist-Server` and `LocalHoist-Client`), so each can be built, versioned, and deployed on its own.

When someone hits your public tunnel URL, the request travels over the open WebSocket connection to your machine, gets executed locally, and the response is relayed back out — all in real time.

## ⚙️ How It Works

```
                     ┌─────────────────────┐
   Public Request    │   LocalHoist Server  │
  ───────────────────▶   (Gin, port 3000)   │
                     │                      │
                     │  GET /tunnel/:id     │
                     └──────────┬───────────┘
                                │  WebSocket
                                │  (persistent)
                                ▼
                     ┌──────────────────────┐
                     │  LocalHoist CLI       │
                     │  (your machine)       │
                     └──────────┬───────────┘
                                │  HTTP
                                ▼
                     ┌──────────────────────┐
                     │  Your Local App       │
                     │  (localhost:<port>)  │
                     └──────────────────────┘
```

1. You run the **client** with the local port you want to expose: `connect <port> [url]`.
2. The client connects to the server via `GET /ws`, which upgrades the connection to a WebSocket.
3. The server generates a random hex **client ID** and replies with a public tunnel URL: `http(s)://<server-host>/tunnel/<clientID>`.
4. Any HTTP request sent to that public URL is captured by the server, matched to the right client by ID, and forwarded down the WebSocket as JSON.
5. The client decodes the payload, replays it as a real HTTP request against `localhost:<port>`, and writes the response back over the socket.
6. The server picks up that response and returns it to the original caller.

## ✨ Features

- 🔌 **WebSocket-based tunneling** — a single persistent connection carries all traffic for a client, no repeated handshakes.
- 🆔 **Per-client public URLs** — each connected client gets a unique `/tunnel/<id>` route, generated with a cryptographically random ID.
- 🎯 **Configurable target port** — `connect <port>` points the tunnel at any local port, no hardcoded target.
- 🖥️ **Simple CLI** — built on [Cobra](https://github.com/spf13/cobra), so commands and flags are easy to extend.
- 📦 **Raw request forwarding** — method, URL, headers, and body are captured and reconstructed faithfully between server and client.
- 🪶 **Minimal dependencies** — no external infrastructure needed to run the server or client.

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | [Go](https://go.dev) 1.25 |
| HTTP Server | [Gin](https://github.com/gin-gonic/gin) |
| Tunnel Transport | [Gorilla WebSocket](https://github.com/gorilla/websocket) |
| CLI Framework | [Cobra](https://github.com/spf13/cobra) |

## 📂 Project Structure

The server, client, and test app are each their own Go module:

```
LocalHoist/
├── server/                    # module: LocalHoist-Server
│   ├── cmd/
│   │   └── main.go            # Server entrypoint — starts Gin, registers /ws and /tunnel routes
│   ├── ws.go                  # WebSocket upgrade + public URL generation + request forwarding
│   ├── http.go                # Tunnel HTTP handler — parses & relays incoming public requests
│   ├── models.go              # Server-side request/response payload structs
│   ├── utils/
│   │   └── helper.go          # Client ID generation, route parsing
│   ├── go.mod
│   └── go.sum
├── client/                    # module: LocalHoist-Client
│   ├── main.go                # CLI entrypoint (Cobra) — `connect <port> [url]`
│   ├── handle_request.go      # Decodes payloads & executes them against localhost:<port>
│   ├── models.go              # Client-side request/response payload structs
│   ├── go.mod
│   └── go.sum
├── test/                      # module: test
│   ├── main.go                # A minimal local Gin app to test tunneling against
│   ├── go.mod
│   └── go.sum
└── gitignore
```

## 🚀 Getting Started

### Prerequisites

- [Go](https://go.dev/dl/) 1.25 or later

### Clone the repo

```bash
git clone https://github.com/navneetshukl/LocalHoist.git
cd LocalHoist
```

Each module manages its own dependencies, so run `go mod tidy` inside `server/`, `client/`, and `test/` as needed.

### 1. Start the server

```bash
cd server
go mod tidy
go run ./cmd
```

The server starts listening on `:3000` and logs `Server is listening on :3000`.

### 2. Start a local app to expose

A sample app is included for testing:

```bash
cd test
go mod tidy
go run .
```

This starts a demo app on `:8080`.

### 3. Start the LocalHoist client

Point the client at the port your local app is running on:

```bash
cd client
go mod tidy
go run . connect 8080
```

An optional second argument overrides the server URL (defaults to `ws://localhost:3000/ws`):

```bash
go run . connect 8080 ws://localhost:3000/ws
```

You'll see the client connect and print out a public tunnel URL, e.g.:

```
Connecting to ws://localhost:3000/ws...
Connected to Server! Tunnel is open.
✔ Successfully connected to LocalHoist!
🌍 Public URL: http://localhost:3000/tunnel/4cf2fd
```

### 4. Hit the public URL

```bash
curl http://localhost:3000/tunnel/4cf2fd
```

The request is forwarded over the WebSocket to your client, executed against `localhost:8080`, and the response is returned.

> 💡 You can also build standalone binaries: `go build -o localhoist-server ./cmd` (from `server/`) and `go build -o localhoist .` (from `client/`).

## 🧭 Roadmap / Known Limitations

LocalHoist is under active development. A few things currently on the radar:

- [ ] Response bodies are assumed to be JSON; support for arbitrary content types (HTML, images, binary streams) is planned.
- [ ] No TLS or authentication on tunnel connections yet.
- [ ] Client-to-server connection state is in-memory only — no reconnect/session resumption yet.
- [ ] `WSManager.wsConn` is a plain map with no mutex — concurrent connects could race.

Contributions and ideas around these are very welcome!

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

No license has been specified yet for this repository. Consider adding a `LICENSE` file (e.g. MIT) if you intend for others to use or contribute to this project.

---

<div align="center">

Built with ❤️ and Go by [Navneet Shukla](https://github.com/navneetshukl)

</div>
