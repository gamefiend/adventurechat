package adventurechat_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/gamefiend/adventurechat"
	"github.com/phayes/freeport"
)

func TestServerAllowsClientConnectionsProperly(t *testing.T) {
	s := newTestServer(t)
	s.Start()
	newTestClient(t, s)
	newTestClient(t, s)
	s.Shutdown()
}

func TestLoadRoom_LoadsRoomIntoStartRoom(t *testing.T) {
	s := newTestServer(t)
	err := s.LoadRoom("./data/greeting_room.yaml")
	if err != nil {
		t.Fatal(err)
	}
	s.SetStartRoom("Greeting Room")
	s.Start()
	c := newTestClient(t, s)
	want := "Greetings, welcome to the adventure chat server.\n"
	got := c.GetNextMessage()
	if want != got {
		t.Errorf("Wanted: %v	\nGot: %v", want, got)
	}
	s.Shutdown()
}

func TestSetStartRoom_CausesNewClientsToJoinInGivenRoom(t *testing.T) {
	t.Parallel()
	s := newTestServer(t)
	r := &adventurechat.Room{
		DisplayName: "Test Room.",
		Description: "This is a test room.",
	}
	s.Rooms[r.DisplayName] = r
	s.SetStartRoom(r.DisplayName)
	s.Start()
	c := newTestClient(t, s)
	want := "This is a test room.\n"
	got := c.GetNextMessage()
	if want != got {
		t.Errorf("Wanted: %v	\nGot: %v", want, got)
	}
}

func newTestServer(t *testing.T) *adventurechat.ACServer {
	t.Helper()
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	s := adventurechat.NewACServer(port)
	s.Output = io.Discard
	return s
}

func newTestClient(t *testing.T, s *adventurechat.ACServer) *testClient {
	t.Helper()
	conn, err := net.Dial("tcp", s.Address)
	if err != nil {
		t.Fatal(err)
	}
	return &testClient{
		t:          t,
		connection: conn,
	}
}

type testClient struct {
	connection net.Conn
	t          *testing.T
}

func (c *testClient) simulateCmd(msg string) {
	fmt.Fprintln(c.connection, msg)
}

func (c *testClient) GetNextMessage() string {
	got, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		c.t.Fatal(err)
	}
	return got
}
