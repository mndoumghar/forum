FROM golang:1.20-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o forum ./cmd/main/main.go
FROM alpine:latest
WORKDIR /app
RUN apk add --no-cache tzdata ca-certificates
COPY --from=builder /app/forum .
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./forum"]


