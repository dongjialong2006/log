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

		for i := 0; i < 10000; i++ {
			log.Error("dongcf----------1111---------")
		}

		time.Sleep(time.Second)
	})
})
