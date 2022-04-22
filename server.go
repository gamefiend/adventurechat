package adventurechat

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Room struct {
	DisplayName string   `yaml:"displayName"`
	Description string   `yaml:"description"`
	Objects     []string `yaml:"objects"`
	Exits       []string `yaml:"exits"`
}

type Object struct {
	DisplayName string
	Description string
	Properties  map[string]string
}

var (
	connClients []Client
	Rooms       []*Room
	ObjectList  map[string]*Object
)

type Client struct {
	ID          int
	connection  net.Conn
	currentRoom int
}

type Msg struct {
	sender int
	text   string
	room   int
}

func newClient(id, room int, conn net.Conn) Client {
	return Client{id, conn, room}
}

type ACServer struct {
	Name   string
	Rooms  []*Room
	cancel context.CancelFunc
}

func NewACServer(name string) *ACServer {

	chatServer, err := net.Listen("tcp", "127.0.0.1:4444")
	if err != nil {
		panic(err)
	}
	fmt.Println("Server listening on 4444")
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

				client := newClient(connPool, 0, conn)
				scanner := bufio.NewScanner(conn)
				connClients = append(connClients, client)
				for scanner.Scan() {
					broadcast <- Msg{client.ID, scanner.Text(), client.currentRoom}
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
			for _, client := range connClients {
				if msg.sender == client.ID {
					continue
				}
				if msg.room == client.currentRoom {
					fmt.Fprintln(client.connection, msg.sender, ":", msg.text)
				}
			}
		}

	}
}
func (acs *ACServer) LoadRoom(roomName string) error {
	fmt.Println("Loading room", roomName)
	config, err := os.ReadFile(roomName)
	if err != nil {
		return err
	}
	room := &Room{}
	err = yaml.Unmarshal(config, room)
	Rooms = append(Rooms, room)
	return nil
}

func (acs *ACServer) Shutdown() {
	acs.cancel()
}
