FROM golang:latest as build

WORKDIR /app

COPY go.mod go.sum /app/
RUN go mod download

COPY . /app

RUN go build -o link-shortener ./cmd/link-shortener

FROM alpine:latest as prod

WORKDIR /app
COPY --from=build /app/link-shortener /app/

EXPOSE 4000

RUN ./link-shortener

CMD ["/app/link-shortener"]
