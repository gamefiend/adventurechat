package adventurechat

import (
	"errors"
	"strings"
)

var (
	validCommands = map[string]bool{
		"say": true,
		"go":  true,
	}
)

type Command struct {
	Verb   string
	Object string
}

func Parse(input string) (Command, error) {
	if input == "" {
		return Command{}, errors.New("empty command!")
	}
	parsed := strings.SplitN(input, " ", 2)
	if !validCommands[parsed[0]] {
		return Command{}, errors.New("unknown command")
	}
	return Command{
		Verb:   parsed[0],
		Object: parsed[1],
	}, nil
}
