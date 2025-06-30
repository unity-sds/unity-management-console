# Multi-stage build for Unity Management Console

# Backend builder stage
FROM golang:1.21-alpine AS backend-builder
WORKDIR /app
COPY backend/ .
RUN go mod download
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o management-console cmd/web/main.go

# Frontend builder stage  
FROM node:18-alpine AS frontend-builder
WORKDIR /app
COPY ui/ .
RUN npm ci --only=production
RUN npm run build

# Final runtime stage
FROM alpine:3.18
LABEL maintainer="Unity SDS Team"
LABEL description="Unity Management Console - Containerized deployment"

# Install system dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl \
    wget \
    unzip \
    git \
    bash \
    sqlite \
    && rm -rf /var/cache/apk/*

# Install Terraform
ARG TERRAFORM_VERSION=1.5.7
RUN wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    mv terraform /usr/local/bin/ && \
    chmod +x /usr/local/bin/terraform && \
    rm terraform_${TERRAFORM_VERSION}_linux_amd64.zip

# Install AWS CLI v2
RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip" && \
    unzip awscliv2.zip && \
    ./aws/install && \
    rm -rf awscliv2.zip aws

# Create non-root user
RUN addgroup -g 1001 unity && \
    adduser -D -u 1001 -G unity unity

# Create application directories
RUN mkdir -p /app /data/workdir /data/database /data/config && \
    chown -R unity:unity /app /data

# Set working directory
WORKDIR /app

# Copy built applications
COPY --from=backend-builder /app/management-console ./
COPY --from=frontend-builder /app/build ./ui/build

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