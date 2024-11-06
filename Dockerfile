# Stage 1: Build the Go application
FROM golang:1.19-alpine AS builder

# Set the working directory in the container
WORKDIR /app

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project into the container
COPY . .

# Build the application
RUN go build -o backend-coding-challenge-enhanced ./cmd/app

# Stage 2: Prepare the runtime environment
FROM alpine:3.17

# Set the working directory in the runtime container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/backend-coding-challenge-enhanced .

# Expose the application port
EXPOSE 8080

# Define the command to run the application
CMD ["./backend-coding-challenge-enhanced"]
