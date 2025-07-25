'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { User, MessageSquare, Calendar, MessageCircle } from 'lucide-react'
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
    <>
      {/* Main Content - 3 Columns */}
      <div className="grid lg:grid-cols-12 gap-6">
        {/* Informações Pessoais */}
        <Card className="lg:col-span-3">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-sm">
              <User className="h-4 w-4" />
              Informações Pessoais
            </CardTitle>
          </CardHeader>
          <CardContent className="space-y-3">
              {person.startDate && (
                <div>
                  <Label className="text-xs font-medium text-muted-foreground">Data de Início</Label>
                  <p className="text-sm">{new Date(person.startDate).toLocaleDateString('pt-BR')}</p>
                </div>
              )}
              
              {person.department && (
                <div>
                  <Label className="text-xs font-medium text-muted-foreground">Departamento</Label>
                  <p className="text-sm">{person.department}</p>
                </div>
              )}
              
              {(person.email || person.phone) && (
                <div>
                  <Label className="text-xs font-medium text-muted-foreground">Contato</Label>
                  <div className="space-y-1">
                    {person.email && (
                      <p className="text-sm">📧 {person.email}</p>
                    )}
                    {person.phone && (
                      <p className="text-sm">📞 {person.phone}</p>
                    )}
                  </div>
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
                <div>
                  <Label className="text-xs font-medium text-muted-foreground">Notas</Label>
                  <p className="text-sm text-muted-foreground">{person.notes}</p>
                </div>
              )}
          </CardContent>
        </Card>
        
        {/* Formulário de Registro */}
        <Card className="lg:col-span-5">
          <CardHeader className="pb-3">
            <CardTitle className="flex items-center gap-2 text-base">
              <MessageSquare className="h-5 w-5" />
              Registrar Nova Informação
            </CardTitle>
          </CardHeader>
        <CardContent className="space-y-4">
          {/* Tipo de Registro */}
          <div className="space-y-2">
            <Label htmlFor="record-type">Tipo de registro</Label>
            <Select value={recordType} onValueChange={(value: 'oneOnOne' | 'feedback' | 'note') => setRecordType(value)}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione o tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="oneOnOne">
                  <div className="flex items-center gap-2">
                    <Calendar className="h-4 w-4" />
                    <span>1:1 - Reunião individual</span>
                  </div>
                </SelectItem>
                <SelectItem value="feedback">
                  <div className="flex items-center gap-2">
                    <MessageCircle className="h-4 w-4" />
                    <span>Feedback - Avaliação/opinião</span>
                  </div>
                </SelectItem>
                <SelectItem value="note">
                  <div className="flex items-center gap-2">
                    <MessageSquare className="h-4 w-4" />
                    <span>Anotação - Observação geral</span>
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Campos condicionais para feedback - mais compactos */}
          {recordType === 'feedback' && (
            <div className="grid grid-cols-2 gap-3">
              <div className="space-y-2">
                <Label htmlFor="feedback-type">Tipo</Label>
                <Select value={feedbackType} onValueChange={(value: 'positive' | 'constructive' | 'neutral') => setFeedbackType(value)}>
                  <SelectTrigger>
                    <SelectValue placeholder="Tipo" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="positive">
                      <div className="flex items-center gap-2">
                        <span className="text-green-600">👍</span>
                        <span>Positivo</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="constructive">
                      <div className="flex items-center gap-2">
                        <span className="text-orange-600">🔨</span>
                        <span>Construtivo</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="neutral">
                      <div className="flex items-center gap-2">
                        <span className="text-blue-600">💭</span>
                        <span>Neutro</span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <Label htmlFor="feedback-category">Categoria</Label>
                <Select value={feedbackCategory} onValueChange={(value: 'performance' | 'behavior' | 'skill' | 'collaboration') => setFeedbackCategory(value)}>
                  <SelectTrigger>
                    <SelectValue placeholder="Categoria" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="performance">
                      <div className="flex items-center gap-2">
                        <span>📊</span>
                        <span>Performance</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="behavior">
                      <div className="flex items-center gap-2">
                        <span>🤝</span>
                        <span>Comportamento</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="skill">
                      <div className="flex items-center gap-2">
                        <span>🎯</span>
                        <span>Habilidade</span>
                      </div>
                    </SelectItem>
                    <SelectItem value="collaboration">
                      <div className="flex items-center gap-2">
                        <span>👥</span>
                        <span>Colaboração</span>
                      </div>
                    </SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
          )}

          {/* Campo de Texto com Mentions */}
          <div className="space-y-2">
            <Label htmlFor="note">
              {recordType === 'oneOnOne' ? 'Resumo da reunião 1:1' :
               recordType === 'feedback' ? 'Feedback para a pessoa' :
               'Anotação sobre a pessoa'}
            </Label>
            
            <MentionsTextarea
              value={newNote}
              onChange={handleNoteChange}
              people={filteredPeople}
              placeholder={
                recordType === 'oneOnOne' 
                  ? "Registre os pontos principais da reunião 1:1. Use @nome para mencionar outras pessoas da equipe..." 
                  : recordType === 'feedback'
                  ? "Registre um feedback positivo ou construtivo. Use @nome para mencionar outras pessoas da equipe..."
                  : "Registre uma observação ou anotação geral. Use @nome para mencionar outras pessoas da equipe..."
              }
              disabled={isSubmitting}
              minHeight={100}
            />
            
            <p className="text-xs text-muted-foreground">
              💡 Digite @ para mencionar outras pessoas da equipe
            </p>
          </div>

          {/* Botão de Envio */}
          <Button 
            onClick={handleAddNote} 
            className="w-full" 
            disabled={isSubmitting || !newNote.trim()}
          >
            {isSubmitting ? 'Salvando...' :
             recordType === 'oneOnOne' ? 'Registrar 1:1' :
             recordType === 'feedback' ? 'Registrar Feedback' :
             'Registrar Anotação'}
          </Button>
          </CardContent>
        </Card>
        
        {/* AI Assistant Sidebar */}
        <AIAssistantSidebar 
          person={person} 
          className="lg:col-span-4"
        />
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
    </>
  )
}