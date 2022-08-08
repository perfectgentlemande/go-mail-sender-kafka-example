FROM golang:1.18.0 AS builder

ADD . /app
WORKDIR /app
# GOOS/GOARCH as you build not from go alpine
RUN GOOS=linux GOARCH=amd64 go build -o go-kafka-app ./cmd/go-mail-sender-kafka-example

FROM alpine:3.15 AS app
WORKDIR /app
COPY --from=builder /app/go-kafka-app /app
COPY --from=builder /app/cmd/go-mail-sender-kafka-example/config.yaml /app
CMD ["/app/go-kafka-app"]