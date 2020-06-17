package main

import (
	"fmt"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	helloResult = "Hello World!"
)

func TestMainSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main Suite")
}

var _ = Describe("hello function", func() {
	Context("output", func() {
		It(fmt.Sprintf("should be %s", helloResult), func() {
			Expect(hello()).Should(BeIdenticalTo(helloResult))
		})
	})
})
