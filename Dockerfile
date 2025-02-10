FROM golang:1.22.0 AS build

WORKDIR /app

COPY go.mod /app
COPY go.sum /app
COPY .env /app
COPY .env.development /app
RUN go mod download

COPY . /app

RUN go build -o shortener ./cmd/link-shortener

FROM golang:1.22.0 AS prod

WORKDIR /app
COPY --from=build /app/shortener /app
COPY --from=build /app/.env /app
COPY --from=build /app/.env.development /app

EXPOSE 4000

CMD ["/app/shortener"]
