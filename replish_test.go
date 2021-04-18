package main

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
    "fmt"
    "github.com/valyala/fasthttp"
    "time"
	//"github.com/leon332157/replish/server"
    fiber "github.com/gofiber/fiber/v2"
)

var client = &fasthttp.Client{}
var _ = BeforeSuite(func() {
    go startFiber()
    time.Sleep(3*time.Second)
})

func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	go app.Listen("127.0.0.1:7373")
	fmt.Println("fiber started")
}

func TestReplish(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Replish Main")
}

var _ = Describe("Replish Main", func() {
    It("can make 10000 requests with no error", func() {
        Expect(makeGetRequests(10000)).To(BeNil())
    })
})

func makeGetRequests(n int) error {
	for x := 0; x < n; x++ {
		statusCode, _, err := client.GetTimeout(nil, "http://127.0.0.1:8383", 1000*time.Millisecond)
		if err != nil {
			return fmt.Errorf("Failed on attempt %v err: %v", x, err)
		}
		if statusCode != fasthttp.StatusOK {
			return fmt.Errorf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
	}
	return nil
}