build:
	@go build -o bin/teenypubsub

run: build
	@./bin/teenypubsub

test:
	@go test ./... -v

benchmark:
	@ go run ./examples/benchmark.go