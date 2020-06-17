package pkg01

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	helloResult = "Hello from Pkg01!"
)

func TestPkg01Suite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg01 Suite")
}

var _ = Describe("helloPkg01 function", func() {
	Context("output", func() {
		It(fmt.Sprintf("should be %s", helloResult), func() {
			Expect(HelloPkg01()).Should(BeIdenticalTo(helloResult))
		})
	})
})
