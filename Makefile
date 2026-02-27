.PHONY: help docker-up docker-down db-up db-down migrate-up migrate-down backend frontend clean docker-build docker-run

# Default target
help:
	@echo "Task Manager - Available Commands:"
	@echo ""
	@echo "Docker:"
	@echo "  make docker-up        Start all services (db, redis, backend, frontend)"
	@echo "  make docker-down      Stop all services"
	@echo "  make docker-logs      View logs"
	@echo "  make docker-build     Build all Docker images"
	@echo "  make docker-restart   Restart all services"
	@echo ""
	@echo "Database:"
	@echo "  make db-up            Start database (PostgreSQL + Redis)"
	@echo "  make db-down          Stop database"
	@echo "  make migrate-up       Run database migrations"
	@echo "  make migrate-down     Rollback migrations"
	@echo ""
	@echo "Development:"
	@echo "  make backend          Run backend server"
	@echo "  make frontend         Run frontend dev server"
	@echo ""
	@echo "Utility:"
	@echo "  make clean            Clean build files"
	@echo "  make help             Show this help message"

# Docker commands
docker-up:
	docker-compose up -d
	@echo "✅ All services started"
	@echo "📊 Backend: http://localhost:8081"
	@echo "🌐 Frontend: http://localhost:3000"
	@echo "🗄️ PostgreSQL: localhost:5433"
	@echo "🔴 Redis: localhost:6380"

docker-down:
	docker-compose down
	@echo "✅ All services stopped"

docker-logs:
	docker-compose logs -f

docker-build:
	docker-compose build
	@echo "✅ All Docker images built"

docker-restart:
	docker-compose restart
	@echo "✅ All services restarted"

# Database commands
db-up:
	docker-compose up -d db redis
	@echo "✅ Database started"
	@echo "📊 PostgreSQL: localhost:5433"
	@echo "🔴 Redis: localhost:6380"

db-down:
	docker-compose down db redis
	@echo "✅ Database stopped"

# Migrations
migrate-up:
	@echo "Running migrations..."
	migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" up
	@echo "✅ Migrations completed"

migrate-down:
	@echo "Rolling back migrations..."
	migrate -path backend/migrations -database "postgresql://taskmanager:taskmanager123@localhost:5433/taskmanager?sslmode=disable" down
	@echo "✅ Rollback completed"

# Development commands
backend:
	cd backend && go run cmd/server/main.go

frontend:
	cd frontend && npm run dev

# Utility commands
clean:
	rm -rf backend/tmp
	rm -rf frontend/.next
	rm -rf frontend/node_modules
	@echo "✅ Cleaned build files"
