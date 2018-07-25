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
	Specify("debug test", func() {
		model := New("test")
		Expect(model).ShouldNot(BeNil())
		model.Debug("test")
	})
})
