'use client'

import React from 'react'
import { InfoCard } from '@/components/ui/info-card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { User, MessageSquare } from 'lucide-react'
import { Person } from '@/lib/types'
import { useCreatePerson } from '@/hooks/useCreatePerson'
import CreatePersonDialog from './CreatePersonDialog'
import { useCompanyStore } from '@/lib/stores/companyStore'
import { apiClient } from '@/lib/stores/authStore'
import { MentionsTextarea } from '@/components/ui/mentions-textarea'
import { AIAssistantSidebar } from './AIAssistantSidebar'

interface PersonInfoTabProps {
  person: Person
  allPeople: Person[]
}

export function PersonInfoTab({ person, allPeople }: PersonInfoTabProps) {
  const [newNote, setNewNote] = React.useState('')
  const [recordType, setRecordType] = React.useState<'oneOnOne' | 'feedback' | 'note'>('note')
  const [feedbackType, setFeedbackType] = React.useState<'positive' | 'constructive' | 'neutral'>('positive')
  const [feedbackCategory, setFeedbackCategory] = React.useState<'performance' | 'behavior' | 'skill' | 'collaboration'>('performance')
  const [isSubmitting, setIsSubmitting] = React.useState(false)

  const { activeCompany } = useCompanyStore()

  // Filter people from the same company (excluding current person)
  const filteredPeople = allPeople.filter(p => 
    p.id !== person.id && 
    p.companyId === person.companyId
  )

  const {
    showCreatePersonDialog,
    personToCreate,
    newPersonName,
    newPersonRole,
    setNewPersonName,
    setNewPersonRole,
    closeCreateDialog,
    handleCreatePerson
  } = useCreatePerson()


  const handleAddNote = async () => {
    if (newNote.trim() && !isSubmitting && activeCompany) {
      setIsSubmitting(true)
      
      try {
        // Convert recordType to backend format
        const noteType = recordType === 'oneOnOne' ? 'one_on_one' : 
                        recordType === 'feedback' ? 'feedback' : 
                        'observation'

        // Prepare note data - the content already contains the mentions in the correct format
        const noteData = {
          type: noteType,
          content: newNote, // This already contains {{person:uuid|name}} tokens
          feedback_type: recordType === 'feedback' ? feedbackType : undefined,
          feedback_category: recordType === 'feedback' ? feedbackCategory : undefined
        }

        // Send to backend
        await apiClient.authPost(`/companies/${activeCompany.id}/people/${person.id}/notes`, noteData)

        // Reset form on success
        setNewNote('')
        setRecordType('note')
        setFeedbackType('positive')
        setFeedbackCategory('performance')
        
      } catch (error) {
        console.error('❌ Error creating note:', error)
        // TODO: Add proper error handling/toast notification
      } finally {
        setIsSubmitting(false)
      }
    }
  }

  const handleNoteChange = (value: string) => {
    setNewNote(value)
  }


  return (
    <div className="flex gap-6">
      {/* Sidebar Esquerda - Profile Style */}
      <div className="w-72 flex-shrink-0 hidden lg:block">
        <div className="sticky top-6 space-y-4">
          {/* Info Card Compacta */}
          <InfoCard
            title="Informações"
            icon={User}
            contentClassName="p-4 space-y-3"
          >
            {person.department && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Departamento</Label>
                <p className="text-sm font-medium">{person.department}</p>
              </div>
            )}
            
            {person.startDate && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Na empresa</Label>
                <p className="text-sm">{new Date(person.startDate).toLocaleDateString('pt-BR')}</p>
              </div>
            )}
            
            {person.email && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Email</Label>
                <p className="text-sm break-all">{person.email}</p>
              </div>
            )}
            
            {person.phone && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Telefone</Label>
                <p className="text-sm">{person.phone}</p>
              </div>
            )}
            
            {person.interests && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Interesses</Label>
                <p className="text-sm">{person.interests}</p>
              </div>
            )}
            
            {person.personality && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Personalidade</Label>
                <p className="text-sm">{person.personality}</p>
              </div>
            )}
            
            {person.notes && (
              <div className="pt-3 border-t">
                <Label className="text-xs font-medium text-muted-foreground">Notas</Label>
                <p className="text-sm text-muted-foreground mt-1">{person.notes}</p>
              </div>
            )}
          </InfoCard>
        </div>
      </div>

      {/* Main Content - 2 Colunas */}
      <div className="flex-1 min-w-0">
        <div className="grid lg:grid-cols-2 gap-6">
        {/* Formulário de Registro - Compacto */}
        <InfoCard
          title="Registrar Nova Informação"
          icon={MessageSquare}
          contentClassName="space-y-3 p-4"
        >
          <div className="grid gap-3">
            {/* Tipo de Registro + Campos condicionais em linha */}
            <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
              <div>
                <Label className="text-sm">Tipo</Label>
                <Select value={recordType} onValueChange={(value: 'oneOnOne' | 'feedback' | 'note') => setRecordType(value)}>
                  <SelectTrigger className="h-9">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="oneOnOne">1:1</SelectItem>
                    <SelectItem value="feedback">Feedback</SelectItem>
                    <SelectItem value="note">Anotação</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              {recordType === 'feedback' && (
                <>
                  <div>
                    <Label className="text-sm">Tom</Label>
                    <Select value={feedbackType} onValueChange={(value: 'positive' | 'constructive' | 'neutral') => setFeedbackType(value)}>
                      <SelectTrigger className="h-9">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="positive">👍 Positivo</SelectItem>
                        <SelectItem value="constructive">🔨 Construtivo</SelectItem>
                        <SelectItem value="neutral">💭 Neutro</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>

                  <div>
                    <Label className="text-sm">Categoria</Label>
                    <Select value={feedbackCategory} onValueChange={(value: 'performance' | 'behavior' | 'skill' | 'collaboration') => setFeedbackCategory(value)}>
                      <SelectTrigger className="h-9">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="performance">📊 Performance</SelectItem>
                        <SelectItem value="behavior">🤝 Comportamento</SelectItem>
                        <SelectItem value="skill">🎯 Habilidade</SelectItem>
                        <SelectItem value="collaboration">👥 Colaboração</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </>
              )}
            </div>

            {/* Campo de Texto - Mais compacto */}
            <div>
              <Label className="text-sm">
                {recordType === 'oneOnOne' ? 'Resumo da reunião 1:1' :
                 recordType === 'feedback' ? 'Feedback' : 'Observação'}
              </Label>
              
              <MentionsTextarea
                value={newNote}
                onChange={handleNoteChange}
                people={filteredPeople}
                placeholder={
                  recordType === 'oneOnOne' 
                    ? "Pontos principais da reunião. Use @nome para mencionar pessoas..." 
                    : recordType === 'feedback'
                    ? "Descreva o feedback. Use @nome para mencionar pessoas..."
                    : "Descreva a observação. Use @nome para mencionar pessoas..."
                }
                disabled={isSubmitting}
                minHeight={100}
                className="mt-1"
              />
            </div>

            {/* Botão + Dica */}
            <div className="flex items-center justify-between">
              <p className="text-xs text-muted-foreground">
                💡 Use @ para mencionar pessoas
              </p>
              <Button 
                onClick={handleAddNote} 
                disabled={isSubmitting || !newNote.trim()}
                size="sm"
                className="px-4"
              >
                {isSubmitting ? 'Salvando...' : 'Registrar'}
              </Button>
            </div>
          </div>
        </InfoCard>

          {/* AI Assistant - Altura otimizada */}
          <AIAssistantSidebar 
            person={person}
            className="h-[500px]"
          />
        </div>
        
        {/* Mobile: Info Card quando sidebar está hidden */}
        <div className="lg:hidden mb-6">
          <InfoCard
            title="Informações"
            icon={User}
            contentClassName="p-4 grid grid-cols-2 gap-4"
          >
            {person.department && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Departamento</Label>
                <p className="text-sm font-medium">{person.department}</p>
              </div>
            )}
            
            {person.email && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Email</Label>
                <p className="text-sm break-all">{person.email}</p>
              </div>
            )}
            
            {person.startDate && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Na empresa</Label>
                <p className="text-sm">{new Date(person.startDate).toLocaleDateString('pt-BR')}</p>
              </div>
            )}
            
            {person.phone && (
              <div>
                <Label className="text-xs font-medium text-muted-foreground">Telefone</Label>
                <p className="text-sm">{person.phone}</p>
              </div>
            )}
          </InfoCard>
        </div>
      </div>
      
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