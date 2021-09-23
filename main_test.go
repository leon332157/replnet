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

var _ = Describe("dotreplit loader function (client)", func() {
	It("should fail on invalid toml", func() {
		content := []byte(`broken`)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
	It("should fail when replish isn't present", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"`,
		)

		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
	It("should default to client mode when mode is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]`,
		)
		loadConfigKoanf(content)
		Expect(globalConfig.Mode).To(Equal("client"))
	})
	It("should fail when client is set and remote-port is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"`,
		)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
	It("should fail when client is set and remote is not valid range", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"
		remote-port = 65599`,
		)
		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
	It("should default listen-port to zero when invalid or not set", func() {
		content := []byte(
			`language = "go"
	run = "bash main.sh"
	onBoot="bash bootstrap.sh"
	[replish]
	mode = "server"
	listen-port = 65599`,
		)
		loadConfigKoanf(content)
		Expect(globalConfig.Mode).To(Equal("server"))
		Expect(globalConfig.ListenPort).To(Equal(uint16(0)))
		content = []byte(
			`language = "go"
	run = "bash main.sh"
	onBoot="bash bootstrap.sh"
	[replish]
	mode = "server"`,
		)
		loadConfigKoanf(content)
		Expect(globalConfig.Mode).To(Equal("server"))
		Expect(globalConfig.ListenPort).To(Equal(uint16(0)))
	})
	It("should fail when local application port is invalid", func() {
		content := []byte(
			`language = "go"
	run = "bash main.sh"
	onBoot="bash bootstrap.sh"
	[replish]
	mode = "server"
	LocalAppPort = 65599`,
		)

		Expect(loadConfigKoanf(content)).ToNot(Succeed())
	})
	//TODO:Write more tests
})
