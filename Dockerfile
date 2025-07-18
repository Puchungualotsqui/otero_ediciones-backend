FROM golang:1.24-alpine

ENV CGO_ENABLED=0 \
    GO111MODULE=on

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN go build -o main .

CMD ["./main"]
