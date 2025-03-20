ARG GO_VERSION=1.24.1

FROM golang:${GO_VERSION}-alpine AS builder

WORKDIR /app

COPY go.mod ./

# COPY go.sum ./
# RUN go mod download

COPY . .

RUN go build -o /orbit-api cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /orbit-api .

EXPOSE 8080

CMD ["./orbit-api"]