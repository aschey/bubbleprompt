package test

import (
	"testing"

	"github.com/aschey/bubbleprompt/input/commandinput"
	"github.com/aschey/bubbleprompt/suggestion"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBubbleprompt(t *testing.T) {
	suggestion.DefaultNameForeground = "15"
	suggestion.DefaultSelectedNameForeground = "8"

	suggestion.DefaultDescriptionForeground = "15"
	suggestion.DefaultDescriptionBackground = "13"
	suggestion.DefaultSelectedDescriptionForeground = "8"
	suggestion.DefaultSelectedDescriptionBackground = "13"

	suggestion.DefaultScrollbarColor = "8"
	suggestion.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	RegisterFailHandler(Fail)
	RunSpecs(t, "Bubbleprompt Suite")
}
