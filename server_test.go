package adventurechat_test

import (
	"bufio"
	"fmt"
	"net"
	"testing"

	"github.com/gamefiend/adventurechat"
)

const (
	serverAddress = "127.0.0.1:4444"
)

func TestServerAllowsClientConnectionsProperly(t *testing.T) {
	server := adventurechat.NewACServer("test")
	c1, err := net.Dial("tcp", serverAddress)
	if err != nil {
		t.Fatal(err)
	}
	c1.Close()
	c2, err := net.Dial("tcp", serverAddress)
	if err != nil {
		t.Fatal(err)
	}
	c2.Close()

	server.Shutdown()
}

func TestServerCommunicatesMessagesToAllClients(t *testing.T) {

	server := adventurechat.NewACServer("test")
	fmt.Println("Server listening on 4444")
	c1, err := net.Dial("tcp", serverAddress)

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("c1 Connected to server")
	c2, err := net.Dial("tcp", serverAddress)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("c2 Connected to server")

	var got string
	input := fmt.Sprintln("Hello World")
	want := fmt.Sprintln("1 : Hello World")

	// TODO: When we make a client, rework the code to use it
	_, err = c1.Write([]byte(input))
	fmt.Println("c1 wrote to server")
	if err != nil {
		t.Fatal(err)
	}
	got, err = bufio.NewReader(c2).ReadString('\n')
	if err != nil {
		t.Error(err)
	}

	if string(got) != want {
		t.Errorf("got %s, want %s", string(got), want)
	}

	c1.Close()
	c2.Close()
	server.Shutdown()
}

func TestLoadRoomSendsDescriptionToAllConnectedClients(t *testing.T) {

	server := adventurechat.NewACServer("test")
	server.LoadRoom("data/greeting_room.yaml")
	fmt.Println("Server listening on 4444")
	c1, err := net.Dial("tcp", serverAddress)
	if err != nil {
		t.Fatal(err)
	}
	want := "Greetings, welcome to the adventure chat server.\n"
	got, err := bufio.NewReader(c1).ReadString('\n')
	if err != nil {
		t.Error(err)
	}
	if want != got {
		t.Errorf("Wanted: %v	\nGot: %v", want, got)
	}

	c1.Close()
	server.Shutdown()
}
