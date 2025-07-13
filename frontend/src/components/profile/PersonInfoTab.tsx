'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { User, MessageSquare, Calendar, MessageCircle } from 'lucide-react'
import { Person } from '@/lib/types'
import { useMentions } from '@/hooks/useMentions'
import { useCreatePerson } from '@/hooks/useCreatePerson'
import { getInitials } from '@/lib/utils/names'
import MentionSuggestions from './MentionSuggestions'
import CreatePersonDialog from './CreatePersonDialog'

interface PersonInfoTabProps {
  person: Person
  allPeople: Person[]
}

export function PersonInfoTab({ person, allPeople }: PersonInfoTabProps) {
  const [newNote, setNewNote] = React.useState('')
  const [recordType, setRecordType] = React.useState<'oneOnOne' | 'feedback' | 'note'>('note')

  const {
    showMentionSuggestions,
    handleTextChange,
    getFilteredPeople,
    insertMention,
    setShowMentionSuggestions,
    detectMentions
  } = useMentions({
    allPeople,
    currentPersonCompanyId: person.companyId
  })

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

  const handleAddNote = () => {
    if (newNote.trim()) {
      // TODO: Add note to person's history
      // Process @mentions here
      const mentions = detectMentions(newNote)
      console.log('Adding record:', {
        type: recordType,
        content: newNote,
        mentions: mentions,
        personId: person.id
      })
      setNewNote('')
      setRecordType('note') // Reset to default
    }
  }

  const handleNoteChange = (value: string) => {
    setNewNote(value)
    handleTextChange(value, openCreateDialog)
  }

  const handleSelectMention = (selectedPerson: Person) => {
    const newText = insertMention(newNote, selectedPerson)
    setNewNote(newText)
    setShowMentionSuggestions(false)
  }

  const filteredPeople = getFilteredPeople(person.id)

  return (
    <div className="grid md:grid-cols-2 gap-6">
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="h-5 w-5" />
            Informações Pessoais
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label className="text-sm font-medium text-muted-foreground">Data de Início</Label>
            <p className="text-sm">{person.startDate ? new Date(person.startDate).toLocaleDateString('pt-BR') : 'Não informado'}</p>
          </div>
          <div>
            <Label className="text-sm font-medium text-muted-foreground">Departamento/Squad</Label>
            <p className="text-sm">{person.department || 'Não informado'}</p>
          </div>
          <div>
            <Label className="text-sm font-medium text-muted-foreground">Contato</Label>
            <div className="space-y-1">
              {person.email && (
                <p className="text-sm">📧 {person.email}</p>
              )}
              {person.phone && (
                <p className="text-sm">📞 {person.phone}</p>
              )}
              {!person.email && !person.phone && (
                <p className="text-sm text-muted-foreground">Não informado</p>
              )}
            </div>
          </div>
          {person.interests && (
            <div>
              <Label className="text-sm font-medium text-muted-foreground">Interesses</Label>
              <p className="text-sm">{person.interests}</p>
            </div>
          )}
          {person.personality && (
            <div>
              <Label className="text-sm font-medium text-muted-foreground">Personalidade</Label>
              <p className="text-sm">{person.personality}</p>
            </div>
          )}
          {person.notes && (
            <div>
              <Label className="text-sm font-medium text-muted-foreground">Notas</Label>
              <p className="text-sm text-muted-foreground">{person.notes}</p>
            </div>
          )}
        </CardContent>
      </Card>
      
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <MessageSquare className="h-5 w-5" />
            Registrar Informação
          </CardTitle>
        </CardHeader>
        <CardContent className="space-y-4">
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

          <div className="space-y-2 relative">
            <Label htmlFor="note">
              {recordType === 'oneOnOne' ? 'Resumo da reunião 1:1' :
               recordType === 'feedback' ? 'Feedback para a pessoa' :
               'Anotação sobre a pessoa'}
            </Label>
            <Textarea 
              id="note"
              placeholder={
                recordType === 'oneOnOne' 
                  ? "Registre os pontos principais da reunião 1:1. Use @nome para mencionar outras pessoas da equipe..." 
                  : recordType === 'feedback'
                  ? "Registre um feedback positivo ou construtivo. Use @nome para mencionar outras pessoas da equipe..."
                  : "Registre uma observação ou anotação geral. Use @nome para mencionar outras pessoas da equipe..."
              }
              value={newNote}
              onChange={(e) => handleNoteChange(e.target.value)}
              className="min-h-[100px]"
            />
            
            <MentionSuggestions
              show={showMentionSuggestions}
              people={filteredPeople}
              onSelect={handleSelectMention}
            />
          </div>
          <Button onClick={handleAddNote} className="w-full">
            {recordType === 'oneOnOne' ? 'Registrar 1:1' :
             recordType === 'feedback' ? 'Registrar Feedback' :
             'Registrar Anotação'}
          </Button>
        </CardContent>
      </Card>

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