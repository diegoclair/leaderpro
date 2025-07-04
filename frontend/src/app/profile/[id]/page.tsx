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
  Clock
} from 'lucide-react'
import { useAllPeopleFromStore } from '@/lib/stores/peopleStore'
import { Person } from '@/lib/types'
import { PersonInfoTab } from '@/components/profile/PersonInfoTab'
import { PersonHistoryTab } from '@/components/profile/PersonHistoryTab'
import { PersonFeedbackTab } from '@/components/profile/PersonFeedbackTab'
import { PersonChatTab } from '@/components/profile/PersonChatTab'
import { useCreatePerson } from '@/hooks/useCreatePerson'
import CreatePersonDialog from '@/components/profile/CreatePersonDialog'
import { formatTimeAgoWithoutSuffix, getMockDaysAgo } from '@/lib/utils/dates'
import { getInitials } from '@/lib/utils/names'

export default function ProfilePage() {
  const params = useParams()
  const router = useRouter()
  const allPeople = useAllPeopleFromStore()
  
  // Get tab from URL or default to 'info' (using window for static export compatibility)
  const [activeTab, setActiveTab] = useState('info')
  
  useEffect(() => {
    if (typeof window !== 'undefined') {
      const urlParams = new URLSearchParams(window.location.search)
      const tab = urlParams.get('tab')
      if (tab) {
        setActiveTab(tab)
      }
    }
  }, [])

  // Update URL when tab changes
  const handleTabChange = (tab: string) => {
    setActiveTab(tab)
    const newUrl = `/profile/${params.id}?tab=${tab}`
    router.replace(newUrl, { scroll: false })
  }
  
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
  
  const person = allPeople.find(p => p.id === params.id) as Person | undefined
  
  if (!person) {
    return (
      <div className="min-h-screen bg-background">
        <AppHeader />
        <main className="container mx-auto px-6 py-8">
          <div className="text-center">
            <h1 className="text-2xl font-bold mb-2">Pessoa não encontrada</h1>
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
              <h1 className="text-3xl font-bold mb-2">{person.name}</h1>
              <p className="text-lg text-muted-foreground mb-4">{person.role}</p>
              
              <div className="flex flex-wrap gap-4 text-sm text-muted-foreground">
                {person.personalInfo.location && (
                  <div className="flex items-center gap-1">
                    <MapPin className="h-4 w-4" />
                    {person.personalInfo.location}
                  </div>
                )}
                <div className="flex items-center gap-1">
                  <Calendar className="h-4 w-4" />
                  Na empresa há {formatTimeAgoWithoutSuffix(person.startDate)}
                </div>
                <div className="flex items-center gap-1">
                  <Clock className="h-4 w-4" />
                  Último 1:1: {getMockDaysAgo()} dias atrás
                </div>
              </div>
              
              <div className="flex gap-2 mt-4">
                {person.personalInfo.hasChildren && (
                  <Badge variant="secondary" className="gap-1">
                    <User className="h-3 w-3" />
                    Tem filhos
                  </Badge>
                )}
                {person.personalInfo.hasPets && (
                  <Badge variant="secondary" className="gap-1">
                    <Dog className="h-3 w-3" />
                    Tem pets
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
    </div>
  )
}