FROM golang:1.23-alpine as builder

WORKDIR /app

COPY . .
COPY .env .

RUN go mod tidy

RUN go build -o /app/main ./cmd

FROM alpine:latest

COPY --from=builder /app/main /main
COPY ./cmd/config/config.yaml /app/cmd/config/config.yaml

EXPOSE 50051 14268


CMD [ "/main" ]