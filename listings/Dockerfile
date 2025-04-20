FROM golang:1.23.4

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN ls -la /app

RUN go build -o listings-server ./cmd/server/

EXPOSE 3000

CMD ["./listings-server"]