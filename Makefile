compose-up-attached:
	docker-compose up --build

compose-up:
	docker-compose up -d --build

compose-down:
	docker-compose down -v

run:
	go run ./cmd/link-shortener

build:
	go build -o shortener ./cmd/link-shortener

build-and-run: build
	./shortener

docker-build-and-run:
	docker build -t shortener .
	docker run -p 8080:80 shortener

run-dynamo:
	docker run -p 8000:8000 amazon/dynamodb-local
