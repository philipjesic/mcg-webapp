FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

RUN go build -o bids-service ./cmd/server/

EXPOSE 3000

CMD ["./bids-service"]