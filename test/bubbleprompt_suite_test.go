package test

import (
	"testing"

	"github.com/aschey/bubbleprompt/formatter"
	"github.com/aschey/bubbleprompt/input/commandinput"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBubbleprompt(t *testing.T) {
	formatter.DefaultNameForeground = "15"
	formatter.DefaultSelectedNameForeground = "8"

	formatter.DefaultDescriptionForeground = "15"
	formatter.DefaultDescriptionBackground = "13"
	formatter.DefaultSelectedDescriptionForeground = "8"
	formatter.DefaultSelectedDescriptionBackground = "13"

	formatter.DefaultScrollbarColor = "8"
	formatter.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	RegisterFailHandler(Fail)
	RunSpecs(t, "Bubbleprompt Suite")
}
