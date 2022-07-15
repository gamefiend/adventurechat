package main

import "github.com/gamefiend/adventurechat"

func main() {
	s := adventurechat.NewACServer(4444)
	s.Start()
	select {}
}
