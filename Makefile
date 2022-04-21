all: test medius

test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./...

clean:
	rm -rf bin

medius:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/medius kubevirt.io/containerdisks/cmd/medius

fmt:
	go mod tidy -compat=1.17
	gofmt -s -w .

lint:
	CGO_ENABLED=0 golangci-lint run

.PHONY: cluster-up
cluster-up:
	hack/kubevirtci.sh up

cluster-down:
	hack/kubevirtci.sh down
