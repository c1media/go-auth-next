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

### Option 1: Docker (Recommended)
```bash
# Clone and setup
git clone <repository-url>
cd go-auth-template

# Setup environment
cp .env.example .env
# Edit .env with your values

# Start with Docker Compose
docker-compose up -d

# Or build and run single container
./docker-build.sh single
docker run -p 3000:3000 -p 8080:8080 go-auth-template:latest
```

### Option 2: Manual Setup

#### Prerequisites
- Go 1.21+
- Node.js 18+
- PostgreSQL
- Redis (optional)

#### Backend Setup
```bash
cd authserver
cp .env.example .env
# Edit .env with your database credentials
go mod download
go run cmd/server/main.go
```

#### Frontend Setup
```bash
cd front-end
npm install
cp .env.example .env.local
# Edit .env.local with your API URL
npm run dev
```

### Access the Application
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080

## 📚 Documentation

- [Backend Documentation](./authserver/README.md)
- [Frontend Documentation](./front-end/README.md)

## 🔧 Configuration

### Environment Variables

**Backend (.env)**
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string (optional)
- `JWT_SECRET` - JWT signing secret
- `RESEND_API_KEY` - Email service API key

**Frontend (.env.local)**
- `API_URL` - Backend API URL
- `NEXTAUTH_SECRET` - Session encryption secret

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

### Docker Deployment (Recommended)

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
  --env-file .env \
  go-auth-template:latest
```

### Traditional Deployment

#### Backend
- Build: `go build -o main cmd/server/main.go`
- Deploy to any Go-compatible platform

#### Frontend
- Build: `npm run build`
- Deploy to Vercel, Netlify, or any Node.js platform

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
