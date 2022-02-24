package main

import (
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/pgbytes/grpc-playground/api/chat"
	"google.golang.org/grpc"
)

type Connection struct {
	conn chat.ChatService_ChatServer
	send chan *chat.ChatMessage
	quit chan struct{}
}

func NewConnection(conn chat.ChatService_ChatServer) *Connection {
	c := &Connection{
		conn: conn,
		send: make(chan *chat.ChatMessage),
		quit: make(chan struct{}),
	}
	go c.start()
	return c
}

func (c *Connection) Close() error {
	close(c.quit)
	close(c.send)
	return nil
}

func (c *Connection) Send(msg *chat.ChatMessage) {
	defer func() {
		// Ignore any errors about sending on a closed channel and recover from the panic
		recover()
	}()
	c.send <- msg
}

func (c *Connection) start() {
	running := true
	for running {
		select {
		case msg := <-c.send:
			c.conn.Send(msg)
		case <-c.quit:
			running = false
		}
	}
}

func (c *Connection) GetMessages(broadcast chan<- *chat.ChatMessage) error {
	for {
		msg, err := c.conn.Recv()
		if err == io.EOF {
			c.Close()
			return nil
		} else if err != nil {
			c.Close()
			return err
		}
		// in a separate go routine as it can still receive and then continue receiving more if the other end is not ready
		go func(msg *chat.ChatMessage) {
			select {
			case broadcast <- msg:
			case <-c.quit:
			}
		}(msg)
	}
}

type ChatServer struct {
	broadcast   chan *chat.ChatMessage
	quit        chan struct{}
	connections []*Connection
	connLock    sync.Mutex
}

func newChatServer() *ChatServer {
	srv := &ChatServer{
		broadcast: make(chan *chat.ChatMessage),
		quit:      make(chan struct{}),
	}
	go srv.start()
	return srv
}

func (c *ChatServer) Close() error {
	close(c.quit)
	return nil
}

func (c *ChatServer) start() {
	running := true
	for running {
		select {
		case msg := <-c.broadcast:
			c.connLock.Lock()
			for _, v := range c.connections {
				go v.Send(msg)
			}
			c.connLock.Unlock()
		case <-c.quit:
			running = false
		}
	}
}

func (c *ChatServer) Chat(stream chat.ChatService_ChatServer) error {
	conn := NewConnection(stream)

	c.connLock.Lock()
	c.connections = append(c.connections, conn)
	c.connLock.Unlock()

	fmt.Printf("%s connected \n", "user")

	err := conn.GetMessages(c.broadcast)

	// remove itself from list of connections when disconnecting
	c.connLock.Lock()
	for i, v := range c.connections {
		if v == conn {
			c.connections = append(c.connections[:i], c.connections[i+1:]...)
		}
	}
	c.connLock.Unlock()
	fmt.Printf("%s disconnected \n", "user")

	return err
}

func main() {
	lst, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer()
	chatServer := newChatServer()
	chat.RegisterChatServiceServer(server, chatServer)

	fmt.Println("Serving chat server at port 8080")
	err = server.Serve(lst)
	if err != nil {
		panic(err)
	}
}
