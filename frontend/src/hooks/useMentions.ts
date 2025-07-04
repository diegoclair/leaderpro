import { useState } from 'react'
import { Person } from '@/lib/types'

interface MentionDetection {
  name: string
  person: Person
  context: string
}

interface UseMentionsProps {
  allPeople: Person[]
  currentPersonCompanyId?: string
}

export function useMentions({ allPeople, currentPersonCompanyId }: UseMentionsProps) {
  const [showMentionSuggestions, setShowMentionSuggestions] = useState(false)
  const [mentionQuery, setMentionQuery] = useState('')

  const detectMentions = (text: string): MentionDetection[] => {
    const mentionRegex = /@(\w+)/g
    const mentions: MentionDetection[] = []
    let match

    while ((match = mentionRegex.exec(text)) !== null) {
      const mentionedName = match[1]
      const mentionedPerson = allPeople.find(p => 
        p.name.toLowerCase().includes(mentionedName.toLowerCase()) &&
        p.companyId === currentPersonCompanyId // Only same company
      )
      if (mentionedPerson) {
        mentions.push({
          name: mentionedName,
          person: mentionedPerson,
          context: text
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

    return `${beforeAt}@${selectedPerson.name}${afterMention}`
  }

  return {
    showMentionSuggestions,
    mentionQuery,
    detectMentions,
    handleTextChange,
    getFilteredPeople,
    insertMention,
    setShowMentionSuggestions
  }
}