# Stage 1: Build the Go app
FROM golang:1.24-alpine3.21 AS builder

WORKDIR /app
COPY . .

RUN apk update && apk add git gcc musl-dev zsh

RUN go mod download
# CGO required for sqlite driver
RUN CGO_ENABLED=1 go build -o /app/main ./cmd/

# Stage 2: Run the Go app
FROM alpine:3.21

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Set gin mode
ENV GIN_MODE=release
# Expose port 8080
EXPOSE 8080

# Command to run the app
CMD ["./main"]
