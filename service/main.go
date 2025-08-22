package main

import (
	"io"
	"log"
	"net"
	"sync"

	pb "github/localChatRouteGrpc/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedChatServiceServer
	clients map[string]chan *pb.Chat
	mu      sync.Mutex
}

func newServer() *server {
	return &server{
		clients: make(map[string]chan *pb.Chat),
	}
}

func (s *server) JoinChat(stream pb.ChatService_JoinChatServer) error {
	// Получаем первое сообщение для регистрации пользователя
	firstMsg, err := stream.Recv()
	if err == io.EOF {
		return nil
	}
	if err != nil {
		return err
	}

	username := firstMsg.User
	log.Printf("User %s joined the chat", username)

	// Создаем канал для пользователя
	userChan := make(chan *pb.Chat, 100)

	// Регистрируем пользователя
	s.mu.Lock()
	s.clients[username] = userChan
	s.mu.Unlock()

	// Убираем пользователя при выходе
	defer func() {
		s.mu.Lock()
		delete(s.clients, username)
		s.mu.Unlock()
		log.Printf("User %s left the chat", username)
	}()

	// Горутина для отправки сообщений пользователю
	go func() {
		for msg := range userChan {
			if err := stream.Send(msg); err != nil {
				log.Printf("Failed to send message to %s: %v", username, err)
				break
			}
		}
	}()

	// Рассылаем сообщение о входе всем пользователям
	s.broadcast(&pb.Chat{
		User:    "System",
		Message: username + " joined the chat",
	})

	// Основной цикл для приема сообщений от пользователя
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Printf("Message from %s: %s", msg.User, msg.Message)

		// Рассылаем сообщение всем клиентам
		s.broadcast(msg)
	}
}

// Функция для рассылки сообщений всем клиентам
func (s *server) broadcast(msg *pb.Chat) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for username, clientChan := range s.clients {
		select {
		case clientChan <- msg:
			// Сообщение отправлено
		default:
			log.Printf("Client %s channel is full, skipping message", username)
		}
	}
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterChatServiceServer(grpcServer, newServer())
	log.Printf("Server started on %v", lis.Addr())

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
