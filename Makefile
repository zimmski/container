.PHONY: fmt lint
fmt:
	gofmt -l -w -tabs=true $(GOPATH)/src/github.com/zimmski/container
lint:
	go tool vet -all=true -v=true $(GOPATH)/src/github.com/zimmski/container
	golint $(GOPATH)/src/github.com/zimmski/container
