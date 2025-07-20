'use client'

import React, { useState, useRef, useEffect } from 'react'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { 
  ChevronDown, 
  ChevronUp, 
  MessageSquare, 
  Star, 
  Eye,
  MoreVertical,
  Edit3,
  Trash2
} from 'lucide-react'
import { Person } from '@/lib/types'
import { formatDateRelative, formatDateExact } from '@/lib/utils/dates'

// Função para traduzir categorias de feedback
const translateFeedbackCategory = (category: string): string => {
  const translations: Record<string, string> = {
    'performance': 'Performance',
    'behavior': 'Comportamento', 
    'skill': 'Habilidade',
    'collaboration': 'Colaboração'
  }
  return translations[category] || category
}

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
  onEdit?: (activity: TimelineActivity) => void
  onDelete?: (activity: TimelineActivity) => void
}

export function SimpleActivityCard({ activity, className = '', onEdit, onDelete }: SimpleActivityCardProps) {
  const [isExpanded, setIsExpanded] = useState(false)
  const [shouldShowExpand, setShouldShowExpand] = useState(false)
  const contentRef = useRef<HTMLDivElement>(null)
  const [isClient, setIsClient] = useState(false)
  
  // Ensure we're on the client side
  useEffect(() => {
    setIsClient(true)
  }, [])
  
  // Responsive truncation based on actual rendered height
  useEffect(() => {
    if (!isClient || !contentRef.current || isExpanded) return
    
    const checkHeight = () => {
      if (!contentRef.current) return
      
      // Get the actual height of the rendered content
      const contentHeight = contentRef.current.scrollHeight
      const containerHeight = contentRef.current.clientHeight
      
      // Define max height based on screen size
      const screenWidth = window.innerWidth
      let maxHeight: number
      
      if (screenWidth < 640) { // mobile
        maxHeight = 120 // ~4 lines on mobile
      } else if (screenWidth < 1024) { // tablet  
        maxHeight = 140 // ~4-5 lines on tablet
      } else { // desktop
        maxHeight = 160 // ~5-6 lines on desktop
      }
      
      setShouldShowExpand(contentHeight > maxHeight)
    }
    
    // Check height immediately and on window resize
    checkHeight()
    window.addEventListener('resize', checkHeight)
    
    return () => window.removeEventListener('resize', checkHeight)
  }, [isClient, activity.content, isExpanded])
  
  // Simple content processing without pre-truncation
  const displayContent = activity.content

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
    <div className={`border-l-4 ${typeConfig.borderColor} bg-card p-3 rounded-r-lg shadow-sm hover:shadow-md transition-all duration-200 ${className}`}>
      {/* Compact Header */}
      <div className="flex items-center justify-between mb-2">
        <div className="flex items-center gap-2">
          <div className={`w-6 h-6 rounded-full ${typeConfig.color} flex items-center justify-center text-white`}>
            {typeConfig.icon}
          </div>
          <span className="font-semibold text-foreground">{typeConfig.label}</span>
          {activity.feedback_type && (
            <Badge variant={activity.feedback_type === 'positive' ? 'default' : 'secondary'} className="text-xs">
              {activity.feedback_type === 'positive' ? '✨ Positivo' : 
               activity.feedback_type === 'constructive' ? '🔨 Construtivo' : '⚪ Neutro'}
            </Badge>
          )}
          {activity.feedback_category && (
            <Badge variant="outline" className="text-xs">
              {translateFeedbackCategory(activity.feedback_category)}
            </Badge>
          )}
        </div>
        
        {/* Date and Actions */}
        <div className="flex items-center gap-2">
          <div 
            className="text-sm text-muted-foreground cursor-help flex-shrink-0" 
            title={formatDateExact(activity.created_at)}
          >
            {formatDateRelative(activity.created_at)}
          </div>
          
          {/* Actions Dropdown - only show if callbacks are provided */}
          {(onEdit || onDelete) && (
            <DropdownMenu>
              <DropdownMenuTrigger asChild>
                <Button variant="ghost" size="sm" className="h-8 w-8 p-0">
                  <MoreVertical className="h-4 w-4" />
                </Button>
              </DropdownMenuTrigger>
              <DropdownMenuContent align="end">
                {onEdit && (
                  <DropdownMenuItem onClick={() => onEdit(activity)}>
                    <Edit3 className="mr-2 h-4 w-4" />
                    Editar
                  </DropdownMenuItem>
                )}
                {onDelete && (
                  <DropdownMenuItem 
                    onClick={() => onDelete(activity)}
                    className="text-destructive focus:text-destructive"
                  >
                    <Trash2 className="mr-2 h-4 w-4" />
                    Excluir
                  </DropdownMenuItem>
                )}
              </DropdownMenuContent>
            </DropdownMenu>
          )}
        </div>
      </div>

      {/* Content */}
      <div className="pl-8">
        <div 
          ref={contentRef}
          className={`text-foreground leading-relaxed whitespace-pre-wrap transition-all duration-300 ${
            shouldShowExpand && !isExpanded 
              ? 'overflow-hidden' 
              : ''
          }`}
          style={{
            maxHeight: shouldShowExpand && !isExpanded && isClient
              ? (typeof window !== 'undefined' && window.innerWidth < 640 ? '120px' : 
                 typeof window !== 'undefined' && window.innerWidth < 1024 ? '140px' : '160px')
              : 'none'
          }}
          dangerouslySetInnerHTML={{ __html: processedContent }}
        />
        
        {shouldShowExpand && (
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