# 🐳 Docker Deployment Guide

Complete guide for running Task Manager with Docker.

---

## 🚀 Quick Start

### Run Everything (Recommended)

```bash
# Start all services (database + backend + frontend)
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

**Access:**
- Frontend: http://localhost:3000
- Backend: http://localhost:8081
- Database: localhost:5433

### Stop All Services

```bash
docker-compose down
```

---

## 📦 Docker Services

### Available Services

| Service | Port | Description |
|---------|------|-------------|
| **frontend** | 3000 | Next.js web application |
| **backend** | 8081 | Go API server |
| **db** | 5433 | PostgreSQL database |
| **redis** | 6380 | Redis cache |

---

## 🔧 Step-by-Step

### 1. Build Images

```bash
# Build all services
docker-compose build

# Or build specific service
docker-compose build backend
docker-compose build frontend
```

### 2. Start Services

```bash
# Start all services
docker-compose up -d

# Start specific services
docker-compose up -d db redis        # Database only
docker-compose up -d backend          # Backend + dependencies
docker-compose up -d frontend         # All services
```

### 3. Check Status

```bash
# View running containers
docker-compose ps

# Expected output:
# NAME                    STATUS    PORTS
# taskmanager-frontend    Up        0.0.0.0:3000->3000/tcp
# taskmanager-backend     Up        0.0.0.0:8081->8081/tcp
# taskmanager-db          Up        0.0.0.0:5433->5432/tcp
# taskmanager-redis       Up        0.0.0.0:6380->6379/tcp
```

### 4. View Logs

```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f db
```

### 5. Stop Services

```bash
# Stop all services
docker-compose down

# Stop and remove volumes (WARNING: deletes data)
docker-compose down -v

# Stop specific service
docker-compose stop backend
```

---

## 🗄️ Database Setup

### Run Migrations

```bash
# Method 1: Using Makefile
make migrate-up

# Method 2: Manual
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager -c "CREATE TABLE IF NOT EXISTS schema_migrations (version bigint NOT NULL PRIMARY KEY, dirty boolean NOT NULL);"

# Or run migration files manually
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/001_developers.up.sql
```

### Verify Database

```bash
# Connect to database
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager

# List tables
\dt

# Exit
\q
```

---

## 🔐 Production Configuration

### Environment Variables

Create `.env` file in project root:

```bash
# Backend
APP_PORT=8081
APP_ENV=production
DB_HOST=db
DB_PORT=5432
DB_USER=taskmanager
DB_PASS=strong-password-change-this
DB_NAME=taskmanager
REDIS_HOST=redis
REDIS_PORT=6379
JWT_SECRET=use-strong-random-secret-min-32-characters
JWT_EXPIRY=24h
CORS_ALLOWED_ORIGINS=https://your-domain.com

# Frontend
NEXT_PUBLIC_API_URL=https://api.your-domain.com/api/v1
```

### Update docker-compose.yml

```yaml
services:
  backend:
    environment:
      - DB_PASS=strong-password-change-this
      - JWT_SECRET=use-strong-random-secret-min-32-characters
      - CORS_ALLOWED_ORIGINS=https://your-domain.com

  frontend:
    environment:
      - NEXT_PUBLIC_API_URL=https://api.your-domain.com/api/v1
```

---

## 🌐 Reverse Proxy (Optional)

### Nginx Configuration

```nginx
# /etc/nginx/sites-available/taskmanager
server {
    listen 80;
    server_name your-domain.com;

    # Frontend
    location / {
        proxy_pass http://localhost:3000;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_cache_bypass $http_upgrade;
    }

    # Backend API
    location /api/ {
        proxy_pass http://localhost:8081/api/;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 🔧 Troubleshooting

### Port Already in Use

```bash
# Find process using port
lsof -i :3000
lsof -i :8081

# Kill process
kill -9 <PID>

# Or change port in docker-compose.yml
ports:
  - "3001:3000"  # Change 3000 to 3001
```

### Database Connection Failed

```bash
# Check database logs
docker-compose logs db

# Restart database
docker-compose restart db

# Check if migrations ran
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager -c "\dt"
```

### Frontend Can't Connect to Backend

```bash
# Check backend is running
curl http://localhost:8081/health

# Check Docker network
docker network inspect taskmanager-network

# Restart backend
docker-compose restart backend
```

### Rebuild Containers

```bash
# Rebuild all
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# Rebuild specific service
docker-compose build --no-cache backend
docker-compose up -d backend
```

### Clear Everything

```bash
# Stop and remove all
docker-compose down -v

# Remove images
docker rmi taskmanager-backend taskmanager-frontend

# Start fresh
docker-compose up -d --build
```

---

## 📊 Monitoring

### Container Stats

```bash
# View resource usage
docker stats

# Specific container
docker stats taskmanager-backend
```

### Health Checks

```bash
# Backend health
curl http://localhost:8081/health

# Database health
docker exec taskmanager-db pg_isready -U taskmanager

# Redis health
docker exec taskmanager-redis redis-cli ping
```

---

## 🔄 Updates

### Update Code & Redeploy

```bash
# Pull latest code
git pull origin main

# Rebuild and restart
docker-compose down
docker-compose build
docker-compose up -d
```

### Update Specific Service

```bash
# Rebuild backend only
docker-compose build backend
docker-compose up -d backend

# Rebuild frontend only
docker-compose build frontend
docker-compose up -d frontend
```

---

## 🗂️ Docker Commands Cheat Sheet

```bash
# Start services
docker-compose up -d                # Start all
docker-compose up -d backend        # Start backend + deps

# Stop services
docker-compose down                 # Stop all
docker-compose stop backend         # Stop specific

# Logs
docker-compose logs -f              # All logs
docker-compose logs -f backend      # Backend logs

# Rebuild
docker-compose build                # Build all
docker-compose build --no-cache     # Force rebuild

# Status
docker-compose ps                   # List containers
docker stats                        # Resource usage

# Cleanup
docker-compose down -v              # Remove volumes
docker system prune                 # Clean unused resources

# Debug
docker exec -it taskmanager-backend sh   # Shell in container
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager  # DB shell
```

---

## ✅ Quick Checklist

After running `docker-compose up -d`:

- [ ] All 4 containers running (`docker-compose ps`)
- [ ] Backend healthy (`curl http://localhost:8081/health`)
- [ ] Frontend accessible (http://localhost:3000)
- [ ] Database has tables (run migrations)
- [ ] Can register/login user

---

## 🚀 Production Deployment

### 1. Prepare Environment

```bash
# Clone repository
git clone https://github.com/ardani17/task-manager.git
cd task-manager

# Create .env file
cp .env.example .env
nano .env  # Edit with production values
```

### 2. Update Configuration

- Change all passwords
- Set strong JWT secret
- Configure CORS origins
- Update API URL

### 3. Deploy

```bash
# Build and start
docker-compose up -d --build

# Run migrations
make migrate-up

# Check status
docker-compose ps
docker-compose logs -f
```

### 4. Setup SSL (Optional)

Use Let's Encrypt with Nginx reverse proxy.

---

## 📞 Need Help?

- **Logs:** `docker-compose logs -f`
- **Status:** `docker-compose ps`
- **Documentation:** [INSTALLATION.md](INSTALLATION.md)

---

**Docker deployment ready! 🐳**
