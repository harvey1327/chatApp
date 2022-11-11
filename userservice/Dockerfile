FROM golang:1.19.1-alpine3.16 AS builder

# Set environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -o main ./cmd/main.go

# Build a small image
FROM scratch

COPY --from=builder /build/main /

# Command to run
ENTRYPOINT ["/main"]