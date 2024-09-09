FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN go build -o bin/main ./app/.

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/bin/main .
CMD ["/app/main"]
