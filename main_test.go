package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
)

var client = &fasthttp.Client{}

func TestMain(t *testing.T) {
	Convey("Make 10000 requests with no fail", t, func() {
		So(makeGetRequests(10000), ShouldBeNil)
	})
}

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
func BenchmarkMain(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		statusCode, _, err := client.GetTimeout(nil, "http://127.0.0.1:8484", 1000*time.Millisecond)
		if err != nil {
			b.Errorf("Error when loading localhost: %s", err)
		}
		if statusCode != fasthttp.StatusOK {
			b.Errorf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
	}
}
