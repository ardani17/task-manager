'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/lib/auth-store'
import { api } from '@/lib/api'
import { Header } from '@/components/layout/Header'
import { Sidebar } from '@/components/layout/Sidebar'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Search, Mail, Calendar } from 'lucide-react'
import type { User, Task } from '@/lib/types'

export default function TeamPage() {
  const { user, fetchUser } = useAuthStore()
  const [users, setUsers] = useState<User[]>([])
  const [tasks, setTasks] = useState<Task[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  useEffect(() => {
    async function fetchData() {
      try {
        const [usersRes, tasksRes] = await Promise.all([
          api.get('/users'),
          api.get('/tasks'),
        ])
        setUsers(usersRes.data.data?.users || [])
        setTasks(tasksRes.data.data?.tasks || [])
      } catch (error) {
        console.error('Failed to fetch data:', error)
      } finally {
        setIsLoading(false)
      }
    }

    if (user) {
      fetchData()
    }
  }, [user])

  const filteredUsers = users.filter((u) =>
    u.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
    u.email.toLowerCase().includes(searchQuery.toLowerCase())
  )

  const getUserStats = (userId: string) => {
    const userTasks = tasks.filter((t) => t.assigned_to === userId)
    return {
      total: userTasks.length,
      completed: userTasks.filter((t) => t.status === 'done').length,
      inProgress: userTasks.filter((t) => t.status === 'in_progress').length,
    }
  }

  if (!user) return null

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Header title="Team" />
        <main className="flex-1 p-6 overflow-auto">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Header */}
              <div className="flex items-center justify-between">
                <div>
                  <h2 className="text-2xl font-bold">Team Members</h2>
                  <p className="text-muted-foreground">{users.length} members in your team</p>
                </div>
              </div>

              {/* Search */}
              <div className="relative w-64">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                <Input
                  placeholder="Search members..."
                  className="pl-9"
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                />
              </div>

              {/* Team Grid */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
                {filteredUsers.map((member) => {
                  const stats = getUserStats(member.id)
                  return (
                    <Card key={member.id} className="hover:shadow-md transition-shadow">
                      <CardHeader className="pb-2">
                        <div className="flex items-center gap-4">
                          <Avatar className="w-12 h-12">
                            <AvatarFallback className="text-lg">
                              {member.name.charAt(0).toUpperCase()}
                            </AvatarFallback>
                          </Avatar>
                          <div>
                            <CardTitle className="text-base">{member.name}</CardTitle>
                            <div className="flex items-center gap-1 text-sm text-muted-foreground">
                              <Mail className="w-3 h-3" />
                              <span className="truncate max-w-[180px]">{member.email}</span>
                            </div>
                          </div>
                        </div>
                      </CardHeader>
                      <CardContent className="space-y-3">
                        <div className="flex items-center gap-2">
                          <Badge variant={member.status === 'active' ? 'default' : 'secondary'}>
                            {member.status}
                          </Badge>
                          <Badge variant="outline">{member.role}</Badge>
                        </div>
                        <div className="grid grid-cols-3 gap-2 pt-2 border-t">
                          <div className="text-center">
                            <p className="text-2xl font-bold">{stats.total}</p>
                            <p className="text-xs text-muted-foreground">Tasks</p>
                          </div>
                          <div className="text-center">
                            <p className="text-2xl font-bold text-green-500">{stats.completed}</p>
                            <p className="text-xs text-muted-foreground">Done</p>
                          </div>
                          <div className="text-center">
                            <p className="text-2xl font-bold text-blue-500">{stats.inProgress}</p>
                            <p className="text-xs text-muted-foreground">Active</p>
                          </div>
                        </div>
                        <div className="flex items-center gap-1 text-xs text-muted-foreground pt-2 border-t">
                          <Calendar className="w-3 h-3" />
                          <span>Joined {new Date(member.created_at).toLocaleDateString()}</span>
                        </div>
                      </CardContent>
                    </Card>
                  )
                })}
              </div>

              {filteredUsers.length === 0 && (
                <div className="text-center py-12">
                  <p className="text-muted-foreground">No team members found</p>
                </div>
              )}
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
