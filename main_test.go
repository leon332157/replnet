package main_test
/*import (
	"testing"

	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/leon332157/replish/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"

	"github.com/valyala/fasthttp"
)

func TestReplish(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Main")
}

var _ = BeforeSuite(func() {
	go startFiber()
	go server.StartReverseProxy()
	time.Sleep(3 * time.Second)
})
var client = &fasthttp.Client{}

func startFiber() {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.SendString("haha")
	})

	go app.Listen("127.0.0.1:7373")
	fmt.Println("fiber started")
}

var _ = Describe("Replish Server Main", func() {
	Describe("TCP Forwarder", func() {
		It("should serve 10000 requests (POST & GET)", func() {
			Expect(makeRequests(10000)).To(Succeed())
		})
	})
})

func makeRequests(n int) error {
	for x := 0; x < n; x++ {
		statusCode, _, err := client.GetTimeout(nil, "http://127.0.0.1:8383", 1000*time.Millisecond)
		if err != nil {
			return fmt.Errorf("Failed on attempt %v err: %v", x, err)
		}
		if statusCode != fasthttp.StatusOK {
			return fmt.Errorf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
		// Assuming GET didn't fail, POST shouldn't fail either.
		// THERE'S NO TIMEOUT FOR POST???
		statusCode, _, err = client.Post(nil, "http://127.0.0.1:8383", nil)
		if err != nil {
			return fmt.Errorf("Failed on attempt %v err: %v", x, err)
		}
		if statusCode != fasthttp.StatusOK {
			return fmt.Errorf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
	}
	return nil
}
*/