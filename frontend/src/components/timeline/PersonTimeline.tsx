'use client'

import React, { useEffect, useState } from 'react'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Calendar, MessageSquare, MessageCircle, Clock, User } from 'lucide-react'
import { Person } from '@/lib/types'
import { formatTimeAgo } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'
import { apiClient } from '@/lib/stores/authStore'
import { LoadingSpinner } from '@/components/ui/loading-spinner'
import { 
  NOTE_SOURCE_TYPES, 
  getNoteSourceTypeLabel, 
  getFeedbackTypeLabel, 
  getFeedbackCategoryLabel, 
  getFeedbackTypeColor 
} from '@/lib/constants/notes'

interface TimelineNote {
  uuid: string
  type: 'note' | 'mention'  // Backend format: "note" or "mention"
  source_type: 'one_on_one' | 'feedback' | 'observation'  // Real note type
  content: string
  author_name: string
  feedback_type?: 'positive' | 'constructive' | 'neutral'
  feedback_category?: 'performance' | 'behavior' | 'skill' | 'collaboration'
  created_at: string
  source_person_name?: string  // For mentions
}

interface PersonTimelineProps {
  person: Person
  companyId: string
  allPeople: Person[]
  className?: string
}

const getNoteIcon = (sourceType: string) => {
  switch (sourceType) {
    case NOTE_SOURCE_TYPES.ONE_ON_ONE:
      return <Calendar className="h-4 w-4" />
    case NOTE_SOURCE_TYPES.FEEDBACK:
      return <MessageCircle className="h-4 w-4" />
    default:
      return <MessageSquare className="h-4 w-4" />
  }
}

// Function to render content with clickable mention links
const renderContentWithMentions = (
  content: string, 
  mentions: TimelineNote['mentions'], 
  allPeople: Person[]
) => {
  // Extract mentions from content using regex pattern {{person:uuid|name}}
  const mentionRegex = /\{\{person:([^|]+)\|([^}]+)\}\}/g
  const parts: React.ReactNode[] = []
  let lastIndex = 0
  let match
  let mentionIndex = 0

  while ((match = mentionRegex.exec(content)) !== null) {
    const [fullMatch, personUuid, personName] = match
    const startIndex = match.index
    const endIndex = startIndex + fullMatch.length

    // Add text before this mention
    if (startIndex > lastIndex) {
      parts.push(
        <span key={`text-before-${mentionIndex}`}>
          {content.slice(lastIndex, startIndex)}
        </span>
      )
    }

    // Find the person to link to (try both uuid and id)
    const mentionedPerson = allPeople.find(p => p.uuid === personUuid || p.id === personUuid)
    
    if (mentionedPerson) {
      parts.push(
        <button
          key={`mention-${mentionIndex}`}
          onClick={() => {
            // Navigate to person's profile
            window.location.href = `/profile/${mentionedPerson.uuid || mentionedPerson.id}`
          }}
          className="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium text-blue-600 hover:text-blue-800 bg-blue-50 hover:bg-blue-100 dark:bg-blue-900/20 dark:text-blue-400 dark:hover:text-blue-300 rounded-md transition-colors"
        >
          <User className="h-3 w-3" />
          @{personName}
        </button>
      )
    } else {
      // Fallback if person not found
      parts.push(
        <span key={`mention-fallback-${mentionIndex}`} className="text-blue-600 font-medium">
          @{personName}
        </span>
      )
    }

    lastIndex = endIndex
    mentionIndex++
  }

  // Add remaining text after last mention
  if (lastIndex < content.length) {
    parts.push(
      <span key="text-after">
        {content.slice(lastIndex)}
      </span>
    )
  }

  // If no mentions found, return plain content
  if (parts.length === 0) {
    return <p className="text-sm leading-relaxed">{content}</p>
  }

  return <p className="text-sm leading-relaxed">{parts}</p>
}

export function PersonTimeline({ person, companyId, allPeople, className }: PersonTimelineProps) {
  const [notes, setNotes] = useState<TimelineNote[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    const fetchNotes = async () => {
      setIsLoading(true)
      setError(null)
      
      try {
        const response = await apiClient.authGet(`/companies/${companyId}/people/${person.uuid}/timeline`)
        
        // The API returns { pagination: {...}, data: [...] }, so we need to access the data array
        const notesData = Array.isArray(response) ? response : response.data || []
        
        // Sort notes by creation date (newest first)
        const sortedNotes = notesData.sort((a: TimelineNote, b: TimelineNote) => 
          new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
        )
        
        setNotes(sortedNotes)
      } catch (err) {
        console.error('Error fetching notes:', err)
        setError('Não foi possível carregar o histórico de anotações')
      } finally {
        setIsLoading(false)
      }
    }

    fetchNotes()
  }, [person.uuid, companyId])

  if (isLoading) {
    return (
      <div className="flex items-center justify-center py-8">
        <LoadingSpinner />
      </div>
    )
  }

  if (error) {
    return (
      <Card className={className}>
        <CardContent className="p-6 text-center">
          <p className="text-sm text-red-600">{error}</p>
        </CardContent>
      </Card>
    )
  }

  if (notes.length === 0) {
    return (
      <Card className={className}>
        <CardContent className="p-8 text-center">
          <div className="p-4 bg-muted/30 rounded-full w-fit mx-auto mb-4">
            <MessageSquare className="h-6 w-6 text-muted-foreground" />
          </div>
          <h3 className="font-medium mb-2">Nenhuma anotação ainda</h3>
          <p className="text-sm text-muted-foreground">
            As reuniões 1:1, feedbacks e anotações sobre {person.name} aparecerão aqui.
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className={className}>
      <div className="space-y-4">
        {notes.map((note, index) => (
          <Card key={note.uuid} className="relative">
            <CardContent className="p-4">
              {/* Note Header */}
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className="p-2 bg-muted/50 rounded-lg">
                    {getNoteIcon(note.source_type)}
                  </div>
                  <div>
                    <div className="flex items-center gap-2 mb-1">
                      <Badge variant="outline" className="text-xs">
                        {getNoteSourceTypeLabel(note.source_type)}
                      </Badge>
                      
                      {/* Feedback type badge */}
                      {note.source_type === NOTE_SOURCE_TYPES.FEEDBACK && note.feedback_type && (
                        <Badge className={`text-xs ${getFeedbackTypeColor(note.feedback_type)}`}>
                          {getFeedbackTypeLabel(note.feedback_type)}
                        </Badge>
                      )}
                      
                      {/* Feedback category badge */}
                      {note.source_type === NOTE_SOURCE_TYPES.FEEDBACK && note.feedback_category && (
                        <Badge variant="secondary" className="text-xs">
                          {getFeedbackCategoryLabel(note.feedback_category)}
                        </Badge>
                      )}
                    </div>
                    
                    <div className="flex items-center gap-1 text-xs text-muted-foreground">
                      <Clock className="h-3 w-3" />
                      {formatTimeAgo(new Date(note.created_at))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Note Content */}
              <div className="ml-11">
                {renderContentWithMentions(note.content, note.mentions, allPeople)}
              </div>
              
              {/* Timeline connector */}
              {index < notes.length - 1 && (
                <div className="absolute left-7 top-16 bottom-0 w-px bg-border"></div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  )
}