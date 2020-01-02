package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/panshul007/grpc-playground/proto/chat"
	"google.golang.org/grpc"
	"io"
	"os"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Must have connection string and user name")
		return
	}

	ctx := context.Background()

	conn, err := grpc.Dial(os.Args[1], grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chatClient := chat.NewChatClient(conn)
	stream, err := chatClient.Chat(ctx)
	if err != nil {
		panic(err)
	}

	waitC := make(chan struct{})
	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				close(waitC)
				return
			} else if err != nil {
				panic(err)
			}
			fmt.Printf("%s: %s \n", msg.User, msg.Message)
		}
	}()

	fmt.Println("Connection established, type \"quit\" ctrl+c to exit")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		if msg == "quit" {
			err := stream.CloseSend()
			if err != nil {
				panic(err)
			}
			break
		}

		err := stream.Send(&chat.ChatMessage{
			User:    os.Args[2],
			Message: msg,
		})
		if err != nil {
			panic(err)
		}
	}

	<-waitC
}
