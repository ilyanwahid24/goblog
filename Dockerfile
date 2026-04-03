FROM docker.io/library/golang:alpine AS builder

WORKDIR /app

# Copy dependency configs and install
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o goblog .

# Use a minimal alpine image for the runtime stage
FROM docker.io/library/alpine:latest
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary and required assets
COPY --from=builder /app/goblog .
COPY templates/ templates/
COPY static/ static/
COPY migrations/ migrations/

EXPOSE 8080

# Command to run the executable
CMD ["./goblog"]
