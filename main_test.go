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

var _ = Describe("dotreplit loader", func() {
	It("should fail on invalid toml", func() {
		content := []byte(`broken`)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})

	It("should fail when replish tag isn't present", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"`,
		)

		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
})
var _ = Describe("dotreplit loader function (client)", func() {

	It("should fail when remote-app-port is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"`,
		)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})

	It("should fail when remote-app-port is not in valid range", func() {
		content := []byte(
			`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-app-port = 65599`,
		)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})

})

var _ = Describe("dotreplit loader function (server)", func() {
	It("default listen port to 0 is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "server"`,
		)
		Expect(loadConfigKoanf(content)).To(Succeed())
		Expect(globalConfig.ListenPort).To(Equal(uint16(0)))
	})

	It("should fail when local-http-port is not in valid range", func() {
		content := []byte(
			`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "server"
			local-http-port = 65599`,
		)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})

})
