build:
	go build -o bin/library.exe main.go

test:
	go test ./...

test_s:
	go test -race ./services/...
