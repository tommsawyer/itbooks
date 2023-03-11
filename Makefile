.PHONY: build test lint scrape publish

build:
	@ go build -o ./build/itbooks ./cmd/itbooks
	
test:
	@ go test ./... -race

lint:
	@ golangci-lint run

scrape: build
	@ ./build/itbooks scrape

publish: build
	@ ./build/itbooks publish
