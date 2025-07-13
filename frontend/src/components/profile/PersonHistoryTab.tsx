'use client'

import React from 'react'
import { Person } from '@/lib/types'
import { PersonTimeline } from '@/components/timeline/PersonTimeline'
import { useActiveCompany } from '@/lib/stores/companyStore'

interface PersonHistoryTabProps {
  person: Person
  allPeople: Person[]
  onCreatePerson?: (mentionName: string) => void
}

export function PersonHistoryTab({ person, allPeople, onCreatePerson }: PersonHistoryTabProps) {
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
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">
            Histórico de {person.name}
          </h2>
          <p className="text-sm text-muted-foreground">
            Timeline completa de 1:1s e anotações
          </p>
        </div>
      </div>
      
      <PersonTimeline 
        person={person}
        companyId={activeCompany.uuid}
        allPeople={allPeople}
      />
    </div>
  )
}