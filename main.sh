go build -ldflags="-s -w" -o bin/replish
if [ $? -ne 0 ];then
    echo "Build failed!"
    exit 1
fi
python sock.py &
./bin/replish --server
#GOPATH=/home/runner/go go get github.com/onsi/ginkgo/ginkgo
#go clean -testcache
#/home/runner/go/bin/ginkgo -r 