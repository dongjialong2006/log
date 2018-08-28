package log

import (
	"context"
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
	err := InitLocalLogSystem(WithLogLevel("debug"), WithContext(context.Background()), WithTerminal(false), WithLogName("sslvpn-agent"))
	Specify("debug test", func() {
		Expect(err).Should(BeNil())
		model := New("test")
		Expect(model).ShouldNot(BeNil())
		for i := 0; i < 100; i++ {
			model.Error("dongcf----------1111---------")
		}

		time.Sleep(time.Second)
	})
})
