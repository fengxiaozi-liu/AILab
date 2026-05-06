BINARY ?= ferryPilot
DIST_DIR ?= dist

.PHONY: test build clean

test:
	go test ./...

build: test
	go build -o bin/$(BINARY) ./cmd/ferryPilot

clean:
	rm -rf bin $(DIST_DIR)
