compose-up-attached:
	docker-compose up --build

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down -v

run:
	APP_ENV=development go run ./cmd/link-shortener

build:
	go build -o shortener ./cmd/link-shortener

build-and-run: build
	APP_ENV=development ./shortener

run-dynamo:
	docker run -p 8000:8000 amazon/dynamodb-local
