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

# IPostgreSQL client installation
RUN apk add --no-cache postgresql-client

# Copy binary file from build stage
COPY --from=builder /output/linux/restaurant_system /usr/local/bin/restaurant_system

# Optional:
# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Run app
CMD ["/usr/local/bin/restaurant_system"]
