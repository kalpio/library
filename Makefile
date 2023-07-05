build:
    ifeq ($(OS),Windows_NT)
	go build -o bin/library.exe main.go
    else
	go build -o bin/library main.go
    endif

test:
	go test ./...

test_race:
	go test -race ./services/...
