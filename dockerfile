FROM golang:1.22.2-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o kafka-consumer

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/kafka-consumer .
COPY ./config/config.yaml /root/config/config.yaml
COPY ./config/dbinit.sql /root/config/dbinit.sql

EXPOSE 8080

CMD ["./kafka-consumer"]