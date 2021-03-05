build:
	go build -o bin/client
	./bin/client

clean:
	rm -rf bin/client