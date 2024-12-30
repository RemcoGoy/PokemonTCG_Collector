# Use the official Golang image as the base image
FROM golang:1.23.4-alpine

# Install make
RUN apk add --no-cache make

# Set the working directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the application
RUN make build

# Expose port 3000
EXPOSE 3000

# Command to run the application
CMD ["./main"]
