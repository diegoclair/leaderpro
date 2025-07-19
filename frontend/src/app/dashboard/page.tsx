'use client'

import React, { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { PersonCard } from '@/components/person/PersonCard'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Calendar, Clock, TrendingUp, Users } from 'lucide-react'
import { useActiveCompany, useLoadCompanies, useCompanyStore } from '@/lib/stores/companyStore'
import { useAllPeopleFromStore, useAllAISuggestions, useLoadPeopleData, useLoadPeopleFromAPI, useLoadDashboardData, useDashboardStats } from '@/lib/stores/peopleStore'
import { Person } from '@/lib/types'
import { AppHeader } from '@/components/layout/AppHeader'
import { useAuthRedirect } from '@/hooks/useAuthRedirect'
import { OnboardingWizard } from '@/components/onboarding/OnboardingWizard'
import PersonModal, { PersonFormData } from '@/components/person/PersonModal'
import { apiClient } from '@/lib/api/client'
import { useNotificationStore } from '@/lib/stores/notificationStore'

export default function Dashboard() {
  const { isLoading, shouldRender, needsOnboarding, completeOnboarding } = useAuthRedirect({ requireAuth: true })
  const router = useRouter()
  const activeCompany = useActiveCompany()
  const companies = useCompanyStore(state => state.companies)
  const loadCompanies = useLoadCompanies()
  const loadPeopleData = useLoadPeopleData()
  const loadPeopleFromAPI = useLoadPeopleFromAPI()
  const loadDashboardData = useLoadDashboardData()
  const dashboardStats = useDashboardStats()
  const { showSuccess, showError } = useNotificationStore()
  
  // Modal state
  const [isAddPersonModalOpen, setIsAddPersonModalOpen] = useState(false)
  const [isCreatingPerson, setIsCreatingPerson] = useState(false)

  // Carregar empresas uma única vez quando o componente monta
  useEffect(() => {
    // Aguardar o shouldRender estar true (significa que auth foi carregado)
    if (shouldRender && companies.length === 0) {
      loadCompanies()
    }
  }, [shouldRender, companies.length]) // Depende do shouldRender para garantir auth carregado
  
  const allPeople = useAllPeopleFromStore()
  const aiSuggestions = useAllAISuggestions()
  
  // Filter people by active company
  const people = React.useMemo(() => {
    if (!activeCompany) return []
    return allPeople.filter(person => person.companyId === activeCompany.uuid)
  }, [allPeople, activeCompany])
  
  // Calculate upcoming 1:1s locally
  const upcomingOneOnOnes = React.useMemo(() => {
    if (!activeCompany) return []
    
    const companyPeople = allPeople.filter(person => person.companyId === activeCompany.uuid)
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


  useEffect(() => {
    console.log('🔍 Dashboard useEffect executou')
    // Carregar dashboard quando estiver pronto (após empresas carregadas)
    if (shouldRender && !needsOnboarding && activeCompany) {
      console.log('📞 Dashboard chamando loadDashboardData')
      loadDashboardData(activeCompany.uuid)
    }
  }, [loadDashboardData, shouldRender, needsOnboarding, activeCompany])

  // Function to create person
  const handleCreatePerson = async (personData: PersonFormData): Promise<boolean> => {
    if (!activeCompany) {
      showError('Erro', 'Nenhuma empresa selecionada')
      return false
    }

    setIsCreatingPerson(true)
    try {
      // Convert form data to API format
      const apiData = {
        name: personData.name,
        email: personData.email || undefined,
        position: personData.position || undefined,
        department: personData.department || undefined,
        phone: personData.phone || undefined,
        start_date: personData.start_date || undefined,
        notes: personData.notes || undefined
      }

      // Remove undefined fields
      Object.keys(apiData).forEach(key => {
        if (apiData[key as keyof typeof apiData] === undefined) {
          delete apiData[key as keyof typeof apiData]
        }
      })

      await apiClient.authPost(`/companies/${activeCompany.uuid}/people`, apiData)
      
      showSuccess('Sucesso', `${personData.name} foi adicionado(a) ao seu time!`)
      
      // Reload dashboard data to show the new person and updated stats
      await loadDashboardData(activeCompany.uuid)
      
      return true
    } catch (error) {
      console.error('Error creating person:', error)
      showError('Erro', 'Não foi possível adicionar a pessoa. Tente novamente.')
      return false
    } finally {
      setIsCreatingPerson(false)
    }
  }

  // Mostrar loading se estiver carregando auth ou não deve renderizar
  if (isLoading || !shouldRender) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  // Mostrar onboarding se necessário
  if (needsOnboarding) {
    return <OnboardingWizard onComplete={completeOnboarding} />
  }

  if (!activeCompany) {
    return (
      <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
        <AppHeader />
        <main className="container mx-auto px-4 py-6">
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
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <AppHeader />
      <main className="container mx-auto px-4 py-6">
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
                  {dashboardStats?.totalPeople || people.length}
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
                    {dashboardStats?.oneOnOnesCountThisMonth || 0}
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
                    {Math.round(dashboardStats?.averageDaysBetweenOneOnOnes || 0)}
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

        {/* Alert metric - Last meeting */}
        <Card className="md:col-span-1">
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div className="space-y-1">
                <p className="text-sm font-medium text-muted-foreground">
                  Última reunião
                </p>
                <div className="flex items-baseline gap-2">
                  <p className="text-2xl font-bold">
                    {dashboardStats?.lastMeetingDate 
                      ? (() => {
                          const now = new Date()
                          const daysDifference = Math.floor((now.getTime() - dashboardStats.lastMeetingDate.getTime()) / (1000 * 60 * 60 * 24))
                          
                          // Debug logging for dashboard stats
                          console.log('🔍 Dashboard stats debug:', {
                            input: dashboardStats.lastMeetingDate,
                            inputISO: dashboardStats.lastMeetingDate.toISOString(),
                            inputTime: dashboardStats.lastMeetingDate.getTime(),
                            now: now,
                            nowISO: now.toISOString(),
                            nowTime: now.getTime(),
                            timeDiffMs: now.getTime() - dashboardStats.lastMeetingDate.getTime(),
                            daysDifference,
                            context: 'dashboard-stats'
                          })
                          
                          return daysDifference
                        })()
                      : '-'
                    }
                  </p>
                  {dashboardStats?.lastMeetingDate && (
                    <span className="text-sm text-muted-foreground">dias atrás</span>
                  )}
                </div>
                <p className="text-xs text-muted-foreground">
                  {dashboardStats?.lastMeetingDate ? 'Registrada no sistema' : 'Nenhuma registrada'}
                </p>
              </div>
              <div className="p-2 bg-blue-500/10 rounded-lg">
                <TrendingUp className="h-6 w-6 text-blue-600" />
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
            {people.length > 0 && (
              <Button 
                className="gap-2" 
                size="lg"
                onClick={() => setIsAddPersonModalOpen(true)}
              >
                <Users className="h-4 w-4" />
                Adicionar pessoa
              </Button>
            )}
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
                <Button 
                  className="gap-2" 
                  size="lg"
                  onClick={() => setIsAddPersonModalOpen(true)}
                >
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

      </div>
      </main>

      {/* Add Person Modal */}
      <PersonModal
        open={isAddPersonModalOpen}
        onClose={() => setIsAddPersonModalOpen(false)}
        mode="create"
        onSubmit={handleCreatePerson}
        isLoading={isCreatingPerson}
      />
    </div>
  )
}