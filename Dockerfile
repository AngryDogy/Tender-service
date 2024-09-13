FROM golang:1.23

COPY . /app
WORKDIR /app

RUN go mod download

RUN go build "cmd/tenderservice/main.go"

EXPOSE 8080

ENTRYPOINT ["./main"]