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


