load-env:
	source .env

run: load-env
	go run cmd/main.go

build:
	go build -o bin/cryptodashboard cmd/server/main.go

run-build: load-env
	./bin/cryptodashboard
