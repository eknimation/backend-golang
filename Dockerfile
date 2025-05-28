# ========== BUILD STAGE ==========
FROM golang:1.24.2-alpine AS build

WORKDIR /usr/src/app

ARG BUILD_TARGET=api

# Install necessary dependencies for building
RUN apk update && apk add --no-cache gcc musl-dev bash tzdata

# Copy Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Enforce UTC for all operations
ENV TZ=UTC

# Copy the source code
COPY . .

# Build the Go binary
RUN go build -ldflags="-s -w" -o "bin/main" "./cmd/${BUILD_TARGET}"

# ========== FINAL SHIPPING STAGE ==========
FROM alpine:latest

WORKDIR /usr/src/app

# Install tzdata for timezone support
RUN apk add --no-cache tzdata

# Enforce UTC for all operations
ENV TZ=UTC

# Define the ENVIRONMENT variable (default to "development")
ARG ENVIRONMENT=development
ENV ENVIRONMENT=${ENVIRONMENT}

# Copy the built binary from the build stage
COPY --from=build /usr/src/app/bin/main .

# Copy the environment-specific .env file
COPY --from=build /usr/src/app/envs/.env.${ENVIRONMENT} /usr/src/app/.env

# Set the entrypoint to the built binary
CMD ["./main"]