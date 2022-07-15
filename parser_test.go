package adventurechat_test

import (
	"testing"

	"github.com/gamefiend/adventurechat"
	"github.com/google/go-cmp/cmp"
)

func TestParserSplitsValidInputIntoVerbAndObject(t *testing.T) {
	want := adventurechat.Command{
		Verb:   "say",
		Object: "Hello World",
	}
	got, err := adventurechat.Parse("say Hello World")
	if err != nil {
		t.Fatal(err)
	}
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}

}

func TestParserReturnsErrorOnInvalidInputs(t *testing.T) {
	_, err := adventurechat.Parse("")
	if err == nil {
		t.Error("got nil")
	}

}

func TestParserReturnsErrorOnUnknownCommand(t *testing.T) {
	_, err := adventurechat.Parse("blah blah")
	if err == nil {
		t.Error("got nil")
	}
}
