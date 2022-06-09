package prompt

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestBubbleprompt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bubbleprompt Suite")
}
