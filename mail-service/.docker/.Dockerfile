# base go image
FROM golang:1.25.2-alpine AS builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o mailApp ./cmd/api

RUN chmod +x /app/mailApp

# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/mailApp .
COPY --from=builder /app/templates ./templates

CMD [ "/app/mailApp" ]