# Build stage
FROM golang:1.22-alpine AS builder
WORKDIR /app

# 1. Copy go.mod and go.sum from the server folder first (for caching)
COPY server/go.mod server/go.sum* ./
RUN go mod download

# 2. Copy the rest of the server folder's contents into /app
COPY server/ ./

# 3. Build the binary by pointing to the cmd folder where main.go lives
RUN CGO_ENABLED=0 GOOS=linux go build -o server-bin ./cmd

# Run stage
FROM alpine:latest
WORKDIR /app

# 4. Copy the compiled binary from the builder stage
COPY --from=builder /app/server-bin .

# 5. Expose the port Cloud Run will use
ENV PORT=8080
EXPOSE 8080

CMD ["./server-bin"]