# Stage 1, use Goland alpine builder
FROM golang:alpine as builder

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# COPY ./Build the package/ .
RUN go clean && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o gripo-cron-linux-amd64

# Build a small image
FROM alpine

# Copy the build pakcage to "gripo-cron" folder
COPY --from=builder /build/gripo-cron-linux-amd64 /gripo-cron/gripo-cron-linux-amd64

# Set to working directory
WORKDIR /gripo-cron/

# Expose the backend port
EXPOSE 8080

# Docker will run backend Prodgram when we start docker
ENTRYPOINT ["./gripo-cron-linux-amd64"]