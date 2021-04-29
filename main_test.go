package main_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func UNUSED(x ...interface{}) {
}

func TestReplish(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

var _ = Describe("dotreplit loader", func() {
	It("should load valid config with no errors", func() {
		
	})
})