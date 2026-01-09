FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o csvproc ./cmd/csvproc

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/csvproc /usr/local/bin/csvproc
COPY data ./data

ENTRYPOINT ["csvproc"]
