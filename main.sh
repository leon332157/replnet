go build -ldflags="-s -w" -o bin/replish
REPLISH_USERNAME=test REPLISH_PASSWORD=password ./bin/replish
#GOPATH=/home/runner/go go get github.com/onsi/ginkgo/ginkgo
#go clean -testcache
#/home/runner/go/bin/ginkgo -r 