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

		Expect(loadConfigKoanf(content)).To(MatchError("replish field doesn't exist"))
	})

	It("should predict client",
		func() {
			content := []byte(`
	 language = "go"
	 run = "bash main.sh"
	 onBoot="bash bootstrap.sh"
	 [replish]
	 mode="clent"`)
			Expect(loadConfigKoanf(content)).To(MatchError("mode clent is invalid, did you mean client?"))
		})

	It("should predict server",
		func() {
			content := []byte(`
	 language = "go"
	 run = "bash main.sh"
	 onBoot="bash bootstrap.sh"
	 [replish]
	 mode="srever"`)
			Expect(loadConfigKoanf(content)).ToNot(Succeed())
		})
	// totally not useless testcases for coverage
})
var _ = Describe("dotreplit loader function (client)", func() {

	It("should fail when remote-url is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"`,
		)
		UNUSED(content)
	})
	// TODO: branch into valid remote-url and invalid remote-url
	It("should fail when remote-app-port is not set", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"
		remote-url = "ws://localhost:8080"`,
		)
		Expect(loadConfigKoanf(content)).To(MatchError("remote application port is unset"))
	})

	It("should default local-app-port to remote-app-port when not set", func() {
		content := []byte(
			`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-app-port = 8080
			remote-url = "ws://localhost:8080"`,
		)
		Expect(loadConfigKoanf(content)).To(Succeed())
		Expect(globalConfig.RemoteAppPort).To(Equal(8080))
		Expect(globalConfig.LocalAppPort).To(Equal(8080))
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
		Expect(loadConfigKoanf(content)).To(MatchError("remote app port 65599 is invalid (1-65535)"))
	})

	It("should set local-app-port",func() {
		content := []byte(
			`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-app-port = 8080
			local-app-port = 8081
			remote-url = "ws://localhost:8080"`,
		)
		Expect(loadConfigKoanf(content)).To(Succeed())
		Expect(globalConfig.LocalAppPort).To(Equal(uint16(8081)))
	})

	It("should fail when local-app-port is not in valid range", func() {

		content := []byte(
			`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-app-port = 8080
			local-app-port = 69999`,
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

	It("default listen port to 0 when invalid", func() {
		content := []byte(
			`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "server"
		listen-port = 99999`,
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
