# Installation Guide - Task Manager

Complete step-by-step guide to install and run Task Manager on your system.

---

## 📋 Prerequisites

### Required Software
- **Go** 1.22+ ([Download](https://golang.org/dl/))
- **Node.js** 20+ ([Download](https://nodejs.org/))
- **Docker** & Docker Compose ([Install Guide](https://docs.docker.com/get-docker/))
- **Git** ([Download](https://git-scm.com/downloads))

### Verify Installations
```bash
go version          # Go 1.22+
node --version      # Node 20+
npm --version       # npm 10+
docker --version    # Docker 24+
docker-compose --version  # Docker Compose 2+
git --version       # Git 2.x
```

---

## 🚀 Quick Start (5 minutes)

### 1. Clone Repository
```bash
git clone https://github.com/ardani17/task-manager.git
cd task-manager
```

### 2. Setup Environment
```bash
# Copy environment template
cp .env.example .env

# Edit if needed (optional)
nano .env
```

### 3. Start Database Services
```bash
# Start PostgreSQL + Redis
make db-up

# Or using docker-compose directly
docker-compose up -d db redis
```

### 4. Run Database Migrations
```bash
# Install migrate tool (if not installed)
# macOS: brew install golang-migrate
# Linux: curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz && sudo mv migrate /usr/local/bin/

# Run migrations
make migrate-up
```

### 5. Start Backend
```bash
# Terminal 1: Backend
cd backend
go mod download
go run cmd/server/main.go
```

Backend will run on: http://localhost:8081

### 6. Start Frontend
```bash
# Terminal 2: Frontend
cd frontend
npm install
npm run dev
```

Frontend will run on: http://localhost:3000

### 7. Access Application
Open your browser: http://localhost:3000

---

## 📦 Detailed Installation

### Option A: Using Make (Recommended)

```bash
# 1. Setup everything
make setup

# 2. Start services
make docker-up

# 3. Run migrations
make migrate-up

# 4. Start backend (Terminal 1)
make backend

# 5. Start frontend (Terminal 2)
make frontend
```

### Option B: Manual Installation

#### Step 1: Database Setup

**PostgreSQL + Redis via Docker:**
```bash
docker-compose up -d db redis
```

**Verify database is running:**
```bash
docker ps | grep taskmanager
```

Expected output:
```
taskmanager-db      postgres:16-alpine   Up      0.0.0.0:5433->5432/tcp
taskmanager-redis   redis:7-alpine       Up      0.0.0.0:6380->6379/tcp
```

#### Step 2: Database Migrations

**Install migrate tool:**

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

**Windows:**
```bash
scoop install migrate
```

**Run migrations:**
```bash
make migrate-up
# or
migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" up
```

**Verify migrations:**
```bash
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager -c "\dt"
```

Expected output:
```
 Schema |       Name        | Type  |    Owner
--------+-------------------+-------+-------------
 public | activities        | table | taskmanager
 public | developers        | table | taskmanager
 public | projects          | table | taskmanager
 public | schema_migrations | table | taskmanager
 public | tasks             | table | taskmanager
 public | teams             | table | taskmanager
```

#### Step 3: Backend Setup

```bash
cd backend

# Download dependencies
go mod download

# Run server
go run cmd/server/main.go
```

Expected output:
```
2026/02/27 14:00:00 🚀 Server starting on :8081
2026/02/27 14:00:00 ✅ Connected to database
2026/02/27 14:00:00 ✅ Redis connected
```

**Verify backend:**
```bash
curl http://localhost:8081/health
```

Expected response:
```json
{
  "status": "healthy",
  "timestamp": "2026-02-27T14:00:00Z",
  "version": "1.0.0"
}
```

#### Step 4: Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Development mode
npm run dev
```

Expected output:
```
   ▲ Next.js 16.0.0
   - Local:        http://localhost:3000
   - Environments: .env

 ✓ Starting...
 ✓ Ready in 2.5s
```

---

## ⚙️ Configuration

### Environment Variables

**Backend (.env in project root):**
```bash
# Server Configuration
APP_PORT=8081
APP_ENV=development

# Database Configuration
DB_HOST=localhost
DB_PORT=5433
DB_USER=taskmanager
DB_PASS=taskmanager123
DB_NAME=taskmanager

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6380

# JWT Configuration
JWT_SECRET=your-secret-key-change-in-production
JWT_EXPIRY=24h

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

**Frontend (frontend/.env.local):**
```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

### Production Configuration

**Backend (.env):**
```bash
APP_ENV=production
APP_PORT=8081

DB_HOST=your-production-db-host
DB_PORT=5432
DB_USER=production_user
DB_PASS=strong-password-here
DB_NAME=taskmanager_production

REDIS_HOST=your-redis-host
REDIS_PORT=6379

JWT_SECRET=use-strong-random-secret-min-32-chars
JWT_EXPIRY=24h

CORS_ALLOWED_ORIGINS=https://your-domain.com
```

---

## 🗄️ Database Management

### View Database
```bash
# Connect to PostgreSQL
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager

# List all tables
\dt

# Describe table
\d developers

# Exit
\q
```

### Create Migration
```bash
# Create new migration
migrate create -ext sql -dir backend/migrations -seq migration_name

# Example: Add email verification
migrate create -ext sql -dir backend/migrations -seq add_email_verification
```

### Rollback Migration
```bash
# Rollback last migration
make migrate-down

# Or
migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" down 1
```

### Reset Database
```bash
# Stop containers
docker-compose down

# Remove volumes (WARNING: deletes all data)
docker-compose down -v

# Start fresh
docker-compose up -d
make migrate-up
```

---

## 🔧 Development

### Run in Development Mode

**Backend (with auto-reload):**
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd backend
air
```

**Frontend (with Turbopack):**
```bash
cd frontend
npm run dev
```

### Build for Production

**Backend:**
```bash
cd backend
go build -o bin/server cmd/server/main.go
./bin/server
```

**Frontend:**
```bash
cd frontend
npm run build
npm start
```

---

## 🐳 Docker Deployment

### Build Docker Images

**Backend Dockerfile:**
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY backend/go.* ./
RUN go mod download
COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
EXPOSE 8081
CMD ["./server"]
```

**Frontend Dockerfile:**
```dockerfile
FROM node:20-alpine AS builder
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM node:20-alpine
WORKDIR /app
COPY --from=builder /app/.next/standalone ./
COPY --from=builder /app/.next/static ./.next/static
COPY --from=builder /app/public ./public
EXPOSE 3000
CMD ["node", "server.js"]
```

**Build & Run:**
```bash
# Build
docker build -t taskmanager-backend ./backend
docker build -t taskmanager-frontend ./frontend

# Run
docker run -p 8081:8081 taskmanager-backend
docker run -p 3000:3000 taskmanager-frontend
```

---

## 🧪 Testing

### Backend Tests
```bash
cd backend
go test ./... -v
```

### Frontend Tests
```bash
cd frontend
npm run test
```

---

## ❓ Troubleshooting

### Port Already in Use
```bash
# Find process using port
lsof -i :8081
lsof -i :3000

# Kill process
kill -9 <PID>
```

### Database Connection Failed
```bash
# Check if PostgreSQL is running
docker ps | grep taskmanager-db

# Check logs
docker logs taskmanager-db

# Restart
docker-compose restart db
```

### Redis Connection Failed
```bash
# Check if Redis is running
docker ps | grep taskmanager-redis

# Test connection
docker exec -it taskmanager-redis redis-cli ping
# Expected: PONG
```

### Migration Failed
```bash
# Check migration version
migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" version

# Force version (if stuck)
migrate -path backend/migrations -database "..." force <version>
```

### Frontend Build Errors
```bash
# Clear Next.js cache
cd frontend
rm -rf .next node_modules
npm install
npm run dev
```

---

## 📚 Next Steps

After installation:

1. **Create Account:** http://localhost:3000/register
2. **Login:** http://localhost:3000/login
3. **Explore Dashboard:** Create tasks, projects, and team
4. **Read API Docs:** `docs/API.md`
5. **Check Database:** Verify data is being stored

---

## 🆘 Getting Help

- **GitHub Issues:** https://github.com/ardani17/task-manager/issues
- **Documentation:** Check `docs/` folder
- **Logs:** Check server logs for errors

---

## 📝 Checklist

After installation, verify:

- [ ] PostgreSQL running on port 5433
- [ ] Redis running on port 6380
- [ ] Backend API accessible at http://localhost:8081/health
- [ ] Frontend accessible at http://localhost:3000
- [ ] Can register new user
- [ ] Can login successfully
- [ ] Dashboard loads with no errors

---

**Installation complete! 🎉**

For production deployment, see `docs/DEPLOYMENT.md` (if available).
