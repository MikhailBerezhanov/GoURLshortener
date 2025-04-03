# --- Builder Stage ---
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o myapp

# --- Deploy Stage ---
FROM mongo

# Set the working directory inside the container
WORKDIR /app

# Copy the built application from the builder stage
COPY --from=builder /app/myapp /app/myapp

# Grant execution permissions
RUN chmod +x /app/myapp

# Expose MongoDB's default port
EXPOSE 27017

# Run MongoDB and the Go application
CMD ["bash", "-c", "mongod --bind_ip_all & sleep 5 && /app/myapp"]