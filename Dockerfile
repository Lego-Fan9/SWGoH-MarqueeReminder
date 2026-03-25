# Stage 1: Build
FROM golang:1.24 AS builder
WORKDIR /app

ENV CGO_ENABLED=0

COPY src/go.mod src/go.sum ./

RUN go mod download

COPY ./src/ ./

RUN go build -o MarqueeReminder -v -ldflags="-s -w" -tags netgo .

# Stage 2: Certs
FROM debian:bookworm-slim AS certs
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates

# Stage 3: Runtime
FROM scratch
WORKDIR /app

COPY --from=builder /app/MarqueeReminder ./

COPY --from=certs /etc/ssl/certs /etc/ssl/certs

ENV DOCKER=TRUE

ENTRYPOINT ["./MarqueeReminder"]