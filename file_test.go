package log

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNew(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "New Suite")
}

var _ = Describe("New", func() {
	log1, err := NewLog("", WithLogLevel("debug"), WithTerminal(false), WithLogName("sslvpn-agent1"))
	log2, err := NewLog("sslvpn-agent2", WithLogLevel("debug"), WithTerminal(false), WithLogName("sslvpn-agent2"))
	Specify("debug test", func() {
		Expect(err).Should(BeNil())
		model := log1.NewEntry("dcf")
		Expect(model).ShouldNot(BeNil())
		for i := 0; i < 100; i++ {
			model.Error("dongcf----------2222---------")
		}

		model = log2.NewEntry("dcf")
		Expect(model).ShouldNot(BeNil())
		for i := 0; i < 100; i++ {
			model.Error("dongcf----------3333---------")
		}

		time.Sleep(time.Second)
	})
})
