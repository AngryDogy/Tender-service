FROM golang:1.23

COPY . /app
WORKDIR /app

RUN go mod download

EXPOSE 8080

RUN go run "cmd/tenderservice/main.go"

