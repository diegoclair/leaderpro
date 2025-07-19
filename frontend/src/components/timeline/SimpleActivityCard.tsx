'use client'

import React, { useState } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { 
  ChevronDown, 
  ChevronUp, 
  MessageSquare, 
  Star, 
  Eye
} from 'lucide-react'
import { Person } from '@/lib/types'
import { formatDateRelative, formatDateExact } from '@/lib/utils/dates'

export interface TimelineActivity {
  uuid: string
  type: 'feedback' | 'one_on_one' | 'observation' | 'mention'
  content: string
  author_name: string
  created_at: string
  feedback_type?: 'positive' | 'constructive' | 'neutral'
  feedback_category?: string
  person_name?: string
  entry_source?: string
}

interface SimpleActivityCardProps {
  activity: TimelineActivity
  className?: string
}

export function SimpleActivityCard({ activity, className = '' }: SimpleActivityCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  
  // Determine if content is long and needs truncation
  const contentLines = activity.content.split('\n')
  const isLongContent = activity.content.length > 200 || contentLines.length > 3
  const shouldTruncate = isLongContent && !isExpanded
  
  // Display content (truncated or full)
  const displayContent = shouldTruncate 
    ? activity.content.slice(0, 200) + (activity.content.length > 200 ? '...' : '')
    : activity.content

  // Process @mentions in content
  const processedContent = displayContent.replace(
    /\{\{person:([^|]+)\|([^}]+)\}\}/g,
    '<span class="inline-flex items-center px-2 py-1 text-xs font-medium bg-blue-50 text-blue-700 dark:bg-blue-900/20 dark:text-blue-300 rounded-full">@$2</span>'
  )

  // Get icon and color based on type
  const getTypeConfig = () => {
    switch (activity.type) {
      case 'feedback':
        return { 
          icon: <Star className="w-4 h-4" />, 
          color: 'bg-yellow-500',
          label: 'Feedback',
          borderColor: 'border-l-yellow-400'
        }
      case 'one_on_one':
        return { 
          icon: <MessageSquare className="w-4 h-4" />, 
          color: 'bg-blue-500',
          label: '1:1',
          borderColor: 'border-l-blue-400'
        }
      case 'observation':
        return { 
          icon: <Eye className="w-4 h-4" />, 
          color: 'bg-green-500',
          label: 'Observação',
          borderColor: 'border-l-green-400'
        }
      default:
        return { 
          icon: <MessageSquare className="w-4 h-4" />, 
          color: 'bg-gray-500',
          label: 'Nota',
          borderColor: 'border-l-gray-400'
        }
    }
  }

  const typeConfig = getTypeConfig()

  return (
    <div className={`border-l-4 ${typeConfig.borderColor} bg-card p-4 rounded-r-lg shadow-sm hover:shadow-md transition-all duration-200 ${className}`}>
      {/* Header */}
      <div className="flex items-start justify-between mb-3">
        <div className="flex items-center gap-3">
          <div className={`w-8 h-8 rounded-full ${typeConfig.color} flex items-center justify-center text-white`}>
            {typeConfig.icon}
          </div>
          
          <div>
            <div className="flex items-center gap-2 mb-1">
              <span className="font-semibold text-foreground">{typeConfig.label}</span>
              {activity.feedback_type && (
                <Badge variant={activity.feedback_type === 'positive' ? 'default' : 'secondary'} className="text-xs">
                  {activity.feedback_type === 'positive' ? '✨ Positivo' : 
                   activity.feedback_type === 'constructive' ? '🔨 Construtivo' : '⚪ Neutro'}
                </Badge>
              )}
            </div>
            <div 
              className="text-sm text-muted-foreground cursor-help" 
              title={formatDateExact(activity.created_at)}
            >
              {formatDateRelative(activity.created_at)}
            </div>
          </div>
        </div>
        
        {/* Category tag if exists */}
        {activity.feedback_category && (
          <Badge variant="outline" className="text-xs capitalize">
            {activity.feedback_category}
          </Badge>
        )}
      </div>

      {/* Content */}
      <div className="mb-3">
        <div 
          className="text-foreground leading-relaxed whitespace-pre-wrap"
          dangerouslySetInnerHTML={{ __html: processedContent }}
        />
        
        {isLongContent && (
          <Button
            variant="ghost"
            size="sm"
            className="mt-2 h-auto p-0 text-blue-600 hover:text-blue-700 dark:text-blue-400 dark:hover:text-blue-300"
            onClick={() => setIsExpanded(!isExpanded)}
          >
            {isExpanded ? (
              <>
                <ChevronUp className="h-4 w-4 mr-1" />
                Ver menos
              </>
            ) : (
              <>
                <ChevronDown className="h-4 w-4 mr-1" />
                Ver mais
              </>
            )}
          </Button>
        )}
      </div>
    </div>
  )
}