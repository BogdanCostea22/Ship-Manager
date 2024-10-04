# Start from the official Go image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project
COPY . .

# Download all dependencies
RUN go mod download

# Build the application
# Adjust this path to where your main.go file is located
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Start a new stage from scratch
FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the pre-built binary file from the previous stage
COPY --from=builder /app/main .

# Copy necessary files and directories
COPY --from=builder /app/cmd ./cmd
COPY --from=builder /app/internal ./internal

# Expose port 8080 to the outside world
EXPOSE 8080
ENV PORT=8080
# Command to run the executable
CMD ["./main"]