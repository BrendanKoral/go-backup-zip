BINARY_NAME=zip-backup-dbs
.DEFAULT_GOAL := run

build:
	GOARCH=amd64 GOOS=linux go build -o ./bin/${BINARY_NAME}-linux main.go

run:
	go run .