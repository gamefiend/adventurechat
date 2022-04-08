package main

import (
	"bufio"
	"fmt"
	"net"
)

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

func main() {

	chatServer, err := net.Listen("tcp", "127.0.0.1:4444")
	if err != nil {
		panic(err)
	}
	broadcast := make(chan Msg, 1)
	connPool := 0
	go broadcastConn(broadcast)
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

}

func broadcastConn(broadcast <-chan Msg) {

	for {
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
