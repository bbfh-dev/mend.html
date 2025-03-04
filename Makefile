VERSION := v0.1.0
BINARY_NAME := app

run:
	@go run -ldflags "-X main.Version=$(VERSION)" main.go $(ARGS)

test:
	@go test -v ./...

coverage:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o ./tests/coverage.html
	@rm coverage.out

build-linux-x64:
	GOOS=linux GOARCH=amd64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-linux-x64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-linux-arm64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-windows-x64:
	GOOS=windows GOARCH=amd64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-windows-x64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-windows-arm:
	GOOS=windows GOARCH=arm64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-windows-arm64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-darwin-x64:
	GOOS=darwin GOARCH=amd64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-darwin-x64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-darwin-arm:
	GOOS=darwin GOARCH=arm64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-freebsd-x64:
	GOOS=freebsd GOARCH=amd64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-freebsd-x64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-freebsd-arm:
	GOOS=freebsd GOARCH=arm64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-freebsd-arm64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-openbsd-x64:
	GOOS=openbsd GOARCH=amd64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-openbsd-x64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-openbsd-arm:
	GOOS=openbsd GOARCH=arm64 go build -o dist/$(BINARY_NAME) -ldflags "-s -w -X main.Version=$(VERSION)" main.go
	@tar -czvf dist/$(BINARY_NAME)-$(VERSION)-openbsd-arm64.tar.gz -C dist $(BINARY_NAME)
	@rm dist/$(BINARY_NAME)

build-all: build-linux-x64 build-linux-arm build-windows-x64 build-windows-arm build-darwin-x64 build-darwin-arm build-freebsd-x64 build-freebsd-arm build-openbsd-x64 build-openbsd-arm
