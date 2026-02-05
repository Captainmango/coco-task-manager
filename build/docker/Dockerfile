FROM golang:1.24.10-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./

RUN go mod download
COPY . .
RUN go build -o parser /app/cmd/cli
RUN go build -o server /app/cmd/api 


FROM alpine:3.22 AS run
WORKDIR /app
COPY --from=builder /app/parser .
COPY --from=builder /app/server .
COPY --from=builder /app/build/docker/start.sh .
RUN chmod +x start.sh

RUN apk add --no-cache cronie
ENTRYPOINT [ "./start.sh" ]