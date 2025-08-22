package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	pb "github/localChatRouteGrpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewChatServiceClient(conn)

	fmt.Print("Enter your username: ")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	username := scanner.Text()

	stream, err := client.JoinChat(context.Background())
	if err != nil {
		log.Fatalf("Failed to join chat: %v", err)
	}

	// Регистрируем пользователя
	if err := stream.Send(&pb.Chat{User: username, Message: "joined the chat"}); err != nil {
		log.Fatalf("Failed to register: %v", err)
	}

	// Горутина для получения сообщений
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				log.Printf("Connection closed: %v", err)
				return
			}
			fmt.Printf("\r[%s]: %s\n> ", msg.User, msg.Message)
		}
	}()

	// Основной цикл для отправки сообщений
	fmt.Println("Type messages (type 'exit' to quit):")
	for {
		fmt.Print("> ")
		scanner.Scan()
		message := scanner.Text()

		if message == "exit" {
			break
		}

		if err := stream.Send(&pb.Chat{User: username, Message: message}); err != nil {
			log.Printf("Failed to send message: %v", err)
			break
		}
	}
}
