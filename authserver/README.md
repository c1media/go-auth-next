# Simple Auth with Roles (Level 2)

A clean, standalone authentication service with role-based access control. Perfect for small team tools, internal apps, and community platforms.

## Features

✅ **Magic Link Authentication** - Passwordless login via email  
✅ **Role-Based Access Control** - Admin, Moderator, User roles  
✅ **JWT Token Management** - Secure session handling  
✅ **Permission System** - Granular permission checking  
✅ **RESTful API** - Clean HTTP endpoints  
✅ **PostgreSQL Database** - Reliable data storage  

## Roles & Permissions

### **Admin**
- Full system access
- User management
- All permissions

### **Moderator** 
- Content management
- Read, write, moderate permissions
- Limited user management

### **User**
- Basic access
- Read-only permissions

## Quick Start

### 1. Prerequisites
- Go 1.21+
- PostgreSQL database

### 2. Setup
```bash
# Clone and setup
git clone <repository>
cd simple-auth-roles

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Edit .env with your database and email settings
```

### 3. Database Setup
```bash
# Create database
createdb simple_auth_roles

# Update DATABASE_URL in .env
DATABASE_URL=postgres://username:password@localhost:5432/simple_auth_roles?sslmode=disable
```

### 4. Run Application
```bash
go run main.go
```

The server starts on `http://localhost:8080`

## API Endpoints

### Authentication

#### Send Login Code
```http
POST /api/v1/auth/send-code
Content-Type: application/json

{
  "email": "user@example.com"
}
```

#### Verify Login Code
```http
POST /api/v1/auth/verify-code
Content-Type: application/json

{
  "email": "user@example.com",
  "code": "ABC123"
}
```

Response:
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "",
    "role": "user",
    "is_active": true
  },
  "token": "jwt-token-here",
  "message": "Authentication successful"
}
```

### Protected Routes

#### Get Profile
```http
GET /api/v1/protected/profile
Authorization: Bearer <jwt-token>
```

#### Admin: Get All Users
```http
GET /api/v1/protected/admin/users
Authorization: Bearer <admin-jwt-token>
```

#### Admin: Create User
```http
POST /api/v1/auth/create-user
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "email": "newuser@example.com",
  "name": "New User",
  "company": "Example Corp",
  "role": "moderator"
}
```

#### Admin: Update User Role
```http
PUT /api/v1/auth/users/123/role
Authorization: Bearer <admin-jwt-token>
Content-Type: application/json

{
  "role": "admin"
}
```

## Environment Variables

```bash
# Server
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable

# JWT
JWT_SECRET=your-secret-key

# Email (for magic links)
FROM_EMAIL=auth@yourapp.com
FROM_NAME=Your App
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USERNAME=your-email
SMTP_PASSWORD=your-password
```

## Usage Examples

### Frontend Integration

```javascript
// Send login code
const sendCode = async (email) => {
  const response = await fetch('/api/v1/auth/send-code', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email })
  });
  return response.json();
};

// Verify code and get token
const verifyCode = async (email, code) => {
  const response = await fetch('/api/v1/auth/verify-code', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email, code })
  });
  const data = await response.json();
  
  if (data.token) {
    localStorage.setItem('auth-token', data.token);
  }
  
  return data;
};

// Make authenticated requests
const fetchProfile = async () => {
  const token = localStorage.getItem('auth-token');
  const response = await fetch('/api/v1/protected/profile', {
    headers: { 'Authorization': `Bearer ${token}` }
  });
  return response.json();
};
```

### Role-Based UI

```javascript
// Check user permissions
const user = getCurrentUser(); // From JWT or API

if (user.role === 'admin') {
  showAdminPanel();
} else if (user.role === 'moderator') {
  showModeratorTools();
} else {
  showUserInterface();
}

// Permission-based features
if (user.hasPermission('moderate')) {
  showModerationButtons();
}
```

## Architecture

```
simple-auth-roles/
├── main.go                    # Application entry point
├── internal/
│   ├── config/               # Configuration management
│   ├── database/             # Database connection & migrations
│   ├── types/                # Data models and types
│   ├── auth/                 # Authentication domain
│   │   ├── domain.go         # Domain layer
│   │   ├── handlers/         # HTTP handlers
│   │   ├── service/          # Business logic
│   │   └── repository/       # Data access
│   └── middleware/           # HTTP middleware
├── .env.example              # Environment template
└── README.md                 # This file
```

## Deployment

### Docker
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]
```

### Environment
- Set `ENVIRONMENT=production` in production
- Use strong `JWT_SECRET` 
- Configure proper SMTP settings
- Set up SSL/TLS for database connections

## Development

### Add New Roles
1. Add role constant in `internal/types/user.go`
2. Update `ValidateRole()` function
3. Update `HasPermission()` logic
4. Add role to middleware checks

### Extend Permissions
1. Define new permissions in `HasPermission()` method
2. Use `RequirePermission()` middleware in routes
3. Update frontend permission checks

This is a **Level 2** auth system - perfect for applications that need role-based access control without the complexity of multi-tenant accounts.