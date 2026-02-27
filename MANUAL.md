# Manual Installation (Without Docker)

Complete guide to run Task Manager without Docker containers.

---

## 📋 Prerequisites

### Required Software

| Software | Version | Download |
|----------|---------|----------|
| **Go** | 1.22+ | https://golang.org/dl/ |
| **Node.js** | 20+ | https://nodejs.org/ |
| **PostgreSQL** | 16+ | https://www.postgresql.org/download/ |
| **Redis** | 7+ | https://redis.io/download |

### Check Installations

```bash
go version          # Should show Go 1.22+
node --version      # Should show Node 20+
npm --version       # Should show npm 10+
psql --version      # Should show PostgreSQL 16+
redis-cli --version # Should show Redis 7+
```

---

## 🚀 Quick Start (Manual)

### 1. Install Dependencies

**Windows (Laragon):**
- Laragon already includes PostgreSQL
- Install Go: https://golang.org/dl/
- Install Node.js: https://nodejs.org/
- Install Redis: Enable in Laragon or use Memurai

**Linux/macOS:**
```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go nodejs npm postgresql redis-server

# macOS (Homebrew)
brew install go node postgresql redis
```

---

## 📦 Setup Database

### Option A: Using Laragon (Windows)

1. **Start PostgreSQL in Laragon**
   - Open Laragon
   - Click "Start All"
   - PostgreSQL should be running

2. **Create Database**
   ```powershell
   # Open HeidiSQL or terminal
   psql -U postgres -c "CREATE DATABASE taskmanager;"
   psql -U postgres -c "CREATE USER taskmanager WITH PASSWORD 'taskmanager123';"
   psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;"
   ```

3. **Run Migrations**
   ```powershell
   # Using psql directly
   psql -U taskmanager -d taskmanager -f backend/migrations/001_developers.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/002_tasks.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/003_projects.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/004_activity_logs.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/005_add_password_hash.up.sql
   ```

### Option B: Fresh PostgreSQL Installation

1. **Start PostgreSQL**
   ```bash
   # Linux
   sudo systemctl start postgresql

   # macOS
   brew services start postgresql

   # Windows (service should auto-start)
   ```

2. **Create Database & User**
   ```bash
   # Connect as postgres user
   sudo -u postgres psql

   # In psql console:
   CREATE DATABASE taskmanager;
   CREATE USER taskmanager WITH PASSWORD 'taskmanager123';
   GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;
   \q
   ```

3. **Run Migrations**
   ```bash
   # Using migrate tool (if installed)
   migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5432/taskmanager?sslmode=disable" up

   # Or manually
   psql -U taskmanager -d taskmanager -f backend/migrations/001_developers.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/002_tasks.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/003_projects.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/004_activity_logs.up.sql
   psql -U taskmanager -d taskmanager -f backend/migrations/005_add_password_hash.up.sql
   ```

---

## 🔧 Setup Redis

### Windows (Laragon/Memurai)

1. **Using Laragon:**
   - Redis might not be included by default
   - Enable Redis extension if available

2. **Using Memurai (Redis for Windows):**
   ```powershell
   # Download: https://www.memurai.com/
   # Install and start service
   memurai-cli ping  # Should return PONG
   ```

3. **Or use Docker for Redis only:**
   ```powershell
   docker run -d -p 6379:6379 --name redis redis:7-alpine
   ```

### Linux/macOS

```bash
# Start Redis
sudo systemctl start redis  # Linux
brew services start redis   # macOS

# Test connection
redis-cli ping  # Should return PONG
```

---

## ⚙️ Configuration

### Backend Environment

Create `.env` in project root:

```bash
# Server Configuration
APP_PORT=8081
APP_ENV=development

# Database Configuration (Laragon default)
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskmanager
DB_PASS=taskmanager123
DB_NAME=taskmanager

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379

# JWT Configuration
JWT_SECRET=dev-secret-key-change-in-production
JWT_EXPIRY=24h

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://127.0.0.1:3000
```

### Frontend Environment

Create `frontend/.env.local`:

```bash
# API Configuration
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

---

## 🚀 Running the Application

### Step 1: Setup Database

```powershell
# Run migrations
cd C:\laragon\www\task-manager

# Method 1: Using psql
psql -U taskmanager -d taskmanager -f backend/migrations/001_developers.up.sql
psql -U taskmanager -d taskmanager -f backend/migrations/002_tasks.up.sql
psql -U taskmanager -d taskmanager -f backend/migrations/003_projects.up.sql
psql -U taskmanager -d taskmanager -f backend/migrations/004_activity_logs.up.sql
psql -U taskmanager -d taskmanager -f backend/migrations/005_add_password_hash.up.sql

# Method 2: Using PowerShell script
.\start.ps1 migrate-up
```

### Step 2: Start Backend

**Terminal 1 (Backend):**
```powershell
cd C:\laragon\www\task-manager\backend

# Download dependencies
go mod download

# Run server
go run cmd/server/main.go
```

Expected output:
```
2026/02/27 15:00:00 🚀 Server starting on :8081
2026/02/27 15:00:00 ✅ Connected to database
2026/02/27 15:00:00 ✅ Redis connected
```

**Verify backend:**
```powershell
curl http://localhost:8081/health
# Or in browser: http://localhost:8081/health
```

### Step 3: Start Frontend

**Terminal 2 (Frontend):**
```powershell
cd C:\laragon\www\task-manager\frontend

# Install dependencies
npm install

# Run development server
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

### Step 4: Access Application

- **Frontend:** http://localhost:3000
- **Backend API:** http://localhost:8081
- **Health Check:** http://localhost:8081/health

---

## 🔍 Complete Workflow (Manual)

### Terminal 1: Database Setup
```powershell
# Create database (one-time)
psql -U postgres
CREATE DATABASE taskmanager;
CREATE USER taskmanager WITH PASSWORD 'taskmanager123';
GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;
\q

# Run migrations
cd C:\laragon\www\task-manager
Get-Content backend/migrations/001_developers.up.sql | psql -U taskmanager -d taskmanager
Get-Content backend/migrations/002_tasks.up.sql | psql -U taskmanager -d taskmanager
Get-Content backend/migrations/003_projects.up.sql | psql -U taskmanager -d taskmanager
Get-Content backend/migrations/004_activity_logs.up.sql | psql -U taskmanager -d taskmanager
Get-Content backend/migrations/005_add_password_hash.up.sql | psql -U taskmanager -d taskmanager
```

### Terminal 2: Backend
```powershell
cd C:\laragon\www\task-manager\backend
go mod download
go run cmd/server/main.go
```

### Terminal 3: Frontend
```powershell
cd C:\laragon\www\task-manager\frontend
npm install
npm run dev
```

### Terminal 4: Redis (if needed)
```powershell
# Start Redis service
# If using Memurai:
memurai-server

# If using Docker:
docker run -d -p 6379:6379 redis:7-alpine
```

---

## 🐛 Troubleshooting

### Database Connection Failed

```powershell
# Check PostgreSQL is running
# Laragon: Check if PostgreSQL is green/started

# Test connection
psql -U taskmanager -d taskmanager -c "SELECT version();"

# If password fails, reset it
psql -U postgres -c "ALTER USER taskmanager WITH PASSWORD 'taskmanager123';"
```

### Redis Connection Failed

```powershell
# Check if Redis is running
redis-cli ping
# Should return: PONG

# If not running, start it
# Windows (Memurai):
memurai-server

# Or use Docker:
docker run -d -p 6379:6379 redis:7-alpine
```

### Port Already in Use

```powershell
# Find process on port 8081
netstat -ano | findstr :8081

# Kill process (replace PID)
taskkill /PID <PID> /F

# Or change port in .env
APP_PORT=8082
```

### Go Dependencies Issue

```powershell
cd backend
# Clean module cache
go clean -modcache

# Download fresh
go mod download
go mod tidy
```

### Frontend Build Errors

```powershell
cd frontend
# Clean and reinstall
rm -r node_modules
rm package-lock.json
npm install
```

---

## 📊 Database Management

### View Data
```powershell
# Connect to database
psql -U taskmanager -d taskmanager

# List tables
\dt

# View developers
SELECT * FROM developers;

# View tasks
SELECT * FROM tasks;

# Exit
\q
```

### Reset Database
```powershell
# Drop and recreate
psql -U postgres -c "DROP DATABASE taskmanager;"
psql -U postgres -c "CREATE DATABASE taskmanager;"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;"

# Run migrations again
cd C:\laragon\www\task-manager
.\start.ps1 migrate-up
```

---

## 🔄 Development Workflow

### Daily Development

```powershell
# Terminal 1: Start backend
cd backend
go run cmd/server/main.go

# Terminal 2: Start frontend
cd frontend
npm run dev

# Terminal 3: Database (if needed)
psql -U taskmanager -d taskmanager
```

### Hot Reload (Optional)

**Backend (using Air):**
```powershell
# Install Air
go install github.com/cosmtrek/air@latest

# Run with hot reload
cd backend
air
```

**Frontend:**
```powershell
# Already has hot reload with npm run dev
cd frontend
npm run dev
```

---

## 🚀 Production Build (Manual)

### Backend
```powershell
cd backend

# Build binary
go build -o bin/server.exe cmd/server/main.go

# Run
.\bin\server.exe
```

### Frontend
```powershell
cd frontend

# Build
npm run build

# Run production
npm start
```

---

## 📝 Environment Checklist

- [ ] PostgreSQL running (port 5432)
- [ ] Redis running (port 6379)
- [ ] Database created
- [ ] Migrations ran
- [ ] `.env` file created
- [ ] Backend dependencies installed
- [ ] Frontend dependencies installed
- [ ] Backend running (port 8081)
- [ ] Frontend running (port 3000)

---

## 🎯 Summary

**Without Docker, you need:**
1. ✅ PostgreSQL installed & running
2. ✅ Redis installed & running
3. ✅ Go 1.22+ installed
4. ✅ Node.js 20+ installed
5. ✅ Database created & migrated
6. ✅ `.env` files configured
7. ✅ Backend running on port 8081
8. ✅ Frontend running on port 3000

**Commands:**
```powershell
# Setup (one-time)
psql -U postgres -c "CREATE DATABASE taskmanager;"
.\start.ps1 migrate-up

# Run backend
cd backend && go run cmd/server/main.go

# Run frontend
cd frontend && npm run dev
```

---

**Manual installation complete! 🚀**
