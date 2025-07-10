# Multi-stage build for Go Auth Template
# This builds both the Go backend and Next.js frontend in a single container

# Stage 1: Build Next.js frontend
FROM node:22-alpine AS frontend-builder

# Update packages for security
RUN apk update && apk upgrade && apk add --no-cache dumb-init

WORKDIR /app/frontend

# Copy package files
COPY front-end/package*.json ./
RUN npm ci

# Copy source code and build
COPY front-end/ ./
RUN npm run build

# Stage 2: Build Go backend
FROM golang:1.23-alpine AS backend-builder

# Update packages for security
RUN apk update && apk upgrade

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
FROM alpine:3.19

# Update packages and install runtime dependencies
RUN apk update && apk upgrade && apk add --no-cache ca-certificates tzdata dumb-init nodejs npm

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nextjs -u 1001

WORKDIR /app

# Copy built applications
COPY --from=backend-builder /app/backend/main ./backend
COPY --from=frontend-builder /app/frontend/.next ./frontend/.next
COPY --from=frontend-builder /app/frontend/package.json ./frontend/
COPY --from=frontend-builder /app/frontend/node_modules ./frontend/node_modules

# Copy public directory - create it first if it doesn't exist
RUN mkdir -p ./frontend/public
COPY --from=frontend-builder /app/frontend/public ./frontend/public

# Copy environment file (if exists)
COPY .env.local* /app/

# Create startup script that loads environment
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'if [ -f /app/.env.local ]; then export $(cat /app/.env.local | sed "s/#.*//g" | xargs); fi' >> /app/start.sh && \
    echo 'echo "Starting backend..."' >> /app/start.sh && \
    echo 'cd /app && MIGRATE=true ./backend &' >> /app/start.sh && \
    echo 'BACKEND_PID=$!' >> /app/start.sh && \
    echo 'echo "Starting frontend..."' >> /app/start.sh && \
    echo 'cd /app/frontend && PORT=3000 npm run start &' >> /app/start.sh && \
    echo 'FRONTEND_PID=$!' >> /app/start.sh && \
    echo 'echo "Both services started. Backend PID: $BACKEND_PID, Frontend PID: $FRONTEND_PID"' >> /app/start.sh && \
    echo 'wait $BACKEND_PID $FRONTEND_PID' >> /app/start.sh && \
    chmod +x /app/start.sh

# Change ownership to non-root user
RUN chown -R nextjs:nodejs /app

# Switch to non-root user
USER nextjs

# Expose ports
EXPOSE 3000 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:3000/ || exit 1

CMD ["/app/start.sh"]