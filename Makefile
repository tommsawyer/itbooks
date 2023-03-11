.PHONY: build test lint scrape publish postgres migrate down

build:
	@ go build -o ./build/itbooks ./cmd/itbooks
	
test:
	@ go test ./... -race

lint:
	@ golangci-lint run --config .golangci.yml

scrape: build
	@ ./build/itbooks scrape

publish: build
	@ ./build/itbooks publish

postgres:
	@ docker-compose up -d postgres

migrate:
	@ docker-compose run migrate

down:
	@ docker-compose down

psql:
	@ docker-compose exec postgres psql -U itbooks
