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

run:
	go run main.go

run_css:
	npx tailwindcss -i ./templates/assets/src/input.css -o ./templates/assets/styles/output.css --watch
