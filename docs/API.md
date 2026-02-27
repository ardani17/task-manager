# API Documentation

Complete REST API documentation for Task Manager backend.

---

## 📍 Base URL

**Development:** `http://localhost:8081/api/v1`  
**Production:** `https://your-domain.com/api/v1`

---

## 🔐 Authentication

All protected endpoints require JWT token in Authorization header:

```http
Authorization: Bearer <your-jwt-token>
```

### Get Token
```bash
# Register
curl -X POST http://localhost:8081/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"password123"}'

# Login
curl -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "uuid",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "developer"
  }
}
```

---

## 📋 Endpoints

### Health Check

#### GET /health
Check server health status.

**Auth Required:** No

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2026-02-27T14:00:00Z",
  "version": "1.0.0",
  "uptime": "1h30m"
}
```

---

### Authentication

#### POST /auth/register
Register a new user.

**Auth Required:** No

**Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "securepassword123",
  "role": "developer"
}
```

**Response (201):**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "developer",
    "status": "offline",
    "created_at": "2026-02-27T14:00:00Z"
  }
}
```

#### POST /auth/login
Login and get JWT token.

**Auth Required:** No

**Body:**
```json
{
  "email": "john@example.com",
  "password": "securepassword123"
}
```

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "John Doe",
    "email": "john@example.com",
    "role": "developer"
  }
}
```

#### GET /auth/me
Get current authenticated user.

**Auth Required:** Yes

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "developer",
  "status": "active",
  "last_active": "2026-02-27T14:00:00Z"
}
```

#### POST /auth/logout
Logout current user.

**Auth Required:** Yes

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

#### POST /auth/refresh
Refresh JWT token.

**Auth Required:** Yes

**Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_at": "2026-02-28T14:00:00Z"
}
```

---

### Tasks

#### GET /tasks
Get all tasks with optional filters.

**Auth Required:** Yes

**Query Parameters:**
- `status` (optional): Filter by status (todo, in_progress, review, done)
- `priority` (optional): Filter by priority (low, medium, high)
- `assignee_id` (optional): Filter by assignee
- `project_id` (optional): Filter by project

**Response (200):**
```json
{
  "tasks": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "title": "Implement user authentication",
      "description": "Add JWT auth to the API",
      "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
      "project_id": "550e8400-e29b-41d4-a716-446655440002",
      "status": "in_progress",
      "priority": "high",
      "estimated_hours": 8.5,
      "actual_hours": 6.0,
      "due_date": "2026-03-01",
      "created_at": "2026-02-27T14:00:00Z",
      "updated_at": "2026-02-27T14:30:00Z"
    }
  ],
  "total": 1
}
```

#### POST /tasks
Create a new task.

**Auth Required:** Yes

**Body:**
```json
{
  "title": "Implement user authentication",
  "description": "Add JWT auth to the API",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
  "project_id": "550e8400-e29b-41d4-a716-446655440002",
  "status": "todo",
  "priority": "high",
  "estimated_hours": 8.5,
  "due_date": "2026-03-01"
}
```

**Response (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Implement user authentication",
  "status": "todo",
  "created_at": "2026-02-27T14:00:00Z"
}
```

#### GET /tasks/:id
Get a specific task by ID.

**Auth Required:** Yes

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Implement user authentication",
  "description": "Add JWT auth to the API",
  "assignee_id": "550e8400-e29b-41d4-a716-446655440001",
  "project_id": "550e8400-e29b-41d4-a716-446655440002",
  "status": "in_progress",
  "priority": "high",
  "estimated_hours": 8.5,
  "actual_hours": 6.0,
  "due_date": "2026-03-01",
  "created_at": "2026-02-27T14:00:00Z",
  "updated_at": "2026-02-27T14:30:00Z"
}
```

#### PUT /tasks/:id
Update a task completely.

**Auth Required:** Yes

**Body:**
```json
{
  "title": "Updated task title",
  "description": "Updated description",
  "status": "in_progress",
  "priority": "medium",
  "actual_hours": 7.5
}
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "title": "Updated task title",
  "updated_at": "2026-02-27T15:00:00Z"
}
```

#### DELETE /tasks/:id
Delete a task.

**Auth Required:** Yes

**Response (200):**
```json
{
  "message": "Task deleted successfully"
}
```

#### PATCH /tasks/:id/status
Update task status only.

**Auth Required:** Yes

**Body:**
```json
{
  "status": "done"
}
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "done",
  "updated_at": "2026-02-27T15:00:00Z"
}
```

---

### Projects

#### GET /projects
Get all projects.

**Auth Required:** Yes

**Response (200):**
```json
{
  "projects": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Task Manager",
      "description": "Developer progress monitoring system",
      "status": "active",
      "created_at": "2026-02-27T14:00:00Z",
      "updated_at": "2026-02-27T14:00:00Z"
    }
  ],
  "total": 1
}
```

#### POST /projects
Create a new project.

**Auth Required:** Yes

**Body:**
```json
{
  "name": "Task Manager",
  "description": "Developer progress monitoring system"
}
```

**Response (201):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Task Manager",
  "status": "active",
  "created_at": "2026-02-27T14:00:00Z"
}
```

#### GET /projects/:id
Get a specific project.

**Auth Required:** Yes

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Task Manager",
  "description": "Developer progress monitoring system",
  "status": "active",
  "created_at": "2026-02-27T14:00:00Z"
}
```

#### PUT /projects/:id
Update a project.

**Auth Required:** Yes

**Body:**
```json
{
  "name": "Updated Project Name",
  "description": "Updated description",
  "status": "completed"
}
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Updated Project Name",
  "updated_at": "2026-02-27T15:00:00Z"
}
```

#### DELETE /projects/:id
Delete a project.

**Auth Required:** Yes

**Response (200):**
```json
{
  "message": "Project deleted successfully"
}
```

---

### Users

#### GET /users
Get all users.

**Auth Required:** Yes

**Response (200):**
```json
{
  "users": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "John Doe",
      "email": "john@example.com",
      "role": "developer",
      "status": "active",
      "last_active": "2026-02-27T14:00:00Z"
    }
  ],
  "total": 1
}
```

#### GET /users/:id
Get a specific user.

**Auth Required:** Yes

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Doe",
  "email": "john@example.com",
  "role": "developer",
  "status": "active",
  "last_active": "2026-02-27T14:00:00Z",
  "created_at": "2026-02-27T14:00:00Z"
}
```

#### PUT /users/:id
Update user profile.

**Auth Required:** Yes

**Body:**
```json
{
  "name": "John Smith",
  "email": "johnsmith@example.com",
  "role": "senior_developer"
}
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "John Smith",
  "updated_at": "2026-02-27T15:00:00Z"
}
```

#### DELETE /users/:id
Delete a user.

**Auth Required:** Yes

**Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

#### PATCH /users/:id/status
Update user status.

**Auth Required:** Yes

**Body:**
```json
{
  "status": "active"
}
```

**Response (200):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "active",
  "updated_at": "2026-02-27T15:00:00Z"
}
```

---

### Activity

#### GET /activity
Get activity feed with filters.

**Auth Required:** Yes

**Query Parameters:**
- `developer_id` (optional): Filter by developer
- `action` (optional): Filter by action type
- `entity_type` (optional): Filter by entity (task, project, user)
- `limit` (optional): Number of results (default: 50)
- `offset` (optional): Pagination offset

**Response (200):**
```json
{
  "activities": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "developer_id": "550e8400-e29b-41d4-a716-446655440001",
      "action": "task_completed",
      "entity_type": "task",
      "entity_id": "550e8400-e29b-41d4-a716-446655440002",
      "metadata": {
        "task_title": "Implement authentication",
        "hours_spent": 6.5
      },
      "created_at": "2026-02-27T14:00:00Z"
    }
  ],
  "total": 1
}
```

---

## ❌ Error Responses

All endpoints follow standard error format:

### 400 Bad Request
```json
{
  "error": "Invalid request body",
  "details": "Field 'email' is required"
}
```

### 401 Unauthorized
```json
{
  "error": "Unauthorized",
  "message": "Invalid or expired token"
}
```

### 404 Not Found
```json
{
  "error": "Not found",
  "message": "Task not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "Internal server error",
  "message": "Database connection failed"
}
```

---

## 📝 Data Models

### Task Status Values
- `todo` - Task created, not started
- `in_progress` - Currently being worked on
- `review` - Completed, awaiting review
- `done` - Completed and approved

### Task Priority Values
- `low` - Low priority
- `medium` - Medium priority (default)
- `high` - High priority

### User Status Values
- `offline` - User not active
- `active` - User is online and active
- `idle` - User is away

### Activity Actions
- `task_created`
- `task_updated`
- `task_completed`
- `task_deleted`
- `project_created`
- `project_updated`
- `user_registered`
- `user_login`

---

## 🔧 Rate Limiting

Currently no rate limiting implemented. Recommended for production:
- 100 requests per minute per user
- 1000 requests per minute per IP

---

## 📚 Examples

### Create and Complete a Task

```bash
# 1. Login
TOKEN=$(curl -s -X POST http://localhost:8081/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@example.com","password":"password"}' | jq -r '.token')

# 2. Create task
TASK_ID=$(curl -s -X POST http://localhost:8081/api/v1/tasks \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"title":"New feature","description":"Add new feature"}' | jq -r '.id')

# 3. Update status to in_progress
curl -X PATCH http://localhost:8081/api/v1/tasks/$TASK_ID/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"in_progress"}'

# 4. Complete task
curl -X PATCH http://localhost:8081/api/v1/tasks/$TASK_ID/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"done"}'
```

---

## 🧪 Testing

### Using cURL
All examples above use cURL. Copy and modify as needed.

### Using Postman
1. Import collection: `docs/postman_collection.json`
2. Set environment variables:
   - `base_url`: http://localhost:8081/api/v1
   - `token`: (auto-set after login)

---

## 📖 OpenAPI Specification

Full OpenAPI 3.0 spec available at: `docs/openapi.yaml`

---

## 🔗 Related Documentation

- [INSTALLATION.md](INSTALLATION.md) - Installation guide
- [README.md](README.md) - Project overview
- Database schema: `backend/migrations/`

---

**Last updated:** 2026-02-27
