FROM golang:1.21
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o forum ./cmd/main
EXPOSE 8080
CMD ["./forum"]