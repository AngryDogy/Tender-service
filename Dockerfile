FROM golang:1.21

COPY . /app
WORKDIR /app

RUN go mod download

RUN go build "cmd/EWallet/main.go"

CMD ["./main"]