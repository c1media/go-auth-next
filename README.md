# Go Auth Template

A complete authentication template with Go backend and Next.js frontend, featuring WebAuthn (passkey) support and magic link email authentication.

## ✨ Features

- 📧 **Magic Link Email Authentication** 
- 🔐 **WebAuthn Passkeys** (biometric/hardware keys)
- 🎯 **Modern Stack**: Go + Next.js 15 + PostgreSQL + Redis
- 🔒 **Security**: JWT sessions, CSRF protection, rate limiting
- 🚀 **Production Ready**: CI/CD, Docker, Railway deployment
- 🧪 **Full Test Suite**: Integration tests, security scanning, load testing

## 🚀 Quick Start

### Prerequisites
- Go 1.23+
- Node.js 22+
- PostgreSQL + Redis (or use Docker)

### Development Setup
```bash
# Clone template
git clone https://github.com/c1media/go-auth-next.git my-project
cd my-project

# Setup environment
cp .env.example .env.local
# Edit .env.local with your database settings

# Start with Docker (easiest)
make docker-dev

# OR start manually
make install-deps
make dev
```

**Access:** Frontend at http://localhost:3000, Backend at http://localhost:8080

## 🛠️ Development Commands

```bash
make dev          # Start both backend + frontend
make docker-dev   # Start full Docker environment  
make test         # Run all tests
make build        # Build for production
make help         # Show all commands
```

## ⚡ Quick Deploy to Railway

### 1. Template Setup (First Time)
```bash
# Clone template for your project
git clone https://github.com/c1media/go-auth-next.git my-project
cd my-project
rm -rf .git && git init

# Update workflows for production (uncomment environment lines)
# In .github/workflows/cd.yml:
# environment: staging     # Uncomment
# environment: production  # Uncomment

# Commit your project
git add . && git commit -m "Initial commit from template"
git remote add origin https://github.com/yourusername/your-project.git
git push -u origin main
```

### 2. Railway Deployment
```bash
# Install Railway CLI
npm install -g @railway/cli

# Login and deploy
railway login
railway init
railway up

# Add services (Railway will auto-connect them)
railway add postgresql
railway add redis
```

### 3. GitHub Environments (Optional)
For deployment protection:
1. Go to GitHub Settings → Environments
2. Create `staging` and `production` environments
3. Add protection rules and required reviewers

### 4. Required Secrets
Add these to your GitHub repository secrets:
- `RESEND_API_KEY` - For email sending
- `RAILWAY_TOKEN` - For deployments (get from Railway dashboard)
- `RAILWAY_PROJECT_ID` - Your Railway project ID

## 🔧 Configuration

All settings in `.env.local`:

```bash
# Database (Railway auto-provides these)
DATABASE_URL=postgres://postgres:password@localhost:5432/localDB?sslmode=disable
REDIS_URL=redis://localhost:6379

# Required
JWT_SECRET=your-secure-jwt-secret-256-bits
RESEND_API_KEY=re_your_resend_api_key
FROM_EMAIL=your-app@domain.com

# Frontend
API_URL=http://localhost:8080
NEXT_PUBLIC_API_URL=http://localhost:8080
```

## 🏗️ Project Structure

```
├── authserver/          # Go backend (Gin + GORM)
├── front-end/           # Next.js 15 frontend  
├── .github/workflows/   # CI/CD pipelines
├── Dockerfile           # Single container for Railway
└── Makefile            # Development commands
```

## 🎯 Authentication Flow

**Magic Link:**
1. User enters email → receives verification code
2. User enters code → authenticated with JWT

**WebAuthn Passkeys:**
1. User registers passkey (Face ID, Touch ID, hardware key)
2. Future logins: instant biometric authentication

## 🧪 CI/CD Pipeline

Automated testing on every push:
- ✅ **Backend Tests**: Go unit tests, security scanning
- ✅ **Frontend Tests**: Next.js build, type checking, linting  
- ✅ **Security**: Trivy vulnerability scanning, gosec analysis
- ✅ **Integration**: Full stack API testing
- ✅ **Docker**: Multi-platform container builds

## 🚀 Deployment

**Railway (Recommended):**
- Single container deployment
- Managed PostgreSQL + Redis
- Auto-scaling and zero-downtime deploys
- Built-in SSL certificates

**Docker:**
- Production-ready Dockerfile included
- Multi-stage build optimized for size
- Supports any container platform

## 🔒 Security

- **WebAuthn** standard implementation
- **JWT** with secure HTTP-only cookies
- **CSRF** protection for web clients
- **Rate limiting** on auth endpoints
- **Input validation** and sanitization
- **Security scanning** in CI/CD

## 🆘 Support

- 📖 Check individual README files in `authserver/` and `front-end/`
- 🐛 Open issues for bugs or feature requests
- 💡 Review the code examples and tests

---

**Ready to build your auth system!** 🎉