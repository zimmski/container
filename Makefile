.PHONY: fmt install lint
fmt:
	gofmt -l -w -tabs=true $(GOPATH)/src/github.com/zimmski/container
install:
	go install github.com/zimmski/container/list
	go install github.com/zimmski/container/list/doublylinkedlist
	go install github.com/zimmski/container/list/linkedlist
	go install github.com/zimmski/container/list/unrolledlinkedlist
lint: install
	go tool vet -all=true -v=true $(GOPATH)/src/github.com/zimmski/container
	golint $(GOPATH)/src/github.com/zimmski/container
