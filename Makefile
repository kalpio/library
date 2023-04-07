build:
	go mod tidy
    ifeq ($(OS), Windows_NT)
		set GOOS=windows
		set GOARCH=amd64
		go build -v -x -o bin/library.exe main.go
    else
		env GOOS=linux GOARCH=amd64 go build -v -x -o bin/library main.go
    endif

build_race:
	go mod tidy
    ifeq ($(OS), Windows_NT)
		set GOOS=windows
		set GOARCH=amd64
		go build -v -x -race -o bin/library.exe main.go
    else
		env GOOS=linux GOARCH=amd64 go build -v -x -race -o bin/library main.go
    endif

release:
	go mod tidy
    ifeq ($(OS), Windows_NT)
		set GOOS=windows
		set GOARCH=amd64
		go build -v -x -ldflags "-s -w" -o bin/release/library.exe main.go
    else
		env GOOS=linux GOARCH=amd64 go build -v -x -ldflags "-s -w" -o bin/release/library main.go
    endif

test:
	go test -race ./...
