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
	DisplayName string            `yaml:"display_name"`
	Description string            `yaml:"description"`
	Objects     []string          `yaml:"objects,omitempty"`
	Exits       map[string]string `yaml:"exits,omitempty"`
	clients     []int             `yaml:"clients,omitempty"`
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
	Debug      io.Writer
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
	Rooms     map[string]*Room
	startRoom *Room
	Output    io.Writer
	Debug     io.Writer
	broadcast chan Msg
	ctx       context.Context
	cancel    context.CancelFunc
	nextID    int
}

func NewACServer(port int) *ACServer {
	ctx, cancel := context.WithCancel(context.Background())
	defaultRoom := &Room{
		DisplayName: "default",
		Description: "default room",
	}
	s := &ACServer{
		Name:    "Default server",
		Output:  os.Stdout,
		Debug:   io.Discard,
		Address: fmt.Sprintf(":%d", port),
		Rooms: map[string]*Room{
			"default": defaultRoom,
		},
		startRoom: defaultRoom,
		ctx:       ctx,
		cancel:    cancel,
		broadcast: make(chan Msg, 1),
	}

	return s
}

func (s *ACServer) Start() {
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
}

func (s *ACServer) Print(args ...any) {
	fmt.Fprintln(s.Output, args...)
}

func (s *ACServer) HandleConn(conn net.Conn) {
	defer conn.Close()

	client := newClient(s.nextID, s.startRoom, conn)
	client.Debug = s.Debug
	s.nextID++
	fmt.Fprintln(client.connection, client.room.Description)
	s.Print("client connected, ID ", client.ID)

	scanner := bufio.NewScanner(conn)
	connClients = append(connClients, client)

	for scanner.Scan() {
		var text string
		command, err := Parse(scanner.Text())
		if err != nil {
			fmt.Fprintln(client.connection, "I don't know what you are saying")
		}
		switch command.Verb {
		case "say":
			text = client.Say(command.Object)
		case "go":
			move := client.room.Exits[command.Object]
			client.room = s.Rooms[move]
			text = client.Go(command.Object)
		}
		s.broadcast <- Msg{client.ID, text, client.room}
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
	s.Rooms[room.DisplayName] = room
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
					fmt.Fprintln(client.connection, msg.text)
				}
			}
		}
	}
}

func (s *ACServer) SetStartRoom(name string) {
	s.startRoom = s.Rooms[name]
}
