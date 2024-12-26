# Start by pulling the official Go image for building the application
FROM golang:1.23 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy all files from the backend folder to the container
COPY ./ .

# Download dependencies and build the Go binary
RUN go mod download
RUN go build -o app

# Use Ubuntu as the base image
FROM ubuntu:22.04

# Update packages and install required libraries
RUN apt-get update && apt-get install -y \
    libc6 \
    wget \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/app .

# Copy the .env file into the container
COPY .env /app/.env

# Expose the application port
EXPOSE 8080

# Load environment variables from .env file
ENV $(cat /app/.env | xargs)

# Command to run the application
CMD ["/app/app"]
