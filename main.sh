go build -ldflags="-s -w" -o bin/replnet
if [ $? -ne 0 ];then
    echo "Build failed!"
    exit 1
fi
#python sock.py &
./bin/replnet #--server
#GOPATH=/home/runner/go go get github.com/onsi/ginkgo/ginkgo
#go clean -testcache
#/home/runner/go/bin/ginkgo -r 