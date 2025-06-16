FROM golang:1.23

RUN apt-get update && apt-get install -y libsqlite3-dev

WORKDIR /forum

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o forum ./cmd/main/main.go

EXPOSE 8080

CMD ["./forum"]


# FROM golang:1.23 AS builder
# WORKDIR /app
# COPY go.mod ./
# RUN go mod download
# COPY . .
# RUN go build -o forum ./cmd/main/main.go
# FROM alpine:latest
# WORKDIR /app
# RUN apk add --no-cache tzdata ca-certificates
# COPY --from=builder /app/forum .
# COPY --from=builder /app/templates ./templates
# COPY --from=builder /app/static ./static
# EXPOSE 8080
# CMD ["./forum"]


