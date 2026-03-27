# Stage 1: Build Astro frontend
FROM node:22-alpine AS web-builder
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.26-alpine AS go-builder
RUN apk add --no-cache git
WORKDIR /app
COPY go.mod ./
COPY go.sum* ./
RUN go mod download
COPY . .
COPY --from=web-builder /app/web/dist ./cmd/server/static
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.version=${VERSION:-dev} -X main.commitSHA=$(git rev-parse --short HEAD 2>/dev/null || echo unknown) -X main.buildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o /poke-store ./cmd/server

# Stage 3: Minimal runtime
FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata && \
    adduser -D -u 1000 appuser
COPY --from=go-builder /poke-store /usr/local/bin/poke-store
USER appuser
EXPOSE 6001
ENV ADDR=:6001
ENTRYPOINT ["poke-store"]
