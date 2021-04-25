git config --local user.name leon332157
git config --local user.email leon332157@gmail.com
GOPATH=/home/runner/go go get github.com/onsi/ginkgo/ginkgo
export PATH=$PATH:/home/runner/go/bin
echo "Bootstrapped"