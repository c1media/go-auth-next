# Multi-stage build for Go Auth Template
# This builds both the Go backend and Next.js frontend in a single container

# Stage 1: Build Next.js frontend
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY front-end/package*.json ./
RUN npm ci --only=production

# Copy source code and build
COPY front-end/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app/backend

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY authserver/go.mod authserver/go.sum ./
RUN go mod download

# Copy source code and build
COPY authserver/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/server/main.go

# Stage 3: Production runtime
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy built applications
COPY --from=backend-builder /app/backend/main ./backend
COPY --from=frontend-builder /app/frontend/.next ./frontend/.next
COPY --from=frontend-builder /app/frontend/public ./frontend/public
COPY --from=frontend-builder /app/frontend/package.json ./frontend/
COPY --from=frontend-builder /app/frontend/node_modules ./frontend/node_modules

# Create startup script
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'cd /app && ./backend &' >> /app/start.sh && \
    echo 'cd /app/frontend && npm start &' >> /app/start.sh && \
    echo 'wait' >> /app/start.sh && \
    chmod +x /app/start.sh

# Expose ports
EXPOSE 3000 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

CMD ["/app/start.sh"]