## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

all: test medius

test: lint
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./...

clean:
	rm -rf bin

medius:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/medius kubevirt.io/containerdisks/cmd/medius

fmt:
	go mod tidy -compat=1.19
	gofmt -s -w .

.PHONY: vendor
vendor:
	go mod vendor

GOLANGCI_LINT ?= $(LOCALBIN)/golangci-lint
GOLANGCI_LINT_VERSION ?= v1.51.1

.PHONY: lint
lint:
	test -s $(GOLANGCI_LINT) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(LOCALBIN) $(GOLANGCI_LINT_VERSION)
	CGO_ENABLED=0 $(GOLANGCI_LINT) run --timeout 5m

.PHONY: cluster-up
cluster-up:
	hack/kubevirtci.sh up

cluster-down:
	hack/kubevirtci.sh down
