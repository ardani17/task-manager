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
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Search, Clock, Filter } from 'lucide-react'
import { format, formatDistanceToNow } from 'date-fns'
import type { Activity, User } from '@/lib/types'

const actionColors: Record<string, 'default' | 'secondary' | 'outline' | 'destructive'> = {
  created: 'default',
  updated: 'secondary',
  deleted: 'destructive',
  completed: 'default',
  assigned: 'outline',
}

export default function ActivityPage() {
  const { user, fetchUser } = useAuthStore()
  const [activities, setActivities] = useState<Activity[]>([])
  const [users, setUsers] = useState<User[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [searchQuery, setSearchQuery] = useState('')
  const [entityFilter, setEntityFilter] = useState<string>('all')

  useEffect(() => {
    fetchUser()
  }, [fetchUser])

  useEffect(() => {
    async function fetchData() {
      try {
        const [activityRes, usersRes] = await Promise.all([
          api.get('/activity?limit=50'),
          api.get('/users'),
        ])
        setActivities(activityRes.data.data?.activities || [])
        setUsers(usersRes.data.data?.users || [])
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

  const filteredActivities = activities.filter((activity) => {
    const matchesSearch = activity.action.toLowerCase().includes(searchQuery.toLowerCase())
    const matchesEntity = entityFilter === 'all' || activity.entity_type === entityFilter
    return matchesSearch && matchesEntity
  })

  const getUser = (userId: string) => {
    return users.find((u) => u.id === userId)
  }

  if (!user) return null

  return (
    <div className="flex h-screen">
      <Sidebar />
      <div className="flex-1 flex flex-col">
        <Header title="Activity" />
        <main className="flex-1 p-6 overflow-auto">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
            </div>
          ) : (
            <div className="space-y-6">
              {/* Header */}
              <div>
                <h2 className="text-2xl font-bold">Activity Feed</h2>
                <p className="text-muted-foreground">Track all team activities in real-time</p>
              </div>

              {/* Filters */}
              <div className="flex gap-4">
                <div className="relative flex-1 max-w-sm">
                  <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-muted-foreground" />
                  <Input
                    placeholder="Search activities..."
                    className="pl-9"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                  />
                </div>
                <Select value={entityFilter} onValueChange={setEntityFilter}>
                  <SelectTrigger className="w-40">
                    <Filter className="w-4 h-4 mr-2" />
                    <SelectValue placeholder="Type" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All Types</SelectItem>
                    <SelectItem value="task">Tasks</SelectItem>
                    <SelectItem value="project">Projects</SelectItem>
                    <SelectItem value="user">Users</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {/* Activity Timeline */}
              <Card>
                <CardHeader>
                  <CardTitle className="text-lg">Recent Activity</CardTitle>
                </CardHeader>
                <CardContent>
                  <ScrollArea className="h-[600px] pr-4">
                    <div className="space-y-6">
                      {filteredActivities.map((activity, index) => {
                        const actor = getUser(activity.user_id)
                        const actionType = activity.action.split(' ')[0].toLowerCase()
                        return (
                          <div key={activity.id} className="relative flex gap-4">
                            {/* Timeline line */}
                            {index !== filteredActivities.length - 1 && (
                              <div className="absolute left-5 top-12 w-px h-full bg-border" />
                            )}
                            
                            {/* Avatar */}
                            <Avatar className="w-10 h-10 z-10">
                              <AvatarFallback>
                                {actor?.name.charAt(0).toUpperCase() || 'U'}
                              </AvatarFallback>
                            </Avatar>

                            {/* Content */}
                            <div className="flex-1 space-y-1 pb-6">
                              <div className="flex items-center gap-2 flex-wrap">
                                <span className="font-medium">{actor?.name || 'Unknown User'}</span>
                                <span className="text-muted-foreground">{activity.action}</span>
                                <Badge variant={actionColors[actionType] || 'secondary'} className="text-xs">
                                  {activity.entity_type}
                                </Badge>
                              </div>
                              
                              {activity.metadata && Object.keys(activity.metadata).length > 0 && (
                                <div className="text-sm text-muted-foreground bg-muted/50 p-2 rounded-lg mt-2">
                                  {Object.entries(activity.metadata).map(([key, value]) => (
                                    <span key={key} className="mr-4">
                                      <span className="font-medium">{key}:</span> {String(value)}
                                    </span>
                                  ))}
                                </div>
                              )}

                              <div className="flex items-center gap-1 text-xs text-muted-foreground pt-1">
                                <Clock className="w-3 h-3" />
                                <span title={format(new Date(activity.created_at), 'PPpp')}>
                                  {formatDistanceToNow(new Date(activity.created_at), { addSuffix: true })}
                                </span>
                              </div>
                            </div>
                          </div>
                        )
                      })}

                      {filteredActivities.length === 0 && (
                        <div className="text-center py-12">
                          <p className="text-muted-foreground">No activity found</p>
                        </div>
                      )}
                    </div>
                  </ScrollArea>
                </CardContent>
              </Card>
            </div>
          )}
        </main>
      </div>
    </div>
  )
}
