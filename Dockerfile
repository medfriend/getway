# Use the official Golang image to create a build artifact.
FROM golang:1.23.3 as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container, incluyendo el .env
COPY . .

# ✅ Copiar el archivo .env a la misma ubicación donde está main.go
COPY .env /app/.env

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# Start a new stage from scratch
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/myapp .

# ✅ Copiar el .env en la misma ubicación donde se ejecutará `myapp`
COPY --from=builder /app/.env .

# ✅ Definir variable de entorno para que la aplicación sepa dónde encontrarlo
ENV ENV_PATH=/root/.env

# Expose port 8070 to the outside world
EXPOSE 8070

# Command to run the executable
CMD ["./myapp"]
