package adventurechat

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"os"

	"gopkg.in/yaml.v2"
)

type Room struct {
	DisplayName string
	Description string
	Objects     []string
	Exits       []string
	clients     []int
}

type Object struct {
	DisplayName string
	Description string
	Properties  map[string]string
}

var (
	connClients []Client
	ObjectList  map[string]*Object
)

type Client struct {
	ID         int
	connection net.Conn
	room       *Room
}

type Msg struct {
	sender int
	text   string
	room   *Room
}

func newClient(id int, room *Room, conn net.Conn) Client {
	return Client{
		ID:         id,
		connection: conn,
		room:       room,
	}
}

type ACServer struct {
	Name      string
	Address   string
	Rooms     []*Room
	Output    io.Writer
	broadcast chan Msg
	ctx       context.Context
	cancel    context.CancelFunc
	nextID    int
}

func NewACServer(port int) *ACServer {
	ctx, cancel := context.WithCancel(context.Background())
	s := &ACServer{
		Name:    "Default server",
		Output:  os.Stdout,
		Address: fmt.Sprintf(":%d", port),
		Rooms: []*Room{
			{
				DisplayName: "The gray limbo",
				Description: "You see nothing interesting here.",
			},
		},
		ctx:       ctx,
		cancel:    cancel,
		broadcast: make(chan Msg, 1),
	}
	listener, err := net.Listen("tcp", s.Address)
	if err != nil {
		panic(err)
	}
	go s.broadcastConn()
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				panic(err)
			}
			go s.HandleConn(conn)

		}
	}()
	return s
}

func (s *ACServer) Print(args ...any) {
	fmt.Fprintln(s.Output, args...)
}

func (s *ACServer) HandleConn(conn net.Conn) {
	defer conn.Close()

	client := newClient(s.nextID, s.Rooms[0], conn)
	s.nextID++
	s.Print("client connected, ID ", client.ID)
	scanner := bufio.NewScanner(conn)
	connClients = append(connClients, client)
	fmt.Fprintln(client.connection, client.room.Description)
	for scanner.Scan() {
		s.broadcast <- Msg{client.ID, scanner.Text(), client.room}
	}
}

func (s *ACServer) LoadRoom(path string) error {
	s.Print("Loading room", path)
	config, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	room := &Room{}
	err = yaml.Unmarshal(config, room)
	if err != nil {
		return fmt.Errorf("parse error %v: %q", err, config)
	}
	s.Rooms = append(s.Rooms, room)
	return nil
}

func (s *ACServer) Shutdown() {
	s.cancel()
}

func (s *ACServer) broadcastConn() {
	for {
		select {

		case <-s.ctx.Done():
			return
		default:
			msg := <-s.broadcast
			for _, client := range connClients {
				if msg.sender == client.ID {
					continue
				}
				if msg.room == client.room {
					fmt.Fprintln(client.connection, msg.sender, ":", msg.text)
				}
			}
		}
	}
}
