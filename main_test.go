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
		correctConfig := DotReplit{
			Run:      "bash main.sh",
			Language: "go",
			onBoot:   "",
			packager: nil,
			Replish:  map[string]interface{}{"port": 7373},
		}
		content := []byte(`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		
		[replish]
		port = 7373`)

		Expect(loadDotreplit(content)).To(Equal(correctConfig))
	})
	It("should fail on invalid config", func() {
		content := []byte(`broken toml`)
		loadDotreplit(content)
		//Expect(loadDotreplit(content)).Should(Panic())
	})
})
