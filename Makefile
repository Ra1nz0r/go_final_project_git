run_server:
	go run ./...
run_test:
	go test -v ./... -count=1

.PHONY: run_serv run_test run_test1