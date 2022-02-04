package main

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
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

	When("toml is invalid", func() {
		It("should fail", func() {
			content := []byte(`broken`)
			Expect(loadConfigKoanf(content)).ToNot(Succeed())
		})
	})

	When("replish field isn't present", func() {
		It("should fail when replish tag isn't present", func() {
			content := []byte(
				`language = "go"
run = "bash main.sh"
onBoot="bash bootstrap.sh"`)
			Expect(loadConfigKoanf(content)).To(MatchError("replish field doesn't exist"))
		})
	})

	/*
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
				Expect(loadConfigKoanf(content)).To(MatchError("mode srever is invalid, did you mean server?"))
			})
	*/

	When("mode is invalid", func() {
		It("should fail", func() {
			content := []byte(`
	 language = "go"
	 run = "bash main.sh"
	 onBoot="bash bootstrap.sh"
	 [replish]
	 mode="test"`)
			Expect(loadConfigKoanf(content)).To(MatchError("mode test is invalid"))

		})
	})

	Describe("dotreplit loader (client)", func() {
		When("remote url is not set", func() {
			It("should fail", func() {
				content := []byte(
					`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "client"`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError(ContainSubstring("remote URL is invalid")))
			})
		})

		When("remote-app-port is not set", func() {
			It("should fail", func() {
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
		})

		When("local-app-port is not set", func() {
			It("should default local-app-port to remote-app-port", func() {
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
				Expect(globalConfig.RemoteAppPort).To(Equal(uint16(8080)))
				Expect(globalConfig.LocalAppPort).To(Equal(uint16(8080)))
			})
		})

		When("remote-app-port is not in valid range", func() {
			It("should fail", func() {
				content := []byte(
					`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-url = "ws://localhost:8080"
			remote-app-port = 65589`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("remote app port 65589 is invalid (1-65535)"))
			})
		})

		It("should set local-app-port", func() {
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

		When("local-app-port is not in valid range", func() {
			It("should fail when local-app-port is not in valid range", func() {

				content := []byte(
					`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "client"
			remote-app-port = 8080
			local-app-port = 69989
			remote-url = "ws://localhost:8080"`,
				)

				Expect(loadConfigKoanf(content)).To(MatchError("local app port 69989 is invalid (1-65535)"))
			})
		})
	})

	Describe("dotreplit loader function (server)", func() {
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

		When("listen port is invalid", func() {
			It("default to 0", func() {
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
		})

		When("local-http-port is not in valid range", func() {
			It("should fail", func() {
				content := []byte(
					`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replish]
			mode = "server"
			local-http-port = 65599`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("local http port 65599 is invalid (1-65535)"))
			})
		})

		It("should set local-http-port", func() {
			content := []byte(
				`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replish]
		mode = "server"
		local-http-port = 7777`,
			)
			Expect(loadConfigKoanf(content)).To(Succeed())
			Expect(globalConfig.AppHttpPort).To(Equal(uint16(7777)))

		})

	})
})
