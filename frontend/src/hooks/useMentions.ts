import { useState } from 'react'
import { Person } from '@/lib/types'

interface MentionDetection {
  uuid: string
  name: string
  person: Person
}

interface UseMentionsProps {
  allPeople: Person[]
  currentPersonCompanyId?: string
}

export function useMentions({ allPeople, currentPersonCompanyId }: UseMentionsProps) {
  const [showMentionSuggestions, setShowMentionSuggestions] = useState(false)
  const [mentionQuery, setMentionQuery] = useState('')

  const detectMentions = (text: string): MentionDetection[] => {
    const tokenRegex = /\{\{person:([^\|]+)\|([^}]+)\}\}/g
    const mentions: MentionDetection[] = []
    let match

    while ((match = tokenRegex.exec(text)) !== null) {
      const uuid = match[1]
      const name = match[2]
      const mentionedPerson = allPeople.find(p => 
        p.id === uuid &&
        p.companyId === currentPersonCompanyId // Only same company
      )
      if (mentionedPerson) {
        mentions.push({
          uuid,
          name,
          person: mentionedPerson
        })
      }
    }

    return mentions
  }

  const handleTextChange = (
    value: string,
    onMentionNotFound: (mentionName: string) => void
  ) => {
    // Detect @ mentions
    const lastAtIndex = value.lastIndexOf('@')
    if (lastAtIndex !== -1) {
      const textAfterAt = value.substring(lastAtIndex + 1)
      const spaceIndex = textAfterAt.indexOf(' ')
      const query = spaceIndex === -1 ? textAfterAt : textAfterAt.substring(0, spaceIndex)

      // Check if user just typed space after @name (mention completion)
      if (spaceIndex !== -1 && query.length > 0) {
        // User typed @name + space, check if person exists
        const mentionedPerson = allPeople.find(p => 
          p.name.toLowerCase().includes(query.toLowerCase()) &&
          p.companyId === currentPersonCompanyId
        )

        if (!mentionedPerson) {
          // Person doesn't exist, trigger callback
          onMentionNotFound(query)
        }

        // Hide suggestions after space
        setShowMentionSuggestions(false)
      } else {
        // Show suggestions immediately when @ is typed, even with empty query
        setMentionQuery(query)
        setShowMentionSuggestions(true)
      }
    } else {
      setShowMentionSuggestions(false)
    }
  }

  const getFilteredPeople = (excludePersonId?: string) => {
    return allPeople.filter(p => 
      p.id !== excludePersonId && // Don't suggest the current person
      p.companyId === currentPersonCompanyId && // Only same company
      (mentionQuery === '' || p.name.toLowerCase().includes(mentionQuery.toLowerCase())) // Show all if empty query
    )
  }

  const insertMention = (text: string, selectedPerson: Person): string => {
    const lastAtIndex = text.lastIndexOf('@')
    const beforeAt = text.substring(0, lastAtIndex)
    const afterAt = text.substring(lastAtIndex + 1)
    const spaceIndex = afterAt.indexOf(' ')
    const afterMention = spaceIndex === -1 ? '' : afterAt.substring(spaceIndex)

    // Insert token with uuid and name: {{person:uuid|name}}
    const token = `{{person:${selectedPerson.id}|${selectedPerson.name}}}`
    return `${beforeAt}${token}${afterMention}`
  }

  const renderTokensForDisplay = (text: string): string => {
    // Replace {{person:uuid|name}} tokens with @name for display
    return text.replace(/\{\{person:[^|]+\|([^}]+)\}\}/g, '@$1')
  }

  return {
    showMentionSuggestions,
    mentionQuery,
    detectMentions,
    handleTextChange,
    getFilteredPeople,
    insertMention,
    renderTokensForDisplay,
    setShowMentionSuggestions
  }
}