package test

import (
	"testing"

	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/commandinput"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBubbleprompt(t *testing.T) {
	input.DefaultNameForeground = "15"
	input.DefaultSelectedNameForeground = "8"

	input.DefaultDescriptionForeground = "15"
	input.DefaultDescriptionBackground = "13"
	input.DefaultSelectedDescriptionForeground = "8"
	input.DefaultSelectedDescriptionBackground = "13"

	input.DefaultScrollbarColor = "8"
	input.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	RegisterFailHandler(Fail)
	RunSpecs(t, "Bubbleprompt Suite")
}
