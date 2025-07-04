'use client'

import React, { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { AppHeader } from '@/components/layout/AppHeader'
import { PersonCard } from '@/components/person/PersonCard'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Calendar, Clock, TrendingUp, Users } from 'lucide-react'
import { useActiveCompany, useLoadCompanies } from '@/lib/stores/companyStore'
import { useAllPeopleFromStore, useAllAISuggestions, useLoadPeopleData } from '@/lib/stores/peopleStore'
import { Person } from '@/lib/types'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

export default function Dashboard() {
  const router = useRouter()
  const activeCompany = useActiveCompany()
  const loadCompanies = useLoadCompanies()
  const loadPeopleData = useLoadPeopleData()
  
  const allPeople = useAllPeopleFromStore()
  const aiSuggestions = useAllAISuggestions()
  
  // Filter people by active company
  const people = React.useMemo(() => {
    if (!activeCompany) return []
    return allPeople.filter(person => person.companyId === activeCompany.id)
  }, [allPeople, activeCompany])
  
  // Calculate upcoming 1:1s locally
  const upcomingOneOnOnes = React.useMemo(() => {
    if (!activeCompany) return []
    
    const companyPeople = allPeople.filter(person => person.companyId === activeCompany.id)
    const upcomingMeetings: Array<{
      id: string
      personId: string
      date: Date
      notes: string
      aiSuggestions: string[]
      mentions: any[]
      status: 'scheduled'
      person: Person
    }> = []
    
    companyPeople.forEach(person => {
      if (person.nextOneOnOne && person.nextOneOnOne > new Date()) {
        upcomingMeetings.push({
          id: `upcoming-${person.id}`,
          personId: person.id,
          date: person.nextOneOnOne,
          notes: '',
          aiSuggestions: aiSuggestions
            .filter(s => s.personId === person.id && !s.isUsed)
            .map(s => s.content),
          mentions: [],
          status: 'scheduled' as const,
          person
        })
      }
    })
    
    return upcomingMeetings.sort((a, b) => a.date.getTime() - b.date.getTime())
  }, [activeCompany, allPeople, aiSuggestions])

  // Calculate additional metrics - moved up to ensure consistent hook ordering
  const oldestMeeting = React.useMemo(() => {
    if (people.length === 0) return null
    
    const lastMeetings = people.map(person => {
      // For now, use mock data for last meeting. In real app, this would come from sessions
      const daysSinceLastMeeting = Math.floor(Math.random() * 30) + 1
      return {
        person,
        daysSince: daysSinceLastMeeting
      }
    })
    
    return lastMeetings.reduce((oldest, current) => 
      current.daysSince > oldest.daysSince ? current : oldest
    )
  }, [people])

  const averageDaysBetweenMeetings = React.useMemo(() => {
    if (people.length === 0) return 0
    // Mock calculation - in real app this would be based on actual meeting history
    return Math.floor(Math.random() * 14) + 7 // 7-21 days average
  }, [people])

  useEffect(() => {
    loadCompanies()
    loadPeopleData()
  }, [loadCompanies, loadPeopleData])

  if (!activeCompany) {
    return (
      <div className="min-h-screen bg-slate-50">
        <AppHeader />
        <main className="container mx-auto px-6 py-8">
          <div className="flex items-center justify-center h-64">
            <div className="text-center">
              <h2 className="text-lg font-semibold text-slate-900 mb-2">
                Carregando...
              </h2>
              <p className="text-slate-600">
                Configurando sua empresa
              </p>
            </div>
          </div>
        </main>
      </div>
    )
  }

  const todayMeetings = upcomingOneOnOnes.filter(meeting => {
    const today = new Date()
    return meeting.date.toDateString() === today.toDateString()
  })

  const tomorrowMeetings = upcomingOneOnOnes.filter(meeting => {
    const tomorrow = new Date()
    tomorrow.setDate(tomorrow.getDate() + 1)
    return meeting.date.toDateString() === tomorrow.toDateString()
  })

  return (
    <div className="min-h-screen bg-background">
      <AppHeader />
      
      <main className="container mx-auto px-6 py-8">
        {/* Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold mb-2">
            Dashboard
          </h1>
          <p className="text-muted-foreground">
            Bem-vindo de volta! Aqui está um resumo do seu time na {activeCompany.name}.
          </p>
        </div>

        {/* Bento Grid Layout */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
          {/* Primary Metrics - Span wider on larger screens */}
          <Card className="md:col-span-1">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <p className="text-sm font-medium text-muted-foreground">
                    Total de Pessoas
                  </p>
                  <p className="text-2xl font-bold">
                    {people.length}
                  </p>
                </div>
                <div className="p-2 bg-blue-500/10 rounded-lg">
                  <Users className="h-6 w-6 text-blue-600" />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* 1:1s realizados este mês */}
          <Card className="md:col-span-1">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <p className="text-sm font-medium text-muted-foreground">
                    1:1s este mês
                  </p>
                  <div className="flex items-baseline gap-2">
                    <p className="text-2xl font-bold text-green-600">
                      {Math.floor(Math.random() * 15) + 3}
                    </p>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Reuniões registradas
                  </p>
                </div>
                <div className="p-2 bg-green-500/10 rounded-lg">
                  <Calendar className="h-6 w-6 text-green-600" />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Média de dias entre 1:1s */}
          <Card className="md:col-span-1">
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <p className="text-sm font-medium text-muted-foreground">
                    Frequência média
                  </p>
                  <div className="flex items-baseline gap-2">
                    <p className="text-2xl font-bold">
                      {averageDaysBetweenMeetings}
                    </p>
                    <span className="text-sm text-muted-foreground">dias</span>
                  </div>
                  <p className="text-xs text-muted-foreground">
                    Entre 1:1s
                  </p>
                </div>
                <div className="p-2 bg-blue-500/10 rounded-lg">
                  <Clock className="h-6 w-6 text-blue-600" />
                </div>
              </div>
            </CardContent>
          </Card>

          {/* Alert metric - Oldest meeting */}
          <Card className={`md:col-span-1 ${oldestMeeting && oldestMeeting.daysSince > 14 ? 'ring-2 ring-red-500/20 bg-red-50/50 dark:bg-red-950/10' : ''}`}>
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <div className="space-y-1">
                  <p className="text-sm font-medium text-muted-foreground">
                    Última reunião
                  </p>
                  <div className="flex items-baseline gap-2">
                    <p className={`text-2xl font-bold ${oldestMeeting && oldestMeeting.daysSince > 14 ? 'text-red-600' : ''}`}>
                      {oldestMeeting ? oldestMeeting.daysSince : 0}
                    </p>
                    <span className="text-sm text-muted-foreground">dias</span>
                  </div>
                  {oldestMeeting && (
                    <p className="text-xs text-muted-foreground truncate">
                      {oldestMeeting.person.name}
                      {oldestMeeting.daysSince > 14 && (
                        <Badge variant="destructive" className="ml-2 text-xs">
                          Negligenciado
                        </Badge>
                      )}
                    </p>
                  )}
                </div>
                <div className={`p-2 rounded-lg ${oldestMeeting && oldestMeeting.daysSince > 14 ? 'bg-red-500/10' : 'bg-blue-500/10'}`}>
                  <TrendingUp className={`h-6 w-6 ${oldestMeeting && oldestMeeting.daysSince > 14 ? 'text-red-600' : 'text-blue-600'}`} />
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Team Focus Layout */}
        <div className="space-y-8">
          {/* Team Members - Main Focus */}
          <div>
            <div className="flex items-center justify-between mb-6">
              <div>
                <h2 className="text-2xl font-bold mb-1">
                  Seu Time
                </h2>
                <p className="text-muted-foreground">
                  Gerencie 1:1s e acompanhe o desenvolvimento de {people.length} pessoa{people.length !== 1 ? 's' : ''} na {activeCompany.name}
                </p>
              </div>
              <Button className="gap-2" size="lg">
                <Users className="h-4 w-4" />
                Adicionar pessoa
              </Button>
            </div>

            {people.length === 0 ? (
              <Card>
                <CardContent className="p-12 text-center">
                  <div className="p-4 bg-muted/30 rounded-full w-fit mx-auto mb-4">
                    <Users className="h-8 w-8 text-muted-foreground" />
                  </div>
                  <h3 className="font-semibold mb-2">
                    Comece a construir seu time
                  </h3>
                  <p className="text-muted-foreground mb-6 max-w-md mx-auto">
                    Adicione as pessoas do seu time para começar a registrar 1:1s e acompanhar o desenvolvimento de cada uma.
                  </p>
                  <Button className="gap-2" size="lg">
                    <Users className="h-4 w-4" />
                    Adicionar primeira pessoa
                  </Button>
                </CardContent>
              </Card>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {people.map((person) => (
                  <PersonCard
                    key={person.id}
                    person={person}
                    onClick={() => {
                      router.push(`/profile/${person.id}`)
                    }}
                  />
                ))}
              </div>
            )}
          </div>

          {/* Quick Stats Summary */}
          {people.length > 0 && (
            <Card>
              <CardHeader>
                <CardTitle className="text-lg">Resumo de Atividade</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                  <div className="text-center">
                    <p className="text-2xl font-bold text-green-600">{people.length}</p>
                    <p className="text-sm text-muted-foreground">Pessoas no time</p>
                  </div>
                  <div className="text-center">
                    <p className="text-2xl font-bold text-blue-600">
                      {Math.floor(Math.random() * 20) + 5}
                    </p>
                    <p className="text-sm text-muted-foreground">1:1s este mês</p>
                  </div>
                  <div className="text-center">
                    <p className="text-2xl font-bold text-orange-600">
                      {oldestMeeting ? oldestMeeting.daysSince : 0}
                    </p>
                    <p className="text-sm text-muted-foreground">Dias desde última reunião</p>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </main>
    </div>
  )
}