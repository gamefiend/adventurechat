package adventurechat

import (
	"bufio"
	"context"
	"fmt"
	"net"
)

type ACServer struct {
	Name   string
	cancel context.CancelFunc
}

var (
	connClients []Client
)

type Client struct {
	ID         int
	connection net.Conn
}

type Msg struct {
	sender int
	text   string
}

func newClient(id int, conn net.Conn) Client {

	return Client{id, conn}
}
func NewACServer(name string) *ACServer {

	chatServer, err := net.Listen("tcp", "127.0.0.1:4444")
	if err != nil {
		panic(err)
	}
	broadcast := make(chan Msg, 1)
	connPool := 0
	ctx, cancel := context.WithCancel(context.Background())

	go broadcastConn(ctx, broadcast)
	go func() {
		for {
			if connPool >= 5 { // limit connections to 5
				continue
			}
			conn, err := chatServer.Accept()
			if err != nil {
				panic(err)
			}
			connPool += 1
			go func() {
				defer conn.Close()
				client := newClient(connPool, conn)
				scanner := bufio.NewScanner(conn)
				fmt.Println("New connection from", client.ID)
				connClients = append(connClients, client)
				for scanner.Scan() {
					broadcast <- Msg{client.ID, scanner.Text()}
				}
			}()

		}
	}()
	return &ACServer{
		Name:   name,
		cancel: cancel,
	}
}

func broadcastConn(ctx context.Context, broadcast <-chan Msg) {

	for {
		select {

		case <-ctx.Done():
			return
		default:
			msg := <-broadcast
			fmt.Println(msg)
			for _, client := range connClients {
				if msg.sender == client.ID {
					continue
				}
				fmt.Fprintln(client.connection, msg.sender, ":", msg.text)
			}
		}

	}
}

func (acs *ACServer) Shutdown() {
	acs.cancel()
}
