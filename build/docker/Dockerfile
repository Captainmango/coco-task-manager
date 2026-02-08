FROM golang:1.25-alpine AS builder

WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./

RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o cli \
    ./cmd/cli && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w -extldflags '-static'" \
    -trimpath \
    -o server \
    ./cmd/api

FROM alpine:3.22 AS run
WORKDIR /app
COPY --from=builder /app/cli .
COPY --from=builder /app/server .
COPY --from=builder /app/build/docker/start.sh .

RUN chmod +x start.sh && \
    apk add --no-cache cronie

RUN touch /etc/cron.d/root && \
    chmod 644 /etc/cron.d/root

ENV CRONTAB_FILE=/etc/cron.d/root
ENTRYPOINT [ "./start.sh" ]