package test

import (
	"time"

	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func getTester(binDir string) *tuitest.Tester {
	tester, err := tuitest.NewTester(binDir,
		tuitest.WithMinInputInterval(10*time.Millisecond),
		tuitest.WithDefaultWaitTimeout(5*time.Second),
		tuitest.WithErrorHandler(func(err error) error {
			defer GinkgoRecover()
			Expect(err).Error().ShouldNot(HaveOccurred())
			return err
		}))
	Expect(err).ShouldNot(HaveOccurred())
	return tester
}

var cmdTester *tuitest.Tester = nil
var parserTester *tuitest.Tester = nil

var _ = BeforeSuite(func() {
	cmdTester = getTester("./cmdtestapp")
	parserTester = getTester("./parsertestapp")
})

var _ = AfterSuite(func() {
	_ = cmdTester.TearDown()
})
