'use client'

import React from 'react'
import { Person } from '@/lib/types'
import { PersonMentions } from '@/components/timeline/PersonMentions'
import { useActiveCompany } from '@/lib/stores/companyStore'

interface PersonFeedbackTabProps {
  person: Person
  allPeople: Person[]
}

export function PersonFeedbackTab({ person, allPeople }: PersonFeedbackTabProps) {
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
            Feedbacks Recebidos
          </h2>
          <p className="text-sm text-muted-foreground">
            Todas as menções de {person.name} em conversas e 1:1s
          </p>
        </div>
      </div>
      
      <PersonMentions 
        person={person}
        companyId={activeCompany.uuid}
        allPeople={allPeople}
      />
    </div>
  )
}