package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"testing"
)

func UNUSED(x ...interface{}) {
}

func TestReplish(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(true)
	log.SetLevel(log.ErrorLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

var _ = Describe("dotreplit loader function", func() {
	XIt("should load valid config with no errors", func() {
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
	XIt("should fail on invalid config", func() {
		content := []byte(`broken`)
		Expect(loadDotreplit(content))
	})
})
