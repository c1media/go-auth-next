# 🚀 Deployment Guide

This comprehensive guide covers all deployment options for the Go Auth Template, from quick cloud deployments to enterprise-grade infrastructure.

## 📋 Prerequisites

- Docker and Docker Compose
- Git
- Domain name (for production)
- SSL certificate (for production)

## 🥇 **Recommended: Modern PaaS (2024 Best Practice)**

### **Frontend → Vercel**
```bash
# 1. Push to GitHub
git add . && git commit -m "Initial commit" && git push

# 2. Connect to Vercel
# Visit vercel.com → Import Project → Select GitHub repo

# 3. Environment Variables (in Vercel dashboard)
NEXT_PUBLIC_API_URL=https://your-backend.railway.app
```

### **Backend → Railway**
```bash
# 1. Connect GitHub repo to Railway
# Visit railway.app → New Project → Deploy from GitHub

# 2. Environment Variables (in Railway dashboard)
DATABASE_URL=${{Postgres.DATABASE_URL}}  # Auto-generated
JWT_SECRET=your-secret-here
RESEND_API_KEY=your-api-key
PORT=${{PORT}}  # Auto-set by Railway
```

### **Database → Railway PostgreSQL**
```bash
# Add PostgreSQL service in Railway dashboard
# DATABASE_URL automatically provided to your backend
```

### **✅ Benefits:**
- 🚀 **Deploy in 5 minutes**
- 🔄 **Auto-deploy on git push**
- 📊 **Built-in monitoring**
- 🌍 **Global CDN**
- 💰 **Free tier available**
- 🔒 **SSL certificates included**

---

## 🏠 Local Development

```bash
# Clone the repository
git clone <your-repo-url>
cd go-auth-template

# Start with Docker Compose
docker-compose up -d

# Or start individual services
cd authserver && go run cmd/server/main.go &
cd front-end && npm run dev
```

**Access Points:**
- Frontend: http://localhost:3000
- Backend API: http://localhost:8080
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## 🐳 Docker Deployment Options

### Single Container (Recommended for simple deployments)

```bash
# Build the combined image
docker build -t go-auth-template .

# Run with environment variables
docker run -d \
  --name go-auth-template \
  -p 3000:3000 \
  -e DATABASE_URL="your-database-url" \
  -e REDIS_URL="your-redis-url" \
  -e JWT_SECRET="your-jwt-secret" \
  go-auth-template
```

### Multi-Service with Docker Compose

```bash
# Development
docker-compose up -d

# Production
docker-compose -f docker-compose.prod.yml up -d
```

## ☁️ Cloud Platform Deployments

### Railway

1. **Connect Repository:**
   ```bash
   # Install Railway CLI
   npm install -g @railway/cli
   railway login
   ```

2. **Deploy Backend:**
   ```bash
   railway new
   railway add --service postgresql
   railway add --service redis
   
   # Set environment variables
   railway variables set JWT_SECRET=your-secret
   railway variables set DATABASE_URL=${{Postgres.DATABASE_URL}}
   railway variables set REDIS_URL=${{Redis.REDIS_URL}}
   
   # Deploy
   railway up
   ```

3. **Deploy Frontend:**
   ```bash
   railway add --service frontend
   railway variables set API_URL=https://your-backend.railway.app
   railway variables set NEXT_PUBLIC_API_URL=https://your-backend.railway.app
   railway up
   ```

### Vercel + PlanetScale

1. **Deploy Frontend to Vercel:**
   ```bash
   npm i -g vercel
   cd front-end
   vercel --prod
   ```

2. **Environment Variables:**
   ```bash
   vercel env add API_URL
   vercel env add NEXT_PUBLIC_API_URL
   ```

3. **Deploy Backend to Railway/Fly.io:**
   ```bash
   # Set production API URL in Vercel
   vercel env add API_URL https://your-backend.railway.app
   ```

### AWS ECS

1. **Build and Push Images:**
   ```bash
   # Build images
   docker build -t your-registry/auth-backend ./authserver
   docker build -t your-registry/auth-frontend ./front-end
   
   # Push to ECR
   aws ecr get-login-password --region us-east-1 | docker login --username AWS --password-stdin your-registry
   docker push your-registry/auth-backend
   docker push your-registry/auth-frontend
   ```

2. **Create ECS Task Definition:**
   ```json
   {
     "family": "go-auth-template",
     "networkMode": "awsvpc",
     "requiresCompatibilities": ["FARGATE"],
     "cpu": "512",
     "memory": "1024",
     "containerDefinitions": [
       {
         "name": "backend",
         "image": "your-registry/auth-backend",
         "portMappings": [{"containerPort": 8080}],
         "environment": [
           {"name": "DATABASE_URL", "value": "your-database-url"},
           {"name": "JWT_SECRET", "value": "your-secret"}
         ]
       }
     ]
   }
   ```

### Google Cloud Run

```bash
# Build and deploy backend
gcloud builds submit --tag gcr.io/PROJECT-ID/auth-backend ./authserver
gcloud run deploy auth-backend --image gcr.io/PROJECT-ID/auth-backend --platform managed

# Build and deploy frontend  
gcloud builds submit --tag gcr.io/PROJECT-ID/auth-frontend ./front-end
gcloud run deploy auth-frontend --image gcr.io/PROJECT-ID/auth-frontend --platform managed
```

## 🏢 **Enterprise: Kubernetes**

```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: go-auth-template
spec:
  replicas: 3
  selector:
    matchLabels:
      app: go-auth-template
  template:
    spec:
      containers:
      - name: backend
        image: go-auth-backend:latest
      - name: frontend
        image: go-auth-frontend:latest
```

## 🔧 Environment Variables

### Backend (.env)
```bash
# Database
DATABASE_URL=postgres://user:pass@host:5432/dbname?sslmode=disable

# Redis
REDIS_URL=redis://localhost:6379

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Server
PORT=8080
GIN_MODE=release

# CORS
ALLOWED_ORIGINS=https://yourapp.com,https://www.yourapp.com

# Email (if using email provider)
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USER=apikey
SMTP_PASS=your-sendgrid-api-key
```

### Frontend (.env.local)
```bash
# API Configuration
API_URL=https://your-backend.railway.app
NEXT_PUBLIC_API_URL=https://your-backend.railway.app


# Production
NODE_ENV=production
```

## 📊 **Deployment Comparison**

| Method | Setup Time | Cost | Maintenance | Scalability | Best For |
|--------|------------|------|-------------|-------------|----------|
| **PaaS** | 5 min | $0-20/mo | None | Auto | Most projects |
| **Serverless** | 10 min | Pay/use | None | Infinite | APIs, functions |
| **Docker VPS** | 1 hour | $5-50/mo | Medium | Manual | Learning, control |
| **Kubernetes** | 1 day | $100+/mo | High | Enterprise | Large apps |

## 🎯 **Quick Start Recommendations**

### **🚀 Just want it deployed fast?**
```bash
# Use Railway (5 minutes)
1. Push code to GitHub
2. Connect to Railway
3. Add PostgreSQL service
4. Deploy frontend to Vercel
```

### **💰 Want it free?**
```bash
# Use free tiers
Frontend: Vercel (free)
Backend: Railway (free $5 credit)
Database: Railway PostgreSQL (free tier)
```

### **🏢 Need enterprise features?**
```bash
# Use Docker + AWS/GCP
./docker-build.sh compose production
# Deploy to ECS/GKE
```

### **🎓 Want to learn infrastructure?**
```bash
# Use Docker on VPS
./docker-build.sh single
# Deploy to DigitalOcean droplet
```

## 📊 Monitoring & Logging

### Application Monitoring
```bash
# Add health check endpoints
curl https://your-backend.com/health
curl https://your-frontend.com/api/health
```

### Logging Configuration
```yaml
# docker-compose.yml
services:
  backend:
    logging:
      driver: "json-file"
      options:
        max-size: "10m"
        max-file: "3"
```

### Metrics Collection
```bash
# Add Prometheus metrics endpoint
curl https://your-backend.com/metrics
```

## 🔒 Security Considerations

### SSL/TLS Configuration
```nginx
# nginx.conf
server {
    listen 443 ssl;
    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;
    
    location / {
        proxy_pass http://frontend:3000;
    }
    
    location /api {
        proxy_pass http://backend:8080;
    }
}
```

### Firewall Rules
```bash
# Only allow HTTPS traffic
ufw allow 443/tcp
ufw allow 80/tcp  # For redirects
ufw deny 8080/tcp # Block direct backend access
```

## 🔄 CI/CD Pipeline

The repository includes GitHub Actions workflows for:

- **CI Pipeline** (`.github/workflows/ci.yml`)
  - Automated testing
  - Security scanning
  - Docker builds

- **CD Pipeline** (`.github/workflows/cd.yml`)
  - Staging deployments
  - Production deployments
  - Container registry publishing

- **Release Management** (`.github/workflows/release.yml`)
  - Automated releases
  - Asset building
  - Changelog generation

## 🚨 Rollback Strategy

### Quick Rollback
```bash
# Docker Compose
docker-compose down
git checkout previous-tag
docker-compose up -d

# Kubernetes
kubectl rollout undo deployment/auth-backend
kubectl rollout undo deployment/auth-frontend
```

### Database Migrations
```bash
# Always backup before migration
pg_dump database_name > backup.sql

# Rollback migration if needed
goose -dir authserver/migrations postgres "connection-string" down
```

## 📋 Post-Deployment Checklist

- [ ] Health checks passing
- [ ] SSL certificate valid
- [ ] Environment variables set
- [ ] Database migrations applied
- [ ] Monitoring alerts configured
- [ ] Backup strategy implemented
- [ ] Security headers configured
- [ ] Performance metrics baseline established

## 🆘 Troubleshooting

### Common Issues

**Frontend can't connect to backend:**
```bash
# Check environment variables
echo $API_URL
echo $NEXT_PUBLIC_API_URL

# Check CORS settings
curl -H "Origin: https://yourfrontend.com" https://yourbackend.com/api/health
```

**Database connection issues:**
```bash
# Test database connection
psql $DATABASE_URL -c "SELECT 1;"

# Check if database exists
psql $DATABASE_URL -c "\l"
```

**Docker build failures:**
```bash
# Clear Docker cache
docker system prune -a

# Check build logs
docker build --no-cache .
```

## 🎉 **Conclusion**

**2024 Recommendation:** Start with **Railway + Vercel** for fastest deployment, switch to Docker/K8s later if you need more control.

The beauty of this template is you can start simple and scale up as needed! 🚀

### Support

For deployment issues:
1. Check GitHub Issues
2. Review logs: `docker-compose logs`
3. Open an issue with deployment details