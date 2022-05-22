MODULES = $(shell find . -name go.mod | sed s:/go.mod::g)

.PHONY: test
test:
	go test $(MODULES)

.PHONY: vet
vet:
	go vet $(MODULES)
