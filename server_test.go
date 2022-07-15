package adventurechat_test

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"testing"

	"github.com/gamefiend/adventurechat"
	"github.com/google/go-cmp/cmp"
	"github.com/phayes/freeport"
)

var debug *bool

func TestMain(m *testing.M) {
	debug = flag.Bool("debug", false, "Enable debug output")
	flag.Parse()
	os.Exit(m.Run())
}

func TestServerAllowsClientConnectionsProperly(t *testing.T) {
	s := newTestServer(t)
	s.Start()
	newTestClient(t, s)
	newTestClient(t, s)
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
	c.expectMessage("Greetings, welcome to the adventure chat server.")
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
	c.expectMessage("This is a test room.")
}

func newTestServer(t *testing.T) *adventurechat.ACServer {
	t.Helper()
	port, err := freeport.GetFreePort()
	if err != nil {
		t.Fatal(err)
	}
	s := adventurechat.NewACServer(port)
	s.Output = io.Discard
	if *debug {
		s.Output = os.Stdout
	}
	t.Cleanup(func() {
		s.Shutdown()
	})
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
		scanner:    bufio.NewScanner(conn),
	}
}

type testClient struct {
	connection net.Conn
	t          *testing.T
	scanner    *bufio.Scanner
}

func (c *testClient) simulateCmd(msg string) {
	c.t.Log("client sending:", msg)
	fmt.Fprintln(c.connection, msg)
}

func (c *testClient) GetNextMessage() string {
	if !c.scanner.Scan() {
		c.t.Fatal(c.scanner.Err())
	}
	got := c.scanner.Text()
	c.t.Log("client received:", got)
	return got
}

func (c *testClient) expectMessage(want string) {
	got := c.GetNextMessage()
	if !cmp.Equal(want, got) {
		c.t.Error(cmp.Diff(want, got))
	}
}
