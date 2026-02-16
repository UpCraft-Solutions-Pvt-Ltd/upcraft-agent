# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY main.go qdrant_client.go ./
RUN go build -o server main.go qdrant_client.go

# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
