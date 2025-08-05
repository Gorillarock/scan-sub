test:
	go clean -testcache && go test -v ./...

build:
	docker compose build

run:
	docker compose up

run-with-test-db-reader:
	docker compose --profile test up