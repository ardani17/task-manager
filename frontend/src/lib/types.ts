// API Types
export interface User {
  id: string
  email: string
  name: string
  role: 'admin' | 'member'
  status: 'active' | 'inactive'
  created_at: string
  updated_at: string
}

export interface Project {
  id: string
  name: string
  description: string
  owner_id: string
  status: 'active' | 'completed' | 'archived'
  due_date: string | null
  created_at: string
  updated_at: string
}

export interface Task {
  id: string
  title: string
  description: string
  project_id: string
  assigned_to: string | null
  priority: 'low' | 'medium' | 'high'
  status: 'todo' | 'in_progress' | 'review' | 'done'
  due_date: string | null
  estimated_hours: number | null
  actual_hours: number | null
  created_at: string
  updated_at: string
}

export interface Activity {
  id: string
  user_id: string
  action: string
  entity_type: string
  entity_id: string
  metadata: Record<string, unknown>
  created_at: string
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: string
  message?: string
}

export interface LoginResponse {
  user: User
  token: string
}

export interface PaginatedResponse<T> {
  items: T[]
  pagination?: {
    page: number
    limit: number
    total: number
  }
}
