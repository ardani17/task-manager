'use client'

import { useEffect, useState } from 'react'
import { useAuthStore } from '@/lib/auth-store'
import { api } from '@/lib/api'
import { Header } from '@/components/layout/Header'
import { Sidebar } from '@/components/layout/Sidebar'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { CheckSquare, FolderKanban, Users, Clock, TrendingUp } from 'lucide-react'
import type { Task, Project, User, Activity } from '@/lib/types'

interface Stats {
  totalTasks: number
  completedTasks: number
  totalProjects: number
  activeUsers: number
}

export default function DashboardPage() {
  const { user, fetchUser } = useAuthStore()
  const [stats, setStats] = useState<Stats>({
    totalTasks: 0,
    completedTasks: 0,
    totalProjects: 0,
    activeUsers: 0,
  })
  const [recentTasks, setRecentTasks] = useState<Task[]>([])
  const [activities, setActivities] = useState<Activity[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  useEffect(() => {
    async function fetchData() {
      try {
        const [tasksRes, projectsRes, usersRes, activityRes] = await Promise.all([
          api.get('/tasks'),
          api.get('/projects'),
          api.get('/users'),
          api.get('/activity?limit=5'),
        ])

        const tasks = tasksRes.data.data?.tasks || []
        const projects = projectsRes.data.data?.projects || []
        const users = usersRes.data.data?.users || []
        const activityData = activityRes.data.data?.activities || []

        setStats({
          totalTasks: tasks.length,
          completedTasks: tasks.filter((t: Task) => t.status === 'done').length,
          totalProjects: projects.length,
          activeUsers: users.filter((u: User) => u.status === 'active').length,
        })

        setRecentTasks(tasks.slice(0, 5))
        setActivities(activityData)
      } catch (error) {
        console.error('Failed to fetch dashboard data:', error)
      } finally {
        setIsLoading(false)
      }
    }

    if (user) {
      fetchData()
    }
  }, [user])

  if (!user) {
    return null
  }

  const statCards = [
    {
      title: 'Total Tasks',
      value: stats.totalTasks,
      icon: CheckSquare,
      color: 'text-blue-500',
    },
    {
      title: 'Completed',
      value: stats.completedTasks,
      icon: TrendingUp,
      color: 'text-green-500',
    },
    {
      title: 'Projects',
      value: stats.totalProjects,
      icon: FolderKanban,
      color: 'text-purple-500',
    },
    {
      title: 'Team Members',
      value: stats.activeUsers,
      icon: Users,
      color: 'text-orange-500',
    },
  ]

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Header title="Dashboard" />
        <main className="flex-1 p-6 overflow-auto">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Welcome */}
              <div>
                <h2 className="text-2xl font-bold">Welcome back, {user.name}!</h2>
                <p className="text-muted-foreground">Here&apos;s what&apos;s happening with your projects.</p>
              </div>

              {/* Stats Grid */}
              <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
                {statCards.map((stat) => {
                  const Icon = stat.icon
                  return (
                    <Card key={stat.title}>
                      <CardHeader className="flex flex-row items-center justify-between pb-2">
                        <CardTitle className="text-sm font-medium text-muted-foreground">
                          {stat.title}
                        </CardTitle>
                        <Icon className={`w-5 h-5 ${stat.color}`} />
                      </CardHeader>
                      <CardContent>
                        <div className="text-2xl font-bold">{stat.value}</div>
                      </CardContent>
                    </Card>
                  )
                })}
              </div>

              {/* Content Grid */}
              <div className="grid gap-6 md:grid-cols-2">
                {/* Recent Tasks */}
                <Card>
                  <CardHeader>
                    <CardTitle>Recent Tasks</CardTitle>
                    <CardDescription>Your latest assigned tasks</CardDescription>
                  </CardHeader>
                  <CardContent>
                    {recentTasks.length === 0 ? (
                      <p className="text-muted-foreground text-sm">No tasks yet</p>
                    ) : (
                      <div className="space-y-4">
                        {recentTasks.map((task) => (
                          <div key={task.id} className="flex items-center justify-between">
                            <div className="space-y-1">
                              <p className="text-sm font-medium">{task.title}</p>
                              <p className="text-xs text-muted-foreground">
                                {task.due_date ? `Due: ${new Date(task.due_date).toLocaleDateString()}` : 'No due date'}
                              </p>
                            </div>
                            <Badge variant={
                              task.status === 'done' ? 'default' :
                              task.status === 'in_progress' ? 'secondary' :
                              task.status === 'review' ? 'outline' : 'destructive'
                            }>
                              {task.status}
                            </Badge>
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>

                {/* Activity Feed */}
                <Card>
                  <CardHeader>
                    <CardTitle>Recent Activity</CardTitle>
                    <CardDescription>Latest updates from your team</CardDescription>
                  </CardHeader>
                  <CardContent>
                    {activities.length === 0 ? (
                      <p className="text-muted-foreground text-sm">No recent activity</p>
                    ) : (
                      <div className="space-y-4">
                        {activities.map((activity) => (
                          <div key={activity.id} className="flex items-start gap-3">
                            <div className="w-8 h-8 rounded-full bg-primary/10 flex items-center justify-center">
                              <Clock className="w-4 h-4 text-primary" />
                            </div>
                            <div className="space-y-1">
                              <p className="text-sm">{activity.action}</p>
                              <p className="text-xs text-muted-foreground">
                                {new Date(activity.created_at).toLocaleString()}
                              </p>
                            </div>
                          </div>
                        ))}
                      </div>
                    )}
                  </CardContent>
                </Card>
              </div>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
