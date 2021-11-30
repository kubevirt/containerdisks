all: test medius

test:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go test ./...

clean:
	rm -rf bin

medius:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/medius kubevirt.io/containerdisks/cmd/medius

fmt:
	go mod tidy
	gofmt -s -w .
