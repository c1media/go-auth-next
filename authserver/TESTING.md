# Testing Simple Auth (Level 2)

## 🚀 **Quick Start**

### 1. Database Setup
```bash
# Create database
createdb simple_auth_roles

# Update .env if needed
cp .env.example .env
```

### 2. Run Server
```bash
# Start with migrations
go run ./cmd/server --migrate

# Or just start
go run ./cmd/server
```

Server starts on: `http://localhost:8080`

### 3. Test API

Open `test-auth.html` in your browser for a visual test interface.

Or use curl:

```bash
# 1. Send login code
curl -X POST http://localhost:8080/api/v1/auth/send-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# 2. Check server logs for the 6-digit code (since Resend isn't configured)
# Example code: ABC123

# 3. Verify code
curl -X POST http://localhost:8080/api/v1/auth/verify-code \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","code":"ABC123"}'

# 4. Test protected route (use token from step 3)
curl -X GET http://localhost:8080/api/v1/protected/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN_HERE"
```

## 🔧 **Configuration**

### Environment Variables
```bash
# Required
DATABASE_URL=postgres://user:pass@localhost:5432/simple_auth_roles?sslmode=disable
JWT_SECRET=your-secret-key

# Optional (for real email sending)
RESEND_API_KEY=re_your_resend_api_key
FROM_EMAIL=auth@yourapp.com
FROM_NAME=Your App

# Optional (for Redis caching)
REDIS_URL=redis://localhost:6379
```

### With Real Email (Resend)
Set `RESEND_API_KEY` in `.env` to: `re_72bo2Cds_6Gja3TYCsHVzQ8WTs6KRrfZE`

## 📋 **API Endpoints**

### Authentication
- `POST /api/v1/auth/send-code` - Send login code
- `POST /api/v1/auth/verify-code` - Verify login code
- `POST /api/v1/auth/create-user` - Create user (admin only)
- `PUT /api/v1/auth/users/:id/role` - Update user role (admin only)

### Protected Routes
- `GET /api/v1/protected/profile` - Get user profile
- `GET /api/v1/protected/admin/users` - Get all users (admin only)

### System
- `GET /health` - Health check

## 🎯 **Test Scenarios**

### 1. **Basic Authentication**
1. Send code to email
2. Verify code and get JWT token
3. Access protected profile route

### 2. **Role-Based Access**
1. Login as admin user
2. Access admin-only routes
3. Try accessing with user role (should fail)

### 3. **Error Handling**
1. Invalid email format
2. Wrong verification code
3. Expired code
4. Missing/invalid JWT token

## 🔒 **Default Roles**

- `admin` - Full system access
- `moderator` - Content management
- `user` - Basic access (default)

## 📁 **Project Structure**

```
simple-auth-roles/
├── cmd/server/main.go          # Application entry point
├── internal/
│   ├── auth/                   # Authentication domain
│   │   ├── domain.go
│   │   ├── handlers/
│   │   ├── service/
│   │   └── repository/
│   ├── config/                 # Configuration
│   ├── middleware/             # HTTP middleware
│   └── types/                  # Data models
├── pkg/
│   ├── cache/                  # Redis + in-memory cache
│   ├── database/               # Database connection
│   └── email/                  # Resend email service
├── .env                        # Environment variables
├── test-auth.html              # Test interface
└── README.md
```

## 🚁 **Railway Deployment**

Ready to deploy to Railway:

```bash
# Deploy to Railway
railway up
```

The `railway.json` configuration is already set up for production deployment.

## 💡 **Tips**

1. **Development**: Use in-memory cache (no Redis needed)
2. **Testing**: Check server logs for login codes
3. **Production**: Set strong JWT_SECRET and configure Resend
4. **Scaling**: Add Redis for distributed caching

This is a **Level 2** auth system - perfect for small team tools and internal applications!