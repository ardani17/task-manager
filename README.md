# Task Manager - Developer Progress Monitoring

## 🎯 Project Overview

**Purpose:** Web application untuk monitoring progress developer tim  
**Type:** Real-time Dashboard + Task Management  
**Timeline:** 6 Days

---

## 🚀 Features

### Core Features
- ✅ Developer activity tracking
- ✅ Task assignment & progress monitoring
- ✅ Team performance dashboard
- ✅ Real-time updates
- ✅ Activity timeline
- ✅ Statistics & analytics

### Dashboard Views
- **Team Overview** - Semua developer status
- **Individual Progress** - Per-developer detail
- **Task Board** - Kanban view
- **Activity Feed** - Real-time timeline
- **Analytics** - Charts & metrics

---

## 🛠️ Tech Stack

| Layer | Technology |
|-------|------------|
| **Backend** | Go + Chi router |
| **Frontend** | Next.js 16 + React 19 |
| **Database** | PostgreSQL 16 |
| **Cache/Queue** | Redis 7 |
| **Real-time** | WebSocket |
| **UI** | Tailwind CSS + shadcn/ui |

---

## 📊 Database Models

### Developers
```go
type Developer struct {
    ID          uuid.UUID
    Name        string
    Email       string
    Role        string
    Avatar      string
    Status      string // active, idle, offline
    LastActive  time.Time
    CreatedAt   time.Time
}
```

### Tasks
```go
type Task struct {
    ID          uuid.UUID
    Title       string
    Description string
    Assignee    uuid.UUID
    Status      string // todo, in_progress, review, done
    Priority    string // low, medium, high
    Project     string
    DueDate     *time.Time
    CreatedAt   time.Time
    UpdatedAt   time.Time
}
```

### Activities
```go
type Activity struct {
    ID          uuid.UUID
    DeveloperID uuid.UUID
    ActionType  string // commit, task_done, review, deploy
    Description string
    Metadata    JSONB
    CreatedAt   time.Time
}
```

### Projects
```go
type Project struct {
    ID          uuid.UUID
    Name        string
    Description string
    Status      string // active, completed, paused
    Members     []uuid.UUID
    CreatedAt   time.Time
}
```

---

## 🔌 API Endpoints

### Developers
```
GET    /api/v1/developers           # List all developers
GET    /api/v1/developers/:id       # Get developer detail
GET    /api/v1/developers/:id/stats # Developer statistics
GET    /api/v1/developers/:id/tasks # Developer's tasks
POST   /api/v1/developers           # Add developer
PUT    /api/v1/developers/:id       # Update developer
DELETE /api/v1/developers/:id       # Remove developer
```

### Tasks
```
GET    /api/v1/tasks                # List all tasks
GET    /api/v1/tasks/:id            # Task detail
POST   /api/v1/tasks                # Create task
PUT    /api/v1/tasks/:id            # Update task
DELETE /api/v1/tasks/:id            # Delete task
PATCH  /api/v1/tasks/:id/status     # Update status
```

### Activities
```
GET    /api/v1/activities           # Activity feed
GET    /api/v1/activities/recent    # Recent activities
POST   /api/v1/activities           # Log activity
```

### Dashboard
```
GET    /api/v1/dashboard/overview   # Team overview
GET    /api/v1/dashboard/stats      # Statistics
GET    /api/v1/dashboard/timeline   # Activity timeline
```

### WebSocket
```
WS     /api/v1/ws                   # Real-time updates
```

---

## 📁 Project Structure

```
task-manager/
├── backend/
│   ├── cmd/
│   │   └── server/
│   │       └── main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── handlers/
│   │   ├── services/
│   │   ├── models/
│   │   ├── repository/
│   │   └── middleware/
│   ├── pkg/
│   │   └── utils/
│   ├── migrations/
│   └── go.mod
├── frontend/
│   ├── app/
│   ├── components/
│   └── lib/
├── docker/
│   └── docker-compose.yml
├── docs/
│   ├── API.md
│   └── DATABASE.md
├── Makefile
└── README.md
```

---

## 📅 Timeline

### Day 1: Setup & Database
- **Flow (DevOps):** Docker + PostgreSQL + Redis
- **Schema (Database):** Database design + migrations

### Day 2: Backend Core
- **Neo (Backend):** HTTP server + basic structure
- **Cipher (Security):** JWT auth + middleware

### Day 3-4: Backend APIs
- **Neo (Backend):** All REST endpoints
- **Schema (Database):** Queries + optimization

### Day 5: Frontend
- **Voxel (Frontend):** Dashboard UI + components

### Day 6: Integration
- **All:** WebSocket + testing + deployment

---

## 🎨 UI Preview

### Dashboard
```
┌────────────────────────────────────────────────┐
│  Task Manager     Team: 6 | Active: 4 | Idle: 2│
├────────────────────────────────────────────────┤
│                                                │
│  Team Activity                                 │
│  ┌──────────────────────────────────────────┐ │
│  │ 🟢 Neo - Working on API endpoints        │ │
│  │ 🟢 Schema - Designing DB schema          │ │
│  │ 🟢 Flow - Docker setup complete          │ │
│  │ 🟡 Atlas - Coordinating tasks            │ │
│  │ ⚪ Voxel - Offline                       │ │
│  └──────────────────────────────────────────┘ │
│                                                │
│  Recent Tasks                                  │
│  ┌──────────────────────────────────────────┐ │
│  │ ✅ Setup Docker (Flow) - 2h ago          │ │
│  │ 🔄 Design schema (Schema) - in progress  │ │
│  │ 📋 Build API (Neo) - waiting             │ │
│  └──────────────────────────────────────────┘ │
│                                                │
└────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

```bash
# Clone/setup
cd /root/.openclaw-dev-team/task-manager

# Setup environment
cp .env.example .env

# Start services
make docker-up

# Run migrations
make migrate-up

# Run backend
make backend

# Run frontend
make frontend
```

---

## 📞 Team Coordination

**Channel:** #dev-general  
**Leader:** Atlas (Coordinator)

**All agents must:**
1. Read `/root/.openclaw-dev-team/PROJECT_STATE.md`
2. Check task for their channel
3. Update progress in PROJECT_STATE.md
4. Report blockers immediately

---

## 🎯 Success Criteria

- ✅ Real-time dashboard working
- ✅ All developers tracked
- ✅ Task CRUD working
- ✅ Activity feed live
- ✅ WebSocket updates
- ✅ Responsive UI
- ✅ Deployed successfully

---

**Let's build it! 🚀**
