package adventurechat_test

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/gamefiend/adventurechat"
	"github.com/google/go-cmp/cmp"
	"github.com/phayes/freeport"
)

func TestServerAllowsClientConnectionsProperly(t *testing.T) {
	s := newTestServer(t)
	newTestClient(t, s)
	newTestClient(t, s)
	s.Shutdown()
}

func TestServerCommunicatesMessagesToAllClients(t *testing.T) {
	s := newTestServer(t)
	c1 := newTestClient(t, s)
	c2 := newTestClient(t, s)
	want := "Hello World\n"
	c1.Say(want)
	c2.GetNextMessage()
	got := c2.GetNextMessage()

	if !strings.HasSuffix(got, want) {
		t.Errorf(cmp.Diff(want, got))
	}
	s.Shutdown()
}

func TestLoadRoomSendsDescriptionToAllConnectedClients(t *testing.T) {
	s := newTestServer(t)
	s.LoadRoom("data/greeting_room.yaml")
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
		DisplayName: "Test Room",
	}
	s.Rooms = append(s.Rooms, r)
	s.SetStartRoom(r)
	c := newTestClient(t, s)
	want := "Test Room.\n"
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

func (c *testClient) Say(msg string) {
	fmt.Fprintln(c.connection, msg)
}

func (c *testClient) GetNextMessage() string {
	got, err := bufio.NewReader(c.connection).ReadString('\n')
	if err != nil {
		c.t.Fatal(err)
	}
	return got
}
