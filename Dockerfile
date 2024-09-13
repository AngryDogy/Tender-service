FROM golang:1.23

COPY . /app
WORKDIR /app

RUN go mod download

RUN go build "cmd/EWallet/main.go"

CMD ["./main"]