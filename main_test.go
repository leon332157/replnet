package main

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

var _ = Describe("dotreplit loader function", func() {
	It("should load valid config with no errors", func() {
		content := []byte(`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		
		[replish]
		port = 7373`)
		loadDotreplit(content)
		Expect(hasReplishField).To(BeTrue())
	})
	It("should fail on invalid config", func() {
		content := []byte(`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		
		[replish]
		port = 7373`)
		loadDotreplit(content)
		Expect(hasReplishField).To(BeFalse()) //  NEED TO USE FUNC RETURN VALUE 
	})
})
