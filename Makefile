clean:
	rm -rf build

build:
	mkdir -p build
	go build -o build/kubedb ./kubedebugger

fmt:
	@go fmt ./...

golint:
	@docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:latest golangci-lint run -v --config .golangci.yaml