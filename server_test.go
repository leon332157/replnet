package main_test

import (
	"testing"

	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/leon332157/replish/server"
	. "github.com/onsi/gomega"
	log "github.com/sirupsen/logrus"
	"time"

	goblin "github.com/franela/goblin"
	"github.com/valyala/fasthttp"
)

func TestServer(t *testing.T) {
	g := goblin.Goblin(t)
	RegisterFailHandler(func(m string, _ ...int) { g.Fail(m) })
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	go startFiber()
	go server.StartForwardServer(7373)
	go server.StartReverseProxy()
	time.Sleep(3 * time.Second)
	g.Describe("TCP Forwarder", func() {
		g.It("should serve 10000 requests (POST & GET)", func() {
			Expect(makeRequests(10000, 8383)).To(Succeed())
		})
	})
	g.Describe("Reverse Proxy", func() {
		g.It("should serve 10000 requests (POST & GET)", func() {
			Expect(makeRequests(10000, 8484)).To(Succeed())
		})
	})
}

var client = &fasthttp.Client{}

func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true, DisableKeepalive: false})

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	app.Post("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	go app.Listen("127.0.0.1:7373")
	fmt.Println("fiber started")
}

func makeRequests(n int, port int) error {
	url := fmt.Sprintf("http://127.0.0.1:%v", port)
	var (
		req  fasthttp.Request
		resp fasthttp.Response
	)
	req.SetRequestURI(url)
	for x := 0; x < n; x++ {
		req.Header.SetMethod(fasthttp.MethodGet)
		err := client.DoTimeout(&req, &resp, 1000*time.Millisecond)
		if err != nil {
			return fmt.Errorf("Failed on attempt %v err: %v", x, err)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("Unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
		}
		// Assuming GET didn't fail, POST shouldn't fail either.
		req.Header.SetMethod(fasthttp.MethodPost)
		err = client.DoTimeout(&req, &resp, 1000*time.Millisecond)
		if err != nil {
			return fmt.Errorf("Failed on attempt %v err: %v", x, err)
		}
		if resp.StatusCode() != fasthttp.StatusOK {
			return fmt.Errorf("Unexpected status code: %d. Expecting %d", resp.StatusCode(), fasthttp.StatusOK)
		}
	}
	return nil
}
