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
	err := InitRemoteLogSystem(WithRemoteAddr("172.24.124.212:55505"), WithRemoteProtocolType("tcp"), WithLogLevel("debug"), WithContext(context.Background()), WithTerminal(false))
	Specify("debug test", func() {
		time.Sleep(time.Second)
		Expect(err).Should(BeNil())
		model := New("test")
		Expect(model).ShouldNot(BeNil())
		model.WithField("identity", "test").Debug("dongcf----------2222---------")
		time.Sleep(time.Second)
	})
})
