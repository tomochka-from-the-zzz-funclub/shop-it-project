FROM golang:1.23-alpine3.20 as builder

WORKDIR /src

RUN apk add build-base
COPY go.mod go.sum ./
RUN go mod tidy
RUN go mod download

COPY . ./

RUN GOOS=linux go build -o app ./cmd

FROM alpine:3.20 as runner

WORKDIR /app

COPY --from=builder /src/app ./

CMD ["./app"]