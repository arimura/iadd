all:
	go build -o iadd cmd/iadd/main.go

run:
	go run cmd/iadd/main.go

test:
	go test