# syntax=docker/dockerfile:1

FROM golang:1.23 AS builder

# Set destination for COPY
WORKDIR /usr/src/app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /restaurant_system

# Compilation for Linux
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /output/linux/restaurant_system

# Compilation for Windows
# RUN CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o /output/windows/restaurant_system.exe

# Compilation for macOS
# RUN CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o /output/macos/restaurant_system

# Final image based on alpine
FROM alpine:latest AS final

# PostgreSQL client installation
RUN apk add --no-cache postgresql-client

# Copy binary file from build stage
COPY --from=builder /output/linux/restaurant_system /usr/local/bin/restaurant_system

# Copy the HTML file to the app directory in the container
COPY index.html /usr/src/app/index.html

# Set the working directory
WORKDIR /usr/src/app

# Expose port 8080
EXPOSE 8080

# Run app
CMD ["/usr/local/bin/restaurant_system"]
