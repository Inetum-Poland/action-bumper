# Copyright (c) 2024 Inetum Poland.

# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY cmd/ cmd/
COPY internal/ internal/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bumper ./cmd/bumper

# Runtime stage
FROM alpine:3.21

# Install git (required for tagging operations)
RUN apk add --no-cache git

# Add user
RUN addgroup -g 1001 runtime && \
  adduser -D -u 1001 -G runtime runtime

WORKDIR /opt

# Copy the binary from builder
COPY --from=builder /build/bumper /opt/bumper

# set the runtime user to a non-root user and the same user as used by the github runners for actions runs.
USER runtime

# Initial command
ENTRYPOINT ["/opt/bumper"]
