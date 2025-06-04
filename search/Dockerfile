FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o search-service ./cmd/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /app/search-service .
EXPOSE 8080

CMD ["./search-service"]
