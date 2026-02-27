import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import type { User } from './types'
import { api } from './api'

interface AuthState {
  user: User | null
  token: string | null
  isLoading: boolean
  error: string | null
  
  // Actions
  login: (email: string, password: string) => Promise<void>
  register: (email: string, password: string, name: string) => Promise<void>
  logout: () => void
  fetchUser: () => Promise<void>
  clearError: () => void
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set, get) => ({
      user: null,
      token: null,
      isLoading: false,
      error: null,

      login: async (email, password) => {
        set({ isLoading: true, error: null })
        try {
          const response = await api.post('/auth/login', { email, password })
          const { user, token } = response.data.data
          localStorage.setItem('token', token)
          set({ user, token, isLoading: false })
        } catch (error: unknown) {
          const message = (error as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Login failed'
          set({ error: message, isLoading: false })
          throw error
        }
      },

      register: async (email, password, name) => {
        set({ isLoading: true, error: null })
        try {
          const response = await api.post('/auth/register', { email, password, name })
          const { user, token } = response.data.data
          localStorage.setItem('token', token)
          set({ user, token, isLoading: false })
        } catch (error: unknown) {
          const message = (error as { response?: { data?: { error?: string } } })?.response?.data?.error || 'Registration failed'
          set({ error: message, isLoading: false })
          throw error
        }
      },

      logout: () => {
        localStorage.removeItem('token')
        set({ user: null, token: null })
      },

      fetchUser: async () => {
        const token = get().token || localStorage.getItem('token')
        if (!token) return
        
        set({ isLoading: true })
        try {
          const response = await api.get('/auth/me')
          set({ user: response.data.data, isLoading: false })
        } catch {
          localStorage.removeItem('token')
          set({ user: null, token: null, isLoading: false })
        }
      },

      clearError: () => set({ error: null }),
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({ token: state.token }),
    }
  )
)
