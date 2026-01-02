# Build stage
FROM golang:alpine AS builder

# Install git required for fetching dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# CGO_ENABLED=0 is important for alpine to ensure static linking
RUN CGO_ENABLED=0 GOOS=linux go build -o fin-news main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the Pre-built binary from the previous stage
COPY --from=builder /app/fin-news .

# Expose port and start
EXPOSE 3001

CMD ["./fin-news"]
