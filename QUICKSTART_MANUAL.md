# Quick Start (Without Docker)

Simplified guide for running Task Manager locally.

---

## 🎯 Prerequisites (Laragon)

Laragon already has:
- ✅ PostgreSQL
- ✅ Apache/Nginx
- ✅ PHP (not needed for this project)

You still need:
- ⬜ Go 1.22+ - https://golang.org/dl/
- ⬜ Node.js 20+ - https://nodejs.org/
- ⬜ Redis - Use Memurai or Docker

---

## 🚀 3-Step Setup

### Step 1: Database

```powershell
# Open Laragon, start PostgreSQL

# Create database (Terminal)
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

### Step 2: Backend

```powershell
# Terminal 1
cd C:\laragon\www\task-manager\backend
go mod download
go run cmd/server/main.go
```

### Step 3: Frontend

```powershell
# Terminal 2
cd C:\laragon\www\task-manager\frontend
npm install
npm run dev
```

**Access:** http://localhost:3000

---

## 🔧 Redis Setup (Choose One)

### Option A: Docker (Easiest)
```powershell
docker run -d -p 6379:6379 --name redis redis:7-alpine
```

### Option B: Memurai
```powershell
# Download: https://www.memurai.com/
# Install and start service
memurai-cli ping
```

### Option C: Disable Redis (Dev Only)
Comment out Redis code in backend temporarily.

---

## ⚙️ Environment Files

### `.env` (project root)
```bash
APP_PORT=8081
DB_HOST=localhost
DB_PORT=5432
DB_USER=taskmanager
DB_PASS=taskmanager123
DB_NAME=taskmanager
REDIS_HOST=localhost
REDIS_PORT=6379
JWT_SECRET=dev-secret-key
```

### `frontend/.env.local`
```bash
NEXT_PUBLIC_API_URL=http://localhost:8081/api/v1
```

---

## ✅ Verify

```powershell
# Backend
curl http://localhost:8081/health
# {"status":"healthy"}

# Frontend
# Open browser: http://localhost:3000
```

---

## 🔄 Daily Commands

```powershell
# Start backend
cd backend
go run cmd/server/main.go

# Start frontend
cd frontend
npm run dev

# Check database
psql -U taskmanager -d taskmanager -c "\dt"
```

---

**That's it! 🎉**
