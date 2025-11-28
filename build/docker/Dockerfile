FROM golang:1.24.10-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -o parser /app/cmd/cli


FROM alpine:3.22 AS run
WORKDIR /app
COPY --from=builder /app/parser .

RUN apk add --no-cache cronie
ENTRYPOINT [ "./parser" ]