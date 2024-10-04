FROM golang:latest as build

WORKDIR /app

# Copy the Go module files
COPY go.mod /app
COPY go.sum /app

# Download the Go module dependencies
RUN go mod download

COPY . /app

RUN go build -o /app/link-shortener
 
FROM alpine:latest as run

# Copy the application executable from the build image
COPY --from=build /app/link-shortener /app/link-shortener

WORKDIR /app
EXPOSE 7001

ENV DOMAIN_NAME=https://go.illiashenko.dev

CMD ["/link-shortener"]
