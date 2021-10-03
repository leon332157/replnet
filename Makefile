build:
	go build -ldflags="-s -w" -o bin/replish
	./bin/replish

clean:
	rm -rf bin/replish