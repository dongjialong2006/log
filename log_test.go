package log

import (
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var _ = Describe("Service", func() {
	FIt("Parse", func() {
		for j := 0; j < 10; j++ {
			log := New("agent", WithLogName(fmt.Sprintf("test%d", j)), WithLogLevel("debug"), WithTerminal(false), WithWatchEnable(true), WithWatchLogsByNum(5))
			Expect(log).ShouldNot(BeNil())

			for i := 0; i < 100000; i++ {
				log.Errorf("dongcf----------%d---------", i)
			}
		}

		time.Sleep(time.Second * 10)
	})
})
