package test

import (
	"testing"

	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/commandinput"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBubbleprompt(t *testing.T) {
	editor.DefaultNameForeground = "15"
	editor.DefaultSelectedNameForeground = "8"

	editor.DefaultDescriptionForeground = "15"
	editor.DefaultDescriptionBackground = "13"
	editor.DefaultSelectedDescriptionForeground = "8"
	editor.DefaultSelectedDescriptionBackground = "13"

	editor.DefaultScrollbarColor = "8"
	editor.DefaultScrollbarThumbColor = "15"

	commandinput.DefaultCurrentPlaceholderSuggestion = "8"

	RegisterFailHandler(Fail)
	RunSpecs(t, "Bubbleprompt Suite")
}
