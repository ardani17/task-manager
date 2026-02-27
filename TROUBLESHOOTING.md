# Troubleshooting Guide

Common issues and solutions when setting up Task Manager.

---

## ❌ Error: Database Connection Failed

### Error Message
```
Failed to connect to database error="failed to ping database: dial tcp [::1]:5433: connectex: No connection could be made because the target machine actively refused it."
```

### Problem
Backend trying to connect to port 5433, but PostgreSQL is running on port 5432 (Laragon default).

### Solution

**Option 1: Update .env file (Recommended)**

```powershell
# Edit .env file in project root
notepad .env

# Change DB_PORT from 5433 to 5432:
DB_PORT=5432

# Save and run backend again
go run cmd/server/main.go
```

**Option 2: Check which port PostgreSQL is running**

```powershell
# In Laragon, check PostgreSQL port
# Usually in: C:\laragon\bin\postgresql\postgresql.conf

# Or query PostgreSQL directly
psql -U postgres -c "SHOW port;"

# Then update .env accordingly
```

**Option 3: Use Docker PostgreSQL on 5433**

```powershell
# Start Docker PostgreSQL on port 5433
docker run -d -p 5433:5432 -e POSTGRES_USER=taskmanager -e POSTGRES_PASSWORD=taskmanager123 -e POSTGRES_DB=taskmanager --name taskmanager-db postgres:16-alpine

# Run migrations
Get-Content backend/migrations/001_developers.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
# ... repeat for all migrations
```

---

## ❌ Error: Path Not Found

### Error Message
```
Get-Content : Cannot find path 'C:\laragon\www\task-manager\frontend\backend\migrations\001_developers.up.sql'
```

### Problem
Running command from wrong directory (inside `frontend/` folder).

### Solution

```powershell
# Go back to project root first
cd C:\laragon\www\task-manager

# Then run migrations
Get-Content backend/migrations/001_developers.up.sql | psql -U taskmanager -d taskmanager
```

**Correct workflow:**
```powershell
# Always from project root
cd C:\laragon\www\task-manager

# Run migrations
Get-Content backend/migrations/001_developers.up.sql | psql -U taskmanager -d taskmanager
Get-Content backend/migrations/002_tasks.up.sql | psql -U taskmanager -d taskmanager
# ... etc

# Then start backend
cd backend
go run cmd/server/main.go
```

---

## ❌ Error: Redis Connection Failed

### Error Message
```
Failed to connect to Redis
```

### Problem
Redis not running or not installed.

### Solutions

**Option 1: Install Redis via Docker (Easiest)**
```powershell
docker run -d -p 6379:6379 --name redis redis:7-alpine
```

**Option 2: Install Memurai (Redis for Windows)**
```powershell
# Download from: https://www.memurai.com/
# Install and start service
memurai-cli ping
```

**Option 3: Disable Redis Temporarily (Development Only)**

Edit `backend/internal/config/config.go`:
```go
// Comment out Redis connection
// cfg.Redis.Host = os.Getenv("REDIS_HOST")
// cfg.Redis.Port = os.Getenv("REDIS_PORT")
```

**Option 4: Use Redis in Docker Compose**
```powershell
# Start only Redis from docker-compose
docker-compose up -d redis
```

---

## ❌ Error: Port Already in Use

### Error Message
```
bind: address already in use
```

### Problem
Another service using port 8080 or 3000.

### Solution

**Find and kill process:**
```powershell
# Find process on port 8080
netstat -ano | findstr :8080

# Kill process (replace PID)
taskkill /PID <PID> /F
```

**Or change port:**
```powershell
# Edit .env
APP_PORT=8081

# Restart backend
go run cmd/server/main.go
```

---

## ❌ Error: Authentication Failed

### Error Message
```
psql: error: connection to server at "localhost", port 5432 failed: FATAL: password authentication failed for user "taskmanager"
```

### Problem
User doesn't exist or password incorrect.

### Solution

```powershell
# Recreate user
psql -U postgres
DROP USER IF EXISTS taskmanager;
CREATE USER taskmanager WITH PASSWORD 'taskmanager123';
GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;
\q
```

---

## ❌ Error: Migration Already Exists

### Error Message
```
ERROR: relation "developers" already exists
```

### Problem
Migration already ran before.

### Solution

**Option 1: Ignore (it's okay)**
The table already exists, migration succeeded.

**Option 2: Reset database**
```powershell
# Drop and recreate
psql -U postgres
DROP DATABASE taskmanager;
CREATE DATABASE taskmanager;
GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;
\q

# Run migrations again
cd C:\laragon\www\task-manager
Get-Content backend/migrations/001_developers.up.sql | psql -U taskmanager -d taskmanager
# ... etc
```

---

## ✅ Verification Checklist

After setup, verify everything works:

### 1. Database
```powershell
# Connect to database
psql -U taskmanager -d taskmanager

# List tables
\dt

# Should show:
#  developers
#  tasks
#  projects
#  activities
#  teams
#  schema_migrations

# Exit
\q
```

### 2. Redis (if using)
```powershell
# Test connection
redis-cli ping
# Should return: PONG
```

### 3. Backend
```powershell
# Start backend
cd C:\laragon\www\task-manager\backend
go run cmd/server/main.go

# Should see:
# ✅ Connected to database
# ✅ Redis connected
# 🚀 Server starting on :8080

# Test health endpoint (new terminal)
curl http://localhost:8080/health
# Should return: {"status":"healthy"}
```

### 4. Frontend
```powershell
# Start frontend
cd C:\laragon\www\task-manager\frontend
npm run dev

# Should see:
# ✓ Ready in 2.5s
# Local: http://localhost:3000

# Open browser
# http://localhost:3000
```

---

## 🔧 Quick Fix Commands

### Complete Reset
```powershell
# Reset database
psql -U postgres -c "DROP DATABASE IF EXISTS taskmanager;"
psql -U postgres -c "CREATE DATABASE taskmanager;"
psql -U postgres -c "CREATE USER taskmanager WITH PASSWORD 'taskmanager123';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE taskmanager TO taskmanager;"

# Run migrations
cd C:\laragon\www\task-manager
Get-Content backend/migrations/*.sql | psql -U taskmanager -d taskmanager
```

### Check Current Config
```powershell
# Check .env file
cat .env

# Check PostgreSQL port
psql -U postgres -c "SHOW port;"

# Check if Redis running
redis-cli ping
```

### Start Fresh
```powershell
# Pull latest code
git pull

# Reinstall dependencies
cd backend
go mod download

cd ../frontend
npm install

# Setup .env
cp .env.example .env
# Edit .env with correct DB_PORT (5432 for Laragon)

# Run migrations
cd ..
Get-Content backend/migrations/*.sql | psql -U taskmanager -d taskmanager

# Start services
# Terminal 1: Backend
cd backend && go run cmd/server/main.go

# Terminal 2: Frontend
cd frontend && npm run dev
```

---

## 📞 Still Having Issues?

1. **Check logs carefully** - Read error message
2. **Verify services running** - PostgreSQL, Redis
3. **Check .env file** - Correct ports and credentials
4. **Try Docker method** - `docker-compose up -d` (easier)
5. **Check documentation** - MANUAL.md, DOCKER.md

---

**Most Common Fix: Change DB_PORT to 5432 in .env file!**
