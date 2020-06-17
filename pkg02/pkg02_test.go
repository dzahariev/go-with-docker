package pkg02

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	helloResult = "Hello from Pkg02!"
)

func TestPkg02Suite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Pkg02 Suite")
}

var _ = Describe("HelloPkg02 function", func() {
	Context("output", func() {
		It(fmt.Sprintf("should be %s", helloResult), func() {
			Expect(HelloPkg02()).Should(BeIdenticalTo(helloResult))
		})
	})
})
