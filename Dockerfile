
FROM golang:1.23 AS builder
RUN apt-get update && apt-get install -y gcc libc6-dev sqlite3 libsqlite3-dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64
RUN go build -o app ./cmd/api


FROM debian:bookworm
RUN apt-get update && apt-get install -y \
  ca-certificates \
  sqlite3 \
  libsqlite3-0 \
  && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/app .
COPY --from=builder /app/cmd/migrations ./cmd/migrations
COPY --from=builder /app/db/gdrive-db ./db/gdrive-db
RUN mkdir -p /app/uploads \
    && chown -R 755 /app

RUN adduser --disabled-password --gecos "" app
USER app

EXPOSE 8000
CMD ["./app"]
    