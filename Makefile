.PHONY: build test lint clean fmt vet

BINARY := oak
CMD := ./cmd/oak

build:
	go build -o $(BINARY) $(CMD)

test:
	go test ./...

lint:
	golangci-lint run

fmt:
	gofmt -w .

vet:
	go vet ./...

clean:
	rm -f $(BINARY)
	rm -rf mirror/ image/ dist/

# Run the full pipeline (dry run)
dry-run: build
	./$(BINARY) build --tier 16 --dry-run

# Download and verify Tor Browser only
tor-verify: build
	./$(BINARY) download tor-browser
	./$(BINARY) verify tor-browser
