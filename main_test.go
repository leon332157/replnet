package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func BenchmarkMain(b *testing.B) {
	c := &fasthttp.Client{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		statusCode, _, err := c.GetTimeout(nil, "http://127.0.0.1:8383", 1000*time.Millisecond)
		if err != nil {
			b.Errorf("Error when loading localhost: %s", err)
		}
		if statusCode != fasthttp.StatusOK {
			b.Errorf("Unexpected status code: %d. Expecting %d", statusCode, fasthttp.StatusOK)
		}
	}
}
