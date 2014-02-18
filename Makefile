.PHONY: fmt install lint

clean:
	go clean ./...
fmt:
	gofmt -l -w -tabs=true $(GOPATH)/src/github.com/zimmski/container
install:
	go install ./...
lint: clean
	go tool vet -all=true -v=true $(GOPATH)/src/github.com/zimmski/container
	golint $(GOPATH)/src/github.com/zimmski/container
test: clean
	go test ./...
testcover: clean
	go test -cover ./...

