package adventurechat_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/gamefiend/adventurechat"
)

func TestServerHandlesClientConnectionsProperly(t *testing.T) {
	t.Parallel()
	server := adventurechat.NewACServer("test")
	fmt.Println(server.Name)
	_, err := net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		t.Fatal(err)
	}
	_, err = net.Dial("tcp", "127.0.0.1:4444")
	if err != nil {
		t.Fatal(err)
	}

}
