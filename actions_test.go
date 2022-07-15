package adventurechat_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSayActionShowsTextToSpeakerAndOthersInTheSameRoom(t *testing.T) {
	t.Parallel()
	s := newTestServer(t)
	s.Start()
	Alice := newTestClient(t, s)
	Alice.GetNextMessage()
	Bob := newTestClient(t, s)
	Bob.GetNextMessage()
	Alice.simulateCmd("say Hello World")

	want := "says 'Hello World'\n"
	got := Bob.GetNextMessage()
	if !strings.HasSuffix(got, want) {
		t.Error(cmp.Diff(want, got))
	}

	s.Shutdown()
}

func TestGoActionSendsClientToAnotherRoom(t *testing.T) {
	t.Parallel()
	s := newTestServer(t)
	err := s.LoadRoom("./data/greeting_room.yaml")
	if err != nil {
		t.Fatal(err)
	}
	err = s.LoadRoom("./data/lobby.yaml")
	if err != nil {
		t.Fatal(err)
	}
	s.SetStartRoom("Greeting Room")
	fmt.Println("starting server")
	s.Start()
	fmt.Println("Alice connecting")
	Alice := newTestClient(t, s)
	Alice.GetNextMessage()

	Alice.simulateCmd("go north")
	fmt.Println("command issued")
	want := "A grandiose room, filled with the spectacle of nothing"
	got := Alice.GetNextMessage()
	
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}
