build:
	go build -o bin/vanity-keygen src/main.go

clean:
	go clean
	rm -rf bin/*
