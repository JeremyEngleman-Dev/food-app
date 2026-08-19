# Build stage
# ─────────────────────────────────────────────
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/server \
    ./cmd/api

# Runtime stage
# ─────────────────────────────────────────────
FROM alpine:latest

WORKDIR /app

RUN addgroup -S foodgroup && \
    adduser -S foodie -G foodgroup

COPY --from=builder /app/server .
COPY --from=builder /app/internal/platform/database/migrations ./internal/platform/database/migrations

USER foodie

EXPOSE 8080

CMD ["./server"]