'use client'

import React from 'react'
import { useRouter } from 'next/navigation'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Person } from '@/lib/types'
import { getSessionsByPerson } from '@/lib/data/mockData'
import { formatTimeAgo, formatShortDate } from '@/lib/utils/dates'

interface PersonHistoryTabProps {
  person: Person
  allPeople: Person[]
  onCreatePerson?: (mentionName: string) => void
}

export function PersonHistoryTab({ person, allPeople, onCreatePerson }: PersonHistoryTabProps) {
  const router = useRouter()
  const personSessions = getSessionsByPerson(person.id)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Histórico de 1:1s de {person.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {personSessions.length > 0 ? (
            personSessions
              .sort((a, b) => b.date.getTime() - a.date.getTime())
              .map((session) => (
                <div 
                  key={session.id} 
                  className="border-l-2 border-blue-500 pl-4 pb-4"
                >
                  <div className="flex items-center justify-between mb-2">
                    <h4 className="font-medium">
                      1:1 com {person.name}
                    </h4>
                    <span className="text-sm text-muted-foreground">
                      {formatTimeAgo(session.date)}
                    </span>
                  </div>
                  
                  <p className="text-sm text-muted-foreground mb-2">
                    {session.notes}
                  </p>
                  
                  {/* Show mentions made in this session */}
                  {session.mentions && session.mentions.length > 0 && (
                    <div className="mt-3 p-2 bg-muted/30 rounded-md">
                      <p className="text-xs font-medium text-muted-foreground mb-1">
                        Pessoas mencionadas nesta reunião:
                      </p>
                      <div className="flex flex-wrap gap-1">
                        {session.mentions.map((mention) => {
                          const mentionedPerson = allPeople.find(p => 
                            p.id === mention.mentionedPersonId && 
                            p.companyId === person.companyId  // Only same company
                          )
                          
                          return (
                            <Badge 
                              key={mention.id} 
                              variant="outline" 
                              className={`text-xs ${mentionedPerson ? 'cursor-pointer hover:bg-primary/10' : 'cursor-pointer opacity-60 hover:opacity-80'}`}
                              onClick={() => {
                                if (mentionedPerson) {
                                  router.push(`/profile/${mentionedPerson.id}`)
                                } else {
                                  // Trigger create person dialog
                                  onCreatePerson?.(mention.mentionedPersonName)
                                }
                              }}
                            >
                              @{mention.mentionedPersonName}
                              {!mentionedPerson && ' (?)'}
                            </Badge>
                          )
                        })}
                      </div>
                    </div>
                  )}
                  
                  <div className="flex gap-2 mt-3">
                    <Badge variant="secondary" className="text-xs">
                      {formatShortDate(session.date)}
                    </Badge>
                    <Badge variant="secondary" className="text-xs">
                      1:1 Reunião
                    </Badge>
                    {session.status === 'completed' && (
                      <Badge variant="outline" className="text-xs text-green-600">
                        Concluída
                      </Badge>
                    )}
                  </div>
                </div>
              ))
          ) : (
            <div className="text-center py-8">
              <p className="text-muted-foreground">
                Nenhuma reunião 1:1 registrada ainda
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}