'use client'

import React from 'react'
import { Person } from '@/lib/types'
import { UnifiedTimeline } from '@/components/timeline/UnifiedTimeline'
import { useActiveCompany } from '@/lib/stores/companyStore'

interface PersonTimelineTabProps {
  person: Person
  allPeople: Person[]
  onCreatePerson?: (mentionName: string) => void
}

export function PersonTimelineTab({ person, allPeople }: PersonTimelineTabProps) {
  const activeCompany = useActiveCompany()

  if (!activeCompany) {
    return (
      <div className="text-center py-8">
        <p className="text-muted-foreground">
          Carregando informações da empresa...
        </p>
      </div>
    )
  }

  return (
    <UnifiedTimeline 
      person={person}
      companyId={activeCompany.uuid}
      allPeople={allPeople}
    />
  )
}