build-alpha-client:
	go build -o bin/alpha-client cmd/alpha-client/main.go

build-alpha-server:
	go build -o bin/alpha-server cmd/alpha-server/main.go

build: build-alpha-client build-alpha-server

compile:
	echo "Compiling alpha-client for multiple platforms"
	GOOS=linux GOARCH=amd64 go build -o bin/alpha-client-linux-amd64 cmd/alpha-client/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/alpha-client-darwin-amd64 cmd/alpha-client/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/alpha-client-windows-amd64 cmd/alpha-client/main.go
	echo "Compiling alpha-server for multiple platforms"
	GOOS=linux GOARCH=amd64 go build -o bin/alpha-server-linux-amd64 cmd/alpha-server/main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/alpha-server-darwin-amd64 cmd/alpha-server/main.go
	GOOS=windows GOARCH=amd64 go build -o bin/alpha-server-windows-amd64 cmd/alpha-server/main.go

run-alpha-client:
	go run cmd/alpha-client/main.go

run-alpha-server:
	go run cmd/alpha-server/main.go
