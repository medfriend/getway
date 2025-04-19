# Use the official Golang image to create a build artifact.
FROM golang:1.23.3 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container, incluyendo el .env
COPY *.go ./

# Copy specific directories with their contents
COPY dto/ ./dto/
COPY consulRegister/ ./consulRegister/
COPY httpServer/ ./httpServer/
COPY jwt/ ./jwt/
COPY util/ ./util/

# ✅ Copiar el archivo .env a la misma ubicación donde está main.go
COPY .env.devprod /app/.env

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# In the Alpine stage, after installing ca-certificates:
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create app directory and set ownership
WORKDIR /app

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# Copy the .env file
COPY --from=builder /app/.env .

# Set proper permissions
RUN chown -R appuser:appgroup /app

# Set environment variable
ENV ENV_PATH=/app/.env

# Switch to non-root user
USER appuser

# Expose port 8070 to the outside world
EXPOSE 8070

# Command to run the executable
CMD ["./myapp"]