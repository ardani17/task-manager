# Docker Commands for Windows (PowerShell)

Use these commands if you don't have `make` installed.

---

## 🚀 Quick Start (Without Make)

### Start Database
```powershell
docker-compose up -d db redis
```

### Run Migrations
```powershell
# If you have migrate tool installed:
migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" up

# Or manually via Docker:
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/001_developers.up.sql
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/002_tasks.up.sql
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/003_projects.up.sql
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/004_activity_logs.up.sql
docker exec -i taskmanager-db psql -U taskmanager -d taskmanager < backend/migrations/005_add_password_hash.up.sql
```

### Start Backend
```powershell
cd backend
go run cmd/server/main.go
```

### Start Frontend
```powershell
cd frontend
npm run dev
```

---

## 🐳 Docker Compose Commands (Make Alternative)

### All Services
```powershell
# Start all
docker-compose up -d

# Stop all
docker-compose down

# View logs
docker-compose logs -f

# Restart all
docker-compose restart

# Build all
docker-compose build

# Rebuild and start
docker-compose up -d --build
```

### Database Only
```powershell
# Start database
docker-compose up -d db redis

# Stop database
docker-compose stop db redis

# View database logs
docker-compose logs -f db

# Connect to PostgreSQL
docker exec -it taskmanager-db psql -U taskmanager -d taskmanager

# Connect to Redis
docker exec -it taskmanager-redis redis-cli
```

### Backend Only
```powershell
# Start backend (with dependencies)
docker-compose up -d backend

# View backend logs
docker-compose logs -f backend

# Restart backend
docker-compose restart backend

# Rebuild backend
docker-compose build backend
docker-compose up -d backend
```

### Frontend Only
```powershell
# Start frontend (with dependencies)
docker-compose up -d frontend

# View frontend logs
docker-compose logs -f frontend

# Restart frontend
docker-compose restart frontend

# Rebuild frontend
docker-compose build frontend
docker-compose up -d frontend
```

---

## 📋 Complete Workflow (Windows)

### 1. Start Database
```powershell
# From project root
docker-compose up -d db redis

# Wait for database to be ready (10-15 seconds)
Start-Sleep -Seconds 15

# Verify database is running
docker-compose ps
```

### 2. Run Migrations
```powershell
# Option A: Using migrate tool (if installed)
migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" up

# Option B: Manual SQL execution
Get-Content backend/migrations/001_developers.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
Get-Content backend/migrations/002_tasks.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
Get-Content backend/migrations/003_projects.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
Get-Content backend/migrations/004_activity_logs.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
Get-Content backend/migrations/005_add_password_hash.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
```

### 3. Start Backend
```powershell
# Terminal 1
cd backend
$env:DB_HOST="localhost"
$env:DB_PORT="5433"
$env:REDIS_HOST="localhost"
$env:REDIS_PORT="6380"
go run cmd/server/main.go
```

### 4. Start Frontend
```powershell
# Terminal 2
cd frontend
npm run dev
```

---

## 🛠️ Install Make on Windows (Optional)

### Option 1: Chocolatey (Recommended)
```powershell
# Install Chocolatey (if not installed)
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))

# Install make
choco install make -y

# Verify
make --version
```

### Option 2: Scoop
```powershell
# Install Scoop (if not installed)
Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
irm get.scoop.sh | iex

# Install make
scoop install make

# Verify
make --version
```

### Option 3: WSL (Windows Subsystem for Linux)
```powershell
# Install WSL
wsl --install

# Restart Windows, then in WSL terminal:
sudo apt update
sudo apt install make

# Use make in WSL
cd /mnt/c/laragon/www/task-manager
make db-up
```

### Option 4: Git Bash (Comes with Git for Windows)
```bash
# Git Bash already has make-like functionality
# Use the docker-compose commands directly instead

# In Git Bash:
docker-compose up -d db redis
```

---

## 📝 PowerShell Script Alternative

Create `start.ps1` in project root:

```powershell
# start.ps1 - PowerShell script for common tasks

param(
    [Parameter(Position=0)]
    [string]$Action = "help"
)

switch ($Action) {
    "db-up" {
        Write-Host "Starting database..."
        docker-compose up -d db redis
        Write-Host "✅ Database started"
    }
    "db-down" {
        Write-Host "Stopping database..."
        docker-compose stop db redis
        Write-Host "✅ Database stopped"
    }
    "migrate-up" {
        Write-Host "Running migrations..."
        Get-Content backend/migrations/001_developers.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
        Get-Content backend/migrations/002_tasks.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
        Get-Content backend/migrations/003_projects.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
        Get-Content backend/migrations/004_activity_logs.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
        Get-Content backend/migrations/005_add_password_hash.up.sql | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager
        Write-Host "✅ Migrations completed"
    }
    "docker-up" {
        Write-Host "Starting all services..."
        docker-compose up -d
        Write-Host "✅ All services started"
    }
    "docker-down" {
        Write-Host "Stopping all services..."
        docker-compose down
        Write-Host "✅ All services stopped"
    }
    "backend" {
        Write-Host "Starting backend..."
        cd backend
        go run cmd/server/main.go
    }
    "frontend" {
        Write-Host "Starting frontend..."
        cd frontend
        npm run dev
    }
    default {
        Write-Host "Task Manager - Available Commands:"
        Write-Host ""
        Write-Host "  .\start.ps1 db-up        Start database"
        Write-Host "  .\start.ps1 db-down      Stop database"
        Write-Host "  .\start.ps1 migrate-up   Run migrations"
        Write-Host "  .\start.ps1 docker-up    Start all services"
        Write-Host "  .\start.ps1 docker-down  Stop all services"
        Write-Host "  .\start.ps1 backend      Run backend server"
        Write-Host "  .\start.ps1 frontend     Run frontend server"
    }
}
```

**Usage:**
```powershell
.\start.ps1 db-up
.\start.ps1 migrate-up
.\start.ps1 docker-up
```

---

## 🎯 Recommended Approach for Windows

### Easiest: Direct Docker Compose
```powershell
# Just use docker-compose directly
docker-compose up -d          # Start all
docker-compose down           # Stop all
docker-compose logs -f        # View logs
```

### Best Long-term: Install Chocolatey + Make
```powershell
# One-time setup
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
choco install make -y

# Then use make commands
make docker-up
make migrate-up
```

---

## ✅ Quick Start for Windows

```powershell
# 1. Start all services (simplest)
docker-compose up -d

# 2. Wait for services to be ready
Start-Sleep -Seconds 20

# 3. Check status
docker-compose ps

# 4. Access
# Frontend: http://localhost:3000
# Backend:  http://localhost:8081
```

---

**Windows commands ready! 🪟**
