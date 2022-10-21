package main

import (
	"testing"

	"github.com/leon332157/replnet/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
)

func UNUSED(x ...interface{}) {
}
func TestReplnet(t *testing.T) {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(true)
	log.SetLevel(log.ErrorLevel)
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

var _ = Describe(".replit file loader", func() {
	BeforeEach(func() {
		globalConfig = common.ReplnetConfig{}
		Expect(globalConfig).To(Equal(common.ReplnetConfig{})) // golang complains so i need a sanity check
		// clear the global config struct before each test
	})

	AfterEach(func() {
		for key := range koanf.Raw() {
			koanf.Delete(key) // clear the koanf map for each test
		}
	})

	When("toml is invalid", func() {
		It("should fail", func() {
			content := []byte(`broken`)
			Expect(loadConfigKoanf(content)).ToNot(Succeed())
		})
	})

	When("replnet field isn't present", func() {
		It("should fail when replnet tag isn't present", func() {
			content := []byte(
				`language = "go"
run = "bash main.sh"
onBoot="bash bootstrap.sh"`)
			Expect(loadConfigKoanf(content)).To(MatchError("replnet field doesn't exist"))
		})
	})

	When("mode is not set", func() {
		It("should return error", func() {
			content := []byte(`
			language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replnet]
			`)
			Expect(loadConfigKoanf(content)).To(MatchError("mode is required"))
		})
	})

	When("mode is invalid", func() {
		It("should fail", func() {
			content := []byte(`
	 language = "go"
	 run = "bash main.sh"
	 onBoot="bash bootstrap.sh"
	 [replnet]
	 mode="test"`)
			Expect(loadConfigKoanf(content)).To(MatchError("mode test is invalid"))
		})
	})

	Describe("config loader (client)", func() {

		When("correct sample config is given", func() {
			It("should succeed", func() {
				content := []byte(
					`[replnet]
				mode="client"
				[replnet.client]
				remote-url="ws://test"
				remote-port = 8080
				local-port = 8081`)
				Expect(loadConfigKoanf(content)).To(Succeed())
				Expect(globalConfig.Mode).To(Equal("client"))
				Expect(globalConfig.Client.RemoteURL).To(Equal("ws://test"))
				Expect(globalConfig.Client.RemotePort).To(Equal(uint16(8080)))
				Expect(globalConfig.Client.LocalPort).To(Equal(uint16(8081)))
			})
		})

		When("client field is not set", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("replnet.client field doesn't exist"))
			})
		})

		When("remote-url is not set", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"
		[replnet.client]`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-url is required"))
			})
		})

		When("remote-url is not parseable", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"
		[replnet.client]
		remote-url = "test"`)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-url is invalid: parse \"test\": invalid URI for request"))
			})
		})

		When("remote-url is not the correct schema", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"
		[replnet.client]
		remote-url = "ftp://localhost:8080"`)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-url must be http, https, ws or wss"))
			})
		})

		When("remote-url host is empty", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"
		[replnet.client]
		remote-url = "http://"`)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-url is invalid: host is empty"))
			})
		})

		When("remote-port is not set", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
		mode = "client"
		[replnet.client]
		remote-url = "ws://localhost:8080"`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-port is required"))
			})
		})

		When("remote-port is not in valid range", func() {
			It("should fail", func() {
				content := []byte(
					`[replnet]
			mode = "client"
			[replnet.client]
			remote-url = "ws://localhost:8080"
			remote-port = 65589`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("remote-port 65589 is invalid (1-65535)"))
			})
		})

		When("local-port is not set", func() {
			It("should default local-port to remote-port", func() {
				content := []byte(
					`[replnet]
			mode = "client"
			[replnet.client]
			remote-port = 8080
			remote-url = "ws://localhost:8080"`,
				)
				Expect(loadConfigKoanf(content)).To(Succeed())
				Expect(globalConfig.Client.RemotePort).To(Equal(uint16(8080)))
				Expect(globalConfig.Client.LocalPort).To(Equal(uint16(8080)))
			})
		})

		When("local-port is not in valid range", func() {
			It("should fail", func() {

				content := []byte(
					`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replnet]
			mode = "client"
			[replnet.client]
			remote-port = 8080
			local-port = 69989
			remote-url = "ws://localhost:8080"`,
				)

				Expect(loadConfigKoanf(content)).To(MatchError("local-port 69989 is invalid (1-65535)"))
			})
		})
	})

	Describe("dotreplit loader function (server)", func() {

		It("default listen-port and reverse-proxy-port to 0 if field is not present", func() {
			content := []byte(
				`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replnet]
		mode = "server"`,
			)
			Expect(loadConfigKoanf(content)).To(Succeed())
			Expect(globalConfig.Server.ListenPort).To(Equal(uint16(0)))
			Expect(globalConfig.Server.ReverseProxyPort).To(Equal(uint16(0)))
		})

		When("listen port or reverse-proxy-port is 0", func() {
			It("should succeed", func() {
				content := []byte(
					`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replnet]
		mode = "server"
		[replnet.server]
		listen-port = 0
		reverse-proxy-port = 0`,
				)
				Expect(loadConfigKoanf(content)).To(Succeed())
				Expect(globalConfig.Server.ListenPort).To(Equal(uint16(0)))
				Expect(globalConfig.Server.ReverseProxyPort).To(Equal(uint16(0)))
			})
		})

		When("listen-port not in valid range", func() {
			It("default to 0", func() {
				content := []byte(
					`language = "go"
		run = "bash main.sh"
		onBoot="bash bootstrap.sh"
		[replnet]
		mode = "server"
		[replnet.server]
		listen-port = 99999`,
				)
				Expect(loadConfigKoanf(content)).To(Succeed())
				Expect(globalConfig.Server.ListenPort).To(Equal(uint16(0)))
			})
		})

		When("reverse-proxy-port is not in valid range", func() {
			It("should fail", func() {
				content := []byte(
					`language = "go"
			run = "bash main.sh"
			onBoot="bash bootstrap.sh"
			[replnet]
			mode = "server"
			[replnet.server]
			reverse-proxy-port = 65599`,
				)
				Expect(loadConfigKoanf(content)).To(MatchError("reverse-proxy-port 65599 is invalid (0-65535)"))
			})
		})

	})
})
