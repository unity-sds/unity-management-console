# Multi-stage build for Unity Management Console

# Backend builder stage
FROM golang:1.21 AS backend-builder
WORKDIR /app
COPY go.mod go.sum ./
COPY backend/ ./backend/
COPY package.json ./
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -buildvcs=false -a -ldflags '-extldflags "-static"' -o management-console ./backend/cmd/web

# Frontend builder stage  
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY package*.json ./
COPY src/ ./src/
COPY static/ ./static/
COPY svelte.config.js vite.config.ts tsconfig.json tailwind.config.js postcss.config.js ./
RUN npm ci
RUN npm run build

# Final runtime stage
FROM debian:bullseye-slim
LABEL maintainer="Unity SDS Team"
LABEL description="Unity Management Console - Containerized deployment"

# Install system dependencies
RUN apt-get update && apt-get install -y \
    ca-certificates \
    tzdata \
    curl \
    wget \
    unzip \
    git \
    bash \
    sqlite3 \
    && rm -rf /var/lib/apt/lists/*

# Install Terraform
ARG TERRAFORM_VERSION=1.5.7
ARG TARGETARCH
RUN ARCH=${TARGETARCH:-amd64} && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_${ARCH}.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_${ARCH}.zip && \
    mv terraform /usr/local/bin/ && \
    chmod +x /usr/local/bin/terraform && \
    rm terraform_${TERRAFORM_VERSION}_linux_${ARCH}.zip

# Install AWS CLI v2
RUN ARCH=${TARGETARCH:-amd64} && \
    if [ "$ARCH" = "amd64" ]; then ARCH_NAME="x86_64"; elif [ "$ARCH" = "arm64" ]; then ARCH_NAME="aarch64"; fi && \
    curl "https://awscli.amazonaws.com/awscli-exe-linux-${ARCH_NAME}.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip aws

# Install Botocore
RUN pip3 install botocore

# Create non-root user
RUN groupadd -g 1001 unity && \
    useradd -u 1001 -g unity -m -s /bin/bash unity

# Create application directories
RUN mkdir -p /app /data/workdir /data/database /data/config && \
    chown -R unity:unity /app /data

# Set working directory
WORKDIR /app

# Copy built applications
COPY --from=backend-builder /app/management-console ./
COPY --from=frontend-builder /app/build ./build

# Set ownership
RUN chown -R unity:unity /app

# Switch to non-root user
USER unity

# Create Unity config directory in container
RUN mkdir -p /home/unity/.unity

# Set environment variables
ENV UNITY_WORKDIR=/data/workdir
ENV UNITY_CONFIG_PATH=/data/config/unity.yaml
ENV PATH="/usr/local/bin:${PATH}"

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# Expose port
EXPOSE 8080

# Default command
CMD ["./management-console", "webapp"]