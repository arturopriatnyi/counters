FROM golang:1.19-alpine

WORKDIR ./counters
COPY . .

RUN go build -o ./build/counters ./cmd/counters/main.go
CMD ["./build/counters"]
