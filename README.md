# Go Auth Template

A complete authentication template with Go backend and Next.js frontend, featuring WebAuthn (passkey) support and magic link email authentication.

## 🚀 Features

- **Dual Authentication Methods**
  - 📧 Magic link email authentication
  - 🔐 WebAuthn passkeys (biometric/hardware keys)
- **Modern Tech Stack**
  - Go backend with Gin framework
  - Next.js 15 frontend with App Router
  - PostgreSQL database with GORM
  - Redis caching
- **Production Ready**
  - JWT session management
  - CSRF protection
  - Multi-client support (web, mobile, API)
  - Role-based access control
  - Type-safe codebase

## 📁 Project Structure

```
go-auth-template/
├── authserver/          # Go backend
│   ├── cmd/server/     # Application entry point
│   ├── internal/       # Private application code
│   └── pkg/           # Public libraries
└── front-end/          # Next.js frontend
    ├── src/app/       # App Router pages
    ├── src/components/ # React components
    └── src/lib/       # Utilities
```

## 🛠️ Quick Start

### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL
- Redis (optional)
- Make (for development commands)

### Setup
```bash
# Clone and setup
git clone <repository-url>
cd go-auth-template

# Setup environment
cp .env.example .env.local
# Edit .env.local with your values

# Install dependencies and tools
make install-deps
make install-tools

# Create database
createdb localDB

# Start development environment
make dev
```

### Development Commands
```bash
# Start both backend and frontend
make dev

# Start individually
make run-backend      # Backend only
make run-backend-air  # Backend with hot reload
make run-frontend     # Frontend only

# Database operations
make migrate          # Run migrations
make migrate-only     # Migrations only (no server)

# Testing
make test            # Run all tests
make test-ci         # Test with act (local CI)

# Utilities
make stop            # Stop all processes
make clean           # Clean build artifacts
make help            # Show all commands
```

### Access the Application
- Frontend: http://localhost:3000 (or 3001 if 3000 is busy)
- Backend API: http://localhost:8080

## 📚 Documentation

- [Backend Documentation](./authserver/README.md)
- [Frontend Documentation](./front-end/README.md)

## 🔧 Configuration

### Environment Variables

All environment variables are configured in the root `.env.local` file:

```bash
# Database
DATABASE_URL=postgres://postgres:password@localhost:5432/localDB?sslmode=disable
REDIS_URL=redis://localhost:6379

# Backend Configuration
JWT_SECRET=your-jwt-secret-change-in-production
RESEND_API_KEY=your-resend-api-key
PORT=8080
HOST=0.0.0.0

# Frontend Configuration
API_URL=http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080

# CI/CD (Optional)
RAILWAY_TOKEN=
RAILWAY_PROJECT_ID=
SLACK_WEBHOOK=
```

## 🎯 Authentication Flow

### Email Magic Link
1. User enters email address
2. System sends verification code via email
3. User enters code to authenticate
4. Session created with JWT token

### WebAuthn Passkeys
1. User registers passkey (optional)
2. User selects passkey authentication
3. Browser prompts for biometric/security key
4. Instant authentication on success

## 🏗️ Architecture

### Backend (Go)
- **Clean Architecture** with domain-driven design
- **GORM** for database operations
- **Gin** for HTTP routing and middleware
- **JWT** for stateless authentication
- **Redis** for caching and sessions

### Frontend (Next.js)
- **App Router** with Server Components
- **TypeScript** for type safety
- **Tailwind CSS** for styling
- **Radix UI** for accessible components
- **Server Actions** for form handling

## 🔒 Security Features

- **CSRF Protection** for web clients
- **Client Type Detection** for multi-platform support
- **Secure Session Management** with HTTP-only cookies
- **Rate Limiting** and request validation
- **WebAuthn** standard for passwordless authentication

## 🚀 Deployment

### Railway Deployment (Recommended)

This template is configured for Railway deployment using a single Docker container:

```bash
# Setup Railway
railway login
railway init

# Deploy
railway up
```

The `railway.json` configuration handles both frontend and backend deployment.

### Docker Deployment

#### Development
```bash
docker-compose up -d
```

#### Production
```bash
# Build images
./docker-build.sh compose production

# Deploy with production compose
docker-compose -f docker-compose.prod.yml up -d
```

#### Single Container
```bash
# Build single container
./docker-build.sh single

# Run in production
docker run -d \
  -p 3000:3000 \
  -p 8080:8080 \
  --env-file .env.local \
  go-auth-template:latest
```

### Manual Deployment

#### Backend Only
```bash
cd authserver
make build
# Deploy the binary in bin/authserver
```

#### Frontend Only
```bash
cd front-end
npm run build
# Deploy the .next folder
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Follow the existing code style
4. Add tests for new features
5. Submit a pull request

## 📄 License

[Your License Here]

## 🆘 Support

For questions and support:
- Check the documentation in each project folder
- Review the example implementations
- Open an issue for bugs or feature requests

---

**Happy coding!** 🎉
