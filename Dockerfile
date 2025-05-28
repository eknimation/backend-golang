# ========== BUILD STAGE ==========
FROM golang:1.24.3-alpine AS build

WORKDIR /usr/src/app

ARG BUILD_TARGET=api

# Install necessary dependencies for building
RUN apk update && apk add --no-cache gcc musl-dev bash tzdata

# Copy Go module files and download dependencies
COPY go.mod .
COPY go.sum .
RUN go mod download

# Install swag for generating Swagger documentation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Enforce UTC for all operations
ENV TZ=UTC

# Copy the source code
COPY . .

# Generate Swagger documentation
RUN swag init -g cmd/api/main.go -o docs

# Build the Go binary
RUN go build -ldflags="-s -w" -o "bin/main" "./cmd/${BUILD_TARGET}"

# ========== FINAL SHIPPING STAGE ==========
FROM alpine:latest

# Install tzdata for timezone support and create non-root user
RUN apk add --no-cache tzdata ca-certificates wget && \
    addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

WORKDIR /usr/src/app

# Enforce UTC for all operations
ENV TZ=UTC

# Define the ENVIRONMENT variable (default to "development")
ARG ENVIRONMENT=development
ENV ENVIRONMENT=${ENVIRONMENT}

# Copy the built binary from the build stage
COPY --from=build /usr/src/app/bin/main .

# Copy the environment-specific .env file
#COPY --from=build /usr/src/app/envs/.env.example /usr/src/app/.env

# Change ownership to non-root user
RUN chown -R appuser:appgroup /usr/src/app

# Switch to non-root user
USER appuser

# Expose the port (assuming your app runs on 5555 based on main.go)
EXPOSE 5555

# Set the entrypoint to the built binary
CMD ["./main"]