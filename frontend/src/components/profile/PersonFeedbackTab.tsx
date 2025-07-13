'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { Person } from '@/lib/types'
import { formatTimeAgo } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'
import { useAllFeedbacks } from '@/lib/stores/peopleStore'

interface PersonFeedbackTabProps {
  person: Person
  allPeople: Person[]
}

export function PersonFeedbackTab({ person, allPeople }: PersonFeedbackTabProps) {
  const allFeedbacks = useAllFeedbacks()
  const personFeedbacks = allFeedbacks.filter(feedback => feedback.personId === person.id)

  return (
    <Card>
      <CardHeader>
        <CardTitle>Feedbacks Recebidos</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {personFeedbacks.length > 0 ? (
            personFeedbacks.map((feedback) => {
              const sourcePerson = feedback.sourcePersonId 
                ? allPeople.find(p => p.id === feedback.sourcePersonId)
                : null
              
              return (
                <div key={feedback.id} className="border rounded-lg p-4">
                  <div className="flex items-center justify-between mb-2">
                    <div className="flex items-center gap-2">
                      <Avatar className="h-8 w-8">
                        <AvatarFallback className="bg-primary/10 text-primary text-xs">
                          {sourcePerson ? getInitials(sourcePerson.name) : 'LI'}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <p className="font-medium text-sm">
                          {sourcePerson?.name || 'Líder'}
                        </p>
                        <p className="text-xs text-muted-foreground">
                          {feedback.source === 'mention' 
                            ? `Via @menção em 1:1` 
                            : 'Feedback direto'
                          }
                        </p>
                      </div>
                    </div>
                    <span className="text-xs text-muted-foreground">
                      {formatTimeAgo(feedback.createdAt)}
                    </span>
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {feedback.source === 'mention' 
                      ? `${sourcePerson?.name || 'Alguém'} disse: "${feedback.content}"`
                      : feedback.content
                    }
                  </p>
                  <div className="flex gap-2 mt-2">
                    <Badge 
                      variant="outline" 
                      className={`text-xs ${
                        feedback.type === 'positive' ? 'border-green-500 text-green-700' :
                        feedback.type === 'negative' ? 'border-red-500 text-red-700' :
                        'border-gray-500 text-gray-700'
                      }`}
                    >
                      {feedback.type === 'positive' ? '👍 Positivo' :
                       feedback.type === 'negative' ? '👎 Atenção' :
                       '💭 Neutro'}
                    </Badge>
                    {feedback.source === 'mention' && (
                      <Badge variant="secondary" className="text-xs">
                        Menção cruzada
                      </Badge>
                    )}
                  </div>
                </div>
              )
            })
          ) : (
            <div className="text-center py-8">
              <p className="text-muted-foreground">
                Nenhum feedback encontrado ainda
              </p>
            </div>
          )}
        </div>
      </CardContent>
    </Card>
  )
}