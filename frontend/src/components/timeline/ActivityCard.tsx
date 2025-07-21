'use client'

import React, { useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { 
  ChevronDown, 
  ChevronUp, 
  MessageSquare, 
  Star, 
  Eye, 
  Target,
  Trophy,
  Users
} from 'lucide-react'
import { formatTimeAgo } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'
import { 
  
  getNoteSourceTypeLabel, 
  getFeedbackTypeColor 
} from '@/lib/constants/notes'
import { Person } from '@/lib/types'

export interface TimelineActivity {
  uuid: string
  type: 'feedback' | 'one_on_one' | 'observation' | 'mention' | 'achievement' | 'goal'
  content: string
  author_name: string
  author_id?: string
  created_at: string
  feedback_type?: 'positive' | 'constructive' | 'neutral'
  feedback_category?: string
  person_id?: string
  person_name?: string
  mentioned_people?: Person[]
  tags?: string[]
  needs_followup?: boolean
  is_resolved?: boolean
}

interface ActivityCardProps {
  activity: TimelineActivity
  targetPerson: Person
  allPeople: Person[]
  className?: string
}

const getActivityIcon = (type: string) => {
  switch (type) {
    case 'feedback':
      return <Star className="h-4 w-4" />
    case 'one_on_one':
      return <MessageSquare className="h-4 w-4" />
    case 'observation':
      return <Eye className="h-4 w-4" />
    case 'mention':
      return <Users className="h-4 w-4" />
    case 'achievement':
      return <Trophy className="h-4 w-4" />
    case 'goal':
      return <Target className="h-4 w-4" />
    default:
      return <MessageSquare className="h-4 w-4" />
  }
}

const getActivityColor = (type: string, feedbackType?: string) => {
  if (type === 'feedback' && feedbackType) {
    return getFeedbackTypeColor(feedbackType)
  }
  
  switch (type) {
    case 'feedback':
      return 'bg-yellow-500'
    case 'one_on_one':
      return 'bg-blue-500'
    case 'observation':
      return 'bg-green-500'
    case 'mention':
      return 'bg-purple-500'
    case 'achievement':
      return 'bg-orange-500'
    case 'goal':
      return 'bg-indigo-500'
    default:
      return 'bg-gray-500'
  }
}

const getDirectionLabel = (activity: TimelineActivity, targetPerson: Person) => {
  const isFromTarget = activity.author_name === targetPerson.name
  const isAboutTarget = activity.person_name === targetPerson.name || !activity.person_name
  
  if (activity.type === 'mention') {
    return `👤 ${activity.author_name} → Sobre: ${targetPerson.name}`
  }
  
  if (isFromTarget && !isAboutTarget) {
    return `👤 ${targetPerson.name} → Sobre: outros`
  }
  
  if (!isFromTarget && isAboutTarget) {
    return `👤 ${activity.author_name} → Para: ${targetPerson.name}`
  }
  
  return `👤 ${activity.author_name} ↔ ${targetPerson.name}`
}

export function ActivityCard({ 
  activity, 
  targetPerson, 
  
  className = '' 
}: ActivityCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  
  const activityColor = getActivityColor(activity.type, activity.feedback_type)
  const activityIcon = getActivityIcon(activity.type)
  const directionLabel = getDirectionLabel(activity, targetPerson)
  
  const typeLabel = activity.type === 'feedback' && activity.feedback_type 
    ? `Feedback ${activity.feedback_type === 'positive' ? 'Positivo' : activity.feedback_type === 'constructive' ? 'Construtivo' : 'Neutro'}`
    : getNoteSourceTypeLabel(activity.type)
  
  // Extract mentioned people from content
  const mentionRegex = /\{\{person:([^|]+)\|([^}]+)\}\}/g
  const mentions = []
  let match
  while ((match = mentionRegex.exec(activity.content)) !== null) {
    const personName = match[2]
    mentions.push(personName)
  }
  
  // Clean content for display
  const cleanContent = activity.content.replace(mentionRegex, '@$2')
  
  const isLongContent = cleanContent.length > 120
  const shouldTruncate = !isExpanded && isLongContent
  const displayContent = shouldTruncate 
    ? cleanContent.substring(0, 120) + '...' 
    : cleanContent

  return (
    <Card className={`timeline-card ${className}`}>
      <CardContent className="p-4">
        {/* Header */}
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-3 flex-1 min-w-0">
            <div className={`p-2 rounded-full text-white ${activityColor}`}>
              {activityIcon}
            </div>
            
            <div className="flex-1 min-w-0">
              <div className="flex items-center gap-2 mb-1">
                <span className="font-medium text-sm">{typeLabel}</span>
                <span className="text-xs text-muted-foreground">•</span>
                <span className="text-xs text-muted-foreground">
                  {formatTimeAgo(new Date(activity.created_at))}
                </span>
              </div>
              
              <div className="text-xs text-muted-foreground">
                {directionLabel}
              </div>
            </div>
          </div>
          
          {/* Tags */}
          <div className="flex gap-1 ml-2">
            {activity.feedback_category && (
              <Badge variant="outline" className="text-xs">
                {activity.feedback_category}
              </Badge>
            )}
            {activity.tags?.map((tag, index) => (
              <Badge key={index} variant="outline" className="text-xs">
                {tag}
              </Badge>
            ))}
            {activity.needs_followup && (
              <Badge variant="destructive" className="text-xs">
                Follow-up
              </Badge>
            )}
          </div>
        </div>

        {/* Content */}
        <div className="mb-3">
          <p className="text-sm leading-relaxed">
            {displayContent}
          </p>
          
          {isLongContent && (
            <Button
              variant="ghost"
              size="sm"
              className="h-6 px-2 mt-1 text-xs"
              onClick={() => setIsExpanded(!isExpanded)}
            >
              {isExpanded ? (
                <>
                  <ChevronUp className="h-3 w-3 mr-1" />
                  Ver menos
                </>
              ) : (
                <>
                  <ChevronDown className="h-3 w-3 mr-1" />
                  Ver mais
                </>
              )}
            </Button>
          )}
        </div>

        {/* Mentions */}
        {mentions.length > 0 && (
          <div className="flex items-center gap-2 mb-3">
            <span className="text-xs text-muted-foreground">Mencionou:</span>
            {mentions.map((mention, index) => (
              <Badge key={index} variant="secondary" className="text-xs">
                @{mention}
              </Badge>
            ))}
          </div>
        )}

        {/* Footer with actions */}
        <div className="flex items-center justify-between pt-2 border-t border-border/50">
          <div className="flex items-center gap-2">
            {activity.author_id && (
              <div className="flex items-center gap-2">
                <Avatar className="h-6 w-6">
                  <AvatarFallback className="text-xs">
                    {getInitials(activity.author_name)}
                  </AvatarFallback>
                </Avatar>
                <span className="text-xs text-muted-foreground">
                  {activity.author_name}
                </span>
              </div>
            )}
          </div>
          
          <div className="flex items-center gap-1">
            {activity.needs_followup && !activity.is_resolved && (
              <Button variant="outline" size="sm" className="h-6 px-2 text-xs">
                Resolver
              </Button>
            )}
          </div>
        </div>
      </CardContent>
    </Card>
  )
}