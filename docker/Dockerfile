FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/server ./cmd

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/bin/server /app/server

EXPOSE 8080
ENV HTTP_ADDR=:8080

CMD ["/app/server"]
