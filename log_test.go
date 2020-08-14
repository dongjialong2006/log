package log

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Log", func() {
	Specify("log", func() {
		log := New("sslvpn-agent", WithLogName("sslvpn-agent"), WithLogLevel("debug"), WithTerminal(false), WithWatchEnable(true))
		Expect(log).ShouldNot(BeNil())

		for i := 0; i < 100000; i++ {
			log.Errorf("dongcf----------%d---------", i)
		}
	})
})
