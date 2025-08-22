# 💬 Local gRPC Chat

A simple **real-time chat application** built with **Go** and **gRPC**.  
This project demonstrates how to use **bidirectional streaming** in gRPC to implement a group chat server and client.

---

## ✨ Features
- 🔗 **gRPC bidirectional streaming**  
- 👥 Multiple clients can join and chat in real-time  
- 📝 Each client is identified by a username  
- 📡 Server broadcasts all messages to connected clients  
- ⚡ Lightweight & concurrent (powered by Go routines and channels)  

---

## 📂 Project Structure

```
.
├── proto/
│   └── chat.proto       # Protocol Buffers definition
├── server/
│   └── main.go          # Chat server implementation
├── client/
│   └── main.go          # Chat client implementation
├── go.mod               # Go module file
└── README.md            # Project documentation
```

---

## 🔧 Prerequisites
Make sure you have installed:
- [Go](https://go.dev/) **v1.18+**
- [Protocol Buffers Compiler (protoc)](https://grpc.io/docs/protoc-installation/)
- gRPC & Protobuf Go plugins:
  ```bash
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
  ```

---

## 📜 Proto Definition

Create `proto/chat.proto`:

```protobuf
syntax = "proto3";

package todo;

option go_package = "./proto";

service ChatService {
    rpc JoinChat(stream Chat) returns (stream Chat);
}

message Chat {
    string user = 1;
    string message = 2;
}
```

Generate Go code from proto:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/chat.proto
```

---

## 🚀 Run the Project

### 1. Initialize Go module
```bash
go mod init github/localChatRouteGrpc
go mod tidy
```

### 2. Start the Server
```bash
cd server
go run main.go
```

Output:
```
2025/08/22 12:00:00 Server started on [::]:8080
```

### 3. Start a Client
```bash
cd client
go run main.go
```

You'll be prompted for a username:
```
Enter your username: Alice
Type messages (type 'exit' to quit):
> 
```

Now open another terminal and start another client:
```
Enter your username: Bob
```

---

## 🖥️ Demo

**Alice terminal:**
```
Enter your username: Alice
Type messages (type 'exit' to quit):
> Hello Bob!
[Bob]: Hi Alice!
> 
```

**Bob terminal:**
```
Enter your username: Bob
Type messages (type 'exit' to quit):
> Hi Alice!
[Alice]: Hello Bob!
> 
```

---

## ⚙️ How It Works

- **Server**:
  - Accepts new client connections via `JoinChat`
  - Registers each client with a username and message channel
  - Broadcasts all messages to connected clients
  - Handles client disconnections gracefully

- **Client**:
  - Connects to the server and registers with a username
  - Runs two goroutines:
    - One for receiving broadcast messages
    - One for sending messages typed by the user
  - Provides clean console interface for chatting

---

## 📌 Notes

- The server listens on port **8080** by default
- To exit a client session, type `exit`
- Each client has a message buffer of 100 messages
- System messages are prefixed with `[System]`

---

## 🤝 Contributing

Pull requests and suggestions are welcome!  
If you'd like to contribute, please fork the repo and create a PR 🚀

---

## 📜 License

This project is licensed under the **MIT License**.

---

## 🔧 Troubleshooting

### Common Issues:

1. **Port already in use**: Change the port in server/main.go
2. **Protoc not found**: Install protocol buffers compiler
3. **Import errors**: Run `go mod tidy` to download dependencies

### Dependencies:
```bash
go get google.golang.org/grpc
go get google.golang.org/protobuf
```

---

Enjoy chatting! 🎉
