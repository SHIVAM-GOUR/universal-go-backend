# ─────────────────────────────────────────────────────────────────────────────
# Stage 1: Builder
# Generates Swagger docs then compiles a statically linked binary.
# ─────────────────────────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install swag CLI for Swagger doc generation.
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Cache dependency layer: copy module files first so Docker can reuse
# the layer on source-only changes.
COPY go.mod go.sum ./
RUN go mod download

# Copy all source.
COPY . .

# Generate Swagger docs from annotations before compiling.
RUN swag init -g cmd/api/main.go -o docs/

# Compile a statically linked, stripped binary.
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/api \
    ./cmd/api/main.go

# ─────────────────────────────────────────────────────────────────────────────
# Stage 2: Final image
# Copies only the binary; final image is well under 20 MB.
# ─────────────────────────────────────────────────────────────────────────────
FROM alpine:3.19

WORKDIR /app

# CA certificates (for outbound HTTPS calls) and timezone data.
RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/bin/api .

EXPOSE 8080

CMD ["/app/api"]
