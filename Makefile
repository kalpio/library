test:
	go test ./...

test_s:
	go test -race ./services/...
