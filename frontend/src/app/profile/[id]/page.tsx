'use client'

import React, { useState, useEffect } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { AppHeader } from '@/components/layout/AppHeader'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Badge } from '@/components/ui/badge'
import { 
  ArrowLeft, 
  MapPin, 
  Calendar, 
  User, 
  Dog,
  Clock,
  Edit3
} from 'lucide-react'
import { useAllPeopleFromStore, useLoadPeopleFromAPI, useUpdatePerson } from '@/lib/stores/peopleStore'
import { useActiveCompany, useLoadCompanies, useCompanyStore } from '@/lib/stores/companyStore'
import { apiClient } from '@/lib/stores/authStore'
import { useAuthRedirect } from '@/hooks/useAuthRedirect'
import { Person } from '@/lib/types'
import { PersonInfoTab } from '@/components/profile/PersonInfoTab'
import { PersonHistoryTab } from '@/components/profile/PersonHistoryTab'
import { PersonFeedbackTab } from '@/components/profile/PersonFeedbackTab'
import { PersonChatTab } from '@/components/profile/PersonChatTab'
import { useCreatePerson } from '@/hooks/useCreatePerson'
import CreatePersonDialog from '@/components/profile/CreatePersonDialog'
import PersonModal, { PersonFormData } from '@/components/person/PersonModal'
import { formatTimeAgoWithoutSuffix, getMockDaysAgo } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'

export default function ProfilePage() {
  const { isLoading: authLoading, shouldRender } = useAuthRedirect({ requireAuth: true })
  const params = useParams()
  const router = useRouter()
  const allPeople = useAllPeopleFromStore()
  const loadPeopleFromAPI = useLoadPeopleFromAPI()
  const updatePerson = useUpdatePerson()
  const activeCompany = useActiveCompany()
  const loadCompanies = useLoadCompanies()
  
  // Get tab from URL or default to 'info' (using window for static export compatibility)
  const [activeTab, setActiveTab] = useState('info')
  const [isLoading, setIsLoading] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  
  // All hooks must be called before any conditional returns
  const {
    showCreatePersonDialog,
    personToCreate,
    newPersonName,
    newPersonRole,
    setNewPersonName,
    setNewPersonRole,
    openCreateDialog,
    closeCreateDialog,
    handleCreatePerson
  } = useCreatePerson()
  
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search)
      const tab = urlParams.get('tab')
      if (tab) {
        setActiveTab(tab)
      }
    }
  }, [])

  // Load companies and people data if not available (happens on page refresh)
  useEffect(() => {
    const loadDataIfNeeded = async () => {
      if (!isLoading && shouldRender) {
        // First load companies if we don't have any
        const { companies } = useCompanyStore.getState()
        if (companies.length === 0) {
          try {
            await loadCompanies()
          } catch (error) {
            console.error('Error loading companies:', error)
            return
          }
        }

        // Then load people if we have an active company but no people
        if (activeCompany && allPeople.length === 0) {
          setIsLoading(true)
          try {
            await loadPeopleFromAPI(activeCompany.uuid)
          } catch (error) {
            console.error('Error loading people:', error)
          } finally {
            setIsLoading(false)
          }
        }
      }
    }

    loadDataIfNeeded()
  }, [activeCompany, allPeople.length, loadPeopleFromAPI, loadCompanies, isLoading, shouldRender])

  // Don't render anything until we have tried to load data
  if (authLoading || !shouldRender) {
    return null
  }

  // Update URL when tab changes
  const handleTabChange = (tab: string) => {
    setActiveTab(tab)
    const newUrl = `/profile/${params.id}?tab=${tab}`
    router.replace(newUrl, { scroll: false })
  }

  const handleEditPerson = async (personData: PersonFormData): Promise<boolean> => {
    if (!activeCompany || !person) return false

    try {
      // Call API to update person
      await apiClient.authPut(`/companies/${activeCompany.uuid}/people/${person.uuid}`, personData)

      // Update local store
      const updatedPersonData: Partial<Person> = {
        name: personData.name,
        email: personData.email,
        position: personData.position,
        department: personData.department,
        phone: personData.phone,
        startDate: personData.start_date ? new Date(personData.start_date.replace('T00:00:00Z', '')) : undefined,
        notes: personData.notes,
      }

      updatePerson(person.id, updatedPersonData)
      return true
    } catch (error) {
      console.error('Error updating person:', error)
      return false
    }
  }
  
  const person = allPeople.find(p => p.id === params.id || p.uuid === params.id) as Person | undefined
  
  // Show loading while fetching people data or if we don't have active company yet
  if (isLoading || !activeCompany || (!person && allPeople.length === 0)) {
    return (
      <div className="min-h-screen bg-background">
        <AppHeader />
        <main className="container mx-auto px-6 py-8">
          <div className="text-center">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <h2 className="text-lg font-semibold mb-2">Carregando...</h2>
            <p className="text-muted-foreground">Buscando informações da pessoa</p>
          </div>
        </main>
      </div>
    )
  }
  
  // Only show "not found" if we have loaded people data and still don't have the person
  if (!person && allPeople.length > 0) {
    return (
      <div className="min-h-screen bg-background">
        <AppHeader />
        <main className="container mx-auto px-6 py-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-2">Pessoa não encontrada</h1>
            <p className="text-muted-foreground mb-4">
              A pessoa que você está procurando não foi encontrada ou pode não existir.
            </p>
            <Button onClick={() => router.back()} variant="outline">
              <ArrowLeft className="h-4 w-4 mr-2" />
              Voltar
            </Button>
          </div>
        </main>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <AppHeader />
      
      <main className="container mx-auto px-6 py-8">
        {/* Header */}
        <div className="mb-8">
          <Button 
            onClick={() => router.back()} 
            variant="ghost" 
            className="mb-4"
          >
            <ArrowLeft className="h-4 w-4 mr-2" />
            Voltar
          </Button>
          
          <div className="flex items-start gap-6">
            <Avatar className="h-24 w-24">
              <AvatarImage src={person.avatar} alt={person.name} />
              <AvatarFallback className="bg-primary/10 text-primary text-xl">
                {getInitials(person.name)}
              </AvatarFallback>
            </Avatar>
            
            <div className="flex-1">
              <div className="flex items-center gap-3 mb-2">
                <h1 className="text-3xl font-bold">{person.name}</h1>
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => setShowEditModal(true)}
                  className="h-8 w-8 p-0"
                >
                  <Edit3 className="h-4 w-4" />
                </Button>
              </div>
              <p className="text-lg text-muted-foreground mb-4">{person.position || person.role || 'Cargo não informado'}</p>
              
              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                {person.department && (
                  <div className="flex items-center gap-1">
                    <MapPin className="h-4 w-4" />
                    {person.department}
                  </div>
                )}
                <div className="flex items-center gap-1">
                  <Calendar className="h-4 w-4" />
                  Na empresa há {person.startDate ? formatTimeAgoWithoutSuffix(person.startDate) : 'não informado'}
                </div>
                <div className="flex items-center gap-1">
                  <Clock className="h-4 w-4" />
                  Último 1:1: {getMockDaysAgo()} dias atrás
                </div>
              </div>
              
              <div className="flex gap-2 mt-4">
                {person.hasKids && (
                  <Badge variant="secondary" className="gap-1">
                    <User className="h-3 w-3" />
                    Tem filhos
                  </Badge>
                )}
                {person.email && (
                  <Badge variant="outline" className="gap-1">
                    📧 {person.email}
                  </Badge>
                )}
                {person.phone && (
                  <Badge variant="outline" className="gap-1">
                    📞 {person.phone}
                  </Badge>
                )}
              </div>
            </div>
          </div>
        </div>

        {/* Tabs */}
        <Tabs value={activeTab} onValueChange={handleTabChange}>
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="info">Informações</TabsTrigger>
            <TabsTrigger value="history">Histórico</TabsTrigger>
            <TabsTrigger value="feedbacks">Feedbacks</TabsTrigger>
            <TabsTrigger value="chat">Chat IA</TabsTrigger>
          </TabsList>
          
          {/* Informações Tab */}
          <TabsContent value="info" className="space-y-6">
            <PersonInfoTab person={person} allPeople={allPeople} />
          </TabsContent>
          
          {/* Histórico Tab */}
          <TabsContent value="history" className="space-y-6">
            <PersonHistoryTab 
              person={person} 
              allPeople={allPeople} 
              onCreatePerson={openCreateDialog}
            />
          </TabsContent>
          
          {/* Feedbacks Tab */}
          <TabsContent value="feedbacks" className="space-y-6">
            <PersonFeedbackTab person={person} allPeople={allPeople} />
          </TabsContent>
          
          {/* Chat IA Tab */}
          <TabsContent value="chat" className="space-y-6">
            <PersonChatTab person={person} />
          </TabsContent>
        </Tabs>
      </main>

      <CreatePersonDialog
        open={showCreatePersonDialog}
        onClose={closeCreateDialog}
        personName={personToCreate}
        newPersonName={newPersonName}
        newPersonRole={newPersonRole}
        setNewPersonName={setNewPersonName}
        setNewPersonRole={setNewPersonRole}
        onCreatePerson={() => handleCreatePerson(person.companyId)}
      />

      <PersonModal
        open={showEditModal}
        onClose={() => setShowEditModal(false)}
        mode="edit"
        person={person}
        onSubmit={handleEditPerson}
      />
    </div>
  )
}