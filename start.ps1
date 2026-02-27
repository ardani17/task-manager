# start.ps1 - PowerShell script for Windows users
# Alternative to Makefile commands

param(
    [Parameter(Position=0)]
    [string]$Action = "help"
)

$ErrorActionPreference = "Stop"

function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

switch ($Action) {
    "help" {
        Write-ColorOutput Green "Task Manager - Available Commands:"
        Write-Output ""
        Write-Output "Docker:"
        Write-Output "  .\start.ps1 docker-up        Start all services (db, redis, backend, frontend)"
        Write-Output "  .\start.ps1 docker-down      Stop all services"
        Write-Output "  .\start.ps1 docker-logs      View logs"
        Write-Output "  .\start.ps1 docker-build     Build all Docker images"
        Write-Output "  .\start.ps1 docker-restart   Restart all services"
        Write-Output ""
        Write-Output "Database:"
        Write-Output "  .\start.ps1 db-up            Start database (PostgreSQL + Redis)"
        Write-Output "  .\start.ps1 db-down          Stop database"
        Write-Output "  .\start.ps1 migrate-up       Run database migrations"
        Write-Output "  .\start.ps1 migrate-down     Rollback migrations"
        Write-Output ""
        Write-Output "Development:"
        Write-Output "  .\start.ps1 backend          Run backend server"
        Write-Output "  .\start.ps1 frontend         Run frontend dev server"
        Write-Output ""
        Write-Output "Utility:"
        Write-Output "  .\start.ps1 clean            Clean build files"
        Write-Output "  .\start.ps1 help             Show this help message"
    }

    "docker-up" {
        Write-ColorOutput Cyan "Starting all services..."
        docker-compose up -d
        Write-ColorOutput Green "✅ All services started"
        Write-Output "📊 Backend: http://localhost:8081"
        Write-Output "🌐 Frontend: http://localhost:3000"
        Write-Output "🗄️ PostgreSQL: localhost:5433"
        Write-Output "🔴 Redis: localhost:6380"
    }

    "docker-down" {
        Write-ColorOutput Cyan "Stopping all services..."
        docker-compose down
        Write-ColorOutput Green "✅ All services stopped"
    }

    "docker-logs" {
        docker-compose logs -f
    }

    "docker-build" {
        Write-ColorOutput Cyan "Building all Docker images..."
        docker-compose build
        Write-ColorOutput Green "✅ All Docker images built"
    }

    "docker-restart" {
        Write-ColorOutput Cyan "Restarting all services..."
        docker-compose restart
        Write-ColorOutput Green "✅ All services restarted"
    }

    "db-up" {
        Write-ColorOutput Cyan "Starting database..."
        docker-compose up -d db redis
        Write-ColorOutput Green "✅ Database started"
        Write-Output "📊 PostgreSQL: localhost:5433"
        Write-Output "🔴 Redis: localhost:6380"
    }

    "db-down" {
        Write-ColorOutput Cyan "Stopping database..."
        docker-compose stop db redis
        Write-ColorOutput Green "✅ Database stopped"
    }

    "migrate-up" {
        Write-ColorOutput Cyan "Running migrations..."

        # Check if migrate tool is available
        $migrateAvailable = Get-Command migrate -ErrorAction SilentlyContinue

        if ($migrateAvailable) {
            migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" up
        } else {
            Write-ColorOutput Yellow "migrate tool not found, running SQL files manually..."

            $migrations = @(
                "001_developers.up.sql",
                "002_tasks.up.sql",
                "003_projects.up.sql",
                "004_activity_logs.up.sql",
                "005_add_password_hash.up.sql"
            )

            foreach ($migration in $migrations) {
                $path = "backend/migrations/$migration"
                if (Test-Path $path) {
                    Write-Output "Running $migration..."
                    Get-Content $path | docker exec -i taskmanager-db psql -U taskmanager -d taskmanager 2>&1 | Out-Null
                }
            }
        }

        Write-ColorOutput Green "✅ Migrations completed"
    }

    "migrate-down" {
        Write-ColorOutput Cyan "Rolling back migrations..."

        $migrateAvailable = Get-Command migrate -ErrorAction SilentlyContinue

        if ($migrateAvailable) {
            migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" down
            Write-ColorOutput Green "✅ Rollback completed"
        } else {
            Write-ColorOutput Red "❌ migrate tool required for rollback. Please install golang-migrate."
        }
    }

    "backend" {
        Write-ColorOutput Cyan "Starting backend server..."
        Push-Location backend
        try {
            go run cmd/server/main.go
        } finally {
            Pop-Location
        }
    }

    "frontend" {
        Write-ColorOutput Cyan "Starting frontend server..."
        Push-Location frontend
        try {
            npm run dev
        } finally {
            Pop-Location
        }
    }

    "clean" {
        Write-ColorOutput Cyan "Cleaning build files..."
        if (Test-Path "backend\tmp") { Remove-Item -Recurse -Force "backend\tmp" }
        if (Test-Path "frontend\.next") { Remove-Item -Recurse -Force "frontend\.next" }
        if (Test-Path "frontend\node_modules") { Remove-Item -Recurse -Force "frontend\node_modules" }
        Write-ColorOutput Green "✅ Cleaned build files"
    }

    default {
        Write-ColorOutput Red "❌ Unknown command: $Action"
        Write-Output "Run '.\start.ps1 help' for available commands"
        exit 1
    }
}
