'use client'

import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Person } from '@/lib/types'
import { getInitials } from '@/lib/utils/names'

interface MentionSuggestionsProps {
  show: boolean
  people: Person[]
  onSelect: (person: Person) => void
}

export default function MentionSuggestions({ show, people, onSelect }: MentionSuggestionsProps) {
  if (!show || people.length === 0) return null

  return (
    <div className="absolute z-10 w-full mt-1 bg-background border rounded-md shadow-lg max-h-32 overflow-y-auto">
      {people.slice(0, 5).map((suggestedPerson) => (
        <button
          key={suggestedPerson.id}
          onClick={() => onSelect(suggestedPerson)}
          className="w-full px-3 py-2 text-left hover:bg-muted flex items-center gap-2"
        >
          <Avatar className="h-6 w-6">
            <AvatarImage src={suggestedPerson.avatar} alt={suggestedPerson.name} />
            <AvatarFallback className="text-xs">
              {getInitials(suggestedPerson.name)}
            </AvatarFallback>
          </Avatar>
          <div>
            <p className="text-sm font-medium">{suggestedPerson.name}</p>
            <p className="text-xs text-muted-foreground">{suggestedPerson.role}</p>
          </div>
        </button>
      ))}
    </div>
  )
}