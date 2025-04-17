FROM golang:1.23.7-alpine3.20 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o user-service ./cmd

FROM alpine:3.18

WORKDIR /app/

COPY --from=builder /app/user-service .

CMD ["./user-service"]