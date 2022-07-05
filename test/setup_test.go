package test

import (
	"time"

	tuitest "github.com/aschey/tui-tester"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var tester *tuitest.Tester = nil

var _ = BeforeSuite(func() {
	var err error
	tester, err = tuitest.NewTester("./_testapp",
		tuitest.WithMinInputInterval(10*time.Millisecond),
		tuitest.WithDefaultWaitTimeout(5*time.Second),
		tuitest.WithErrorHandler(func(err error) error {
			defer GinkgoRecover()
			Expect(err).Error().ShouldNot(HaveOccurred())
			return err
		}))
	Expect(err).ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	_ = tester.TearDown()
})
