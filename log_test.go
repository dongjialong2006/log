package log

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestLog(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Log Suite")
}

var _ = Describe("Log", func() {
	Specify("debug test", func() {
		log := New("sslvpn-agent", WithLogName("sslvpn-agent"), WithLogLevel("debug"), WithTerminal(false), WithWatchEnable(true))
		Expect(log).ShouldNot(BeNil())

		for i := 0; i < 1000000; i++ {
			log.Errorf("dongcf----------%d---------", i)
		}

		time.Sleep(time.Second)
	})
})
