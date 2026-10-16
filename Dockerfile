# Stage 1: Build binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy Go module files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build statically linked Go executable
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o psx-server cmd/server/main.go

# Stage 2: Minimal runtime image (~15MB)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/psx-server .

EXPOSE 8080

ENV PORT=8080
ENV ENVIRONMENT=production

CMD ["./psx-server"]
