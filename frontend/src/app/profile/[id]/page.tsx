'use client'

import React, { useState } from 'react'
import { useParams, useRouter } from 'next/navigation'
import { AppHeader } from '@/components/layout/AppHeader'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Textarea } from '@/components/ui/textarea'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { 
  ArrowLeft, 
  MapPin, 
  Calendar, 
  User, 
  MessageSquare, 
  Heart,
  Dog,
  Home,
  MessageCircle,
  Clock,
  Send
} from 'lucide-react'
import { useAllPeopleFromStore } from '@/lib/stores/peopleStore'
import { Person } from '@/lib/types'
import { formatDistanceToNow, format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { mockFeedbacks, getFeedbacksByPerson, getSessionsByPerson } from '@/lib/data/mockData'

export default function ProfilePage() {
  const params = useParams()
  const router = useRouter()
  const allPeople = useAllPeopleFromStore()
  const [activeTab, setActiveTab] = useState('info')
  const [newNote, setNewNote] = useState('')
  const [chatMessage, setChatMessage] = useState('')
  const [showMentionSuggestions, setShowMentionSuggestions] = useState(false)
  const [mentionQuery, setMentionQuery] = useState('')
  const [recordType, setRecordType] = useState<'oneOnOne' | 'feedback' | 'note'>('note')
  const [showCreatePersonDialog, setShowCreatePersonDialog] = useState(false)
  const [personToCreate, setPersonToCreate] = useState('')
  const [newPersonName, setNewPersonName] = useState('')
  const [newPersonRole, setNewPersonRole] = useState('')
  
  const person = allPeople.find(p => p.id === params.id) as Person | undefined
  const personFeedbacks = person ? getFeedbacksByPerson(person.id) : []
  const personSessions = person ? getSessionsByPerson(person.id) : []
  
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

  const getInitials = (name: string) => {
    return name
      .split(' ')
      .map(n => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2)
  }

  const handleAddNote = () => {
    if (newNote.trim()) {
      // TODO: Add note to person's history
      // Process @mentions here
      const mentions = detectMentions(newNote)
      console.log('Adding record:', {
        type: recordType,
        content: newNote,
        mentions: mentions,
        personId: person?.id
      })
      setNewNote('')
      setRecordType('note') // Reset to default
    }
  }

  const detectMentions = (text: string) => {
    const mentionRegex = /@(\w+)/g
    const mentions = []
    let match
    
    while ((match = mentionRegex.exec(text)) !== null) {
      const mentionedName = match[1]
      const mentionedPerson = allPeople.find(p => 
        p.name.toLowerCase().includes(mentionedName.toLowerCase()) &&
        p.companyId === person?.companyId // Only same company
      )
      if (mentionedPerson) {
        mentions.push({
          name: mentionedName,
          person: mentionedPerson,
          context: text
        })
      }
    }
    
    return mentions
  }

  const handleNoteChange = (value: string) => {
    setNewNote(value)
    
    // Detect @ mentions
    const lastAtIndex = value.lastIndexOf('@')
    if (lastAtIndex !== -1) {
      const textAfterAt = value.substring(lastAtIndex + 1)
      const spaceIndex = textAfterAt.indexOf(' ')
      const query = spaceIndex === -1 ? textAfterAt : textAfterAt.substring(0, spaceIndex)
      
      // Check if user just typed space after @name (mention completion)
      if (spaceIndex !== -1 && query.length > 0) {
        // User typed @name + space, check if person exists
        const mentionedPerson = allPeople.find(p => 
          p.name.toLowerCase().includes(query.toLowerCase()) &&
          p.companyId === person?.companyId
        )
        
        if (!mentionedPerson) {
          // Person doesn't exist, offer to create
          setPersonToCreate(query)
          setNewPersonName(query)
          setShowCreatePersonDialog(true)
        }
        
        // Hide suggestions after space
        setShowMentionSuggestions(false)
      } else {
        // Show suggestions immediately when @ is typed, even with empty query
        setMentionQuery(query)
        setShowMentionSuggestions(true)
      }
    } else {
      setShowMentionSuggestions(false)
    }
  }

  const handleSelectMention = (selectedPerson: Person) => {
    const lastAtIndex = newNote.lastIndexOf('@')
    const beforeAt = newNote.substring(0, lastAtIndex)
    const afterAt = newNote.substring(lastAtIndex + 1)
    const spaceIndex = afterAt.indexOf(' ')
    const afterMention = spaceIndex === -1 ? '' : afterAt.substring(spaceIndex)
    
    const newText = `${beforeAt}@${selectedPerson.name}${afterMention}`
    setNewNote(newText)
    setShowMentionSuggestions(false)
  }

  const filteredPeople = allPeople.filter(p => 
    p.id !== person?.id && // Don't suggest the current person
    p.companyId === person?.companyId && // Only same company
    (mentionQuery === '' || p.name.toLowerCase().includes(mentionQuery.toLowerCase())) // Show all if empty query
  )

  const handleSendMessage = () => {
    if (chatMessage.trim()) {
      // TODO: Send message to AI chat
      console.log('Sending message to AI:', chatMessage)
      setChatMessage('')
    }
  }

  const handleCreatePerson = () => {
    if (newPersonName.trim() && newPersonRole.trim()) {
      // TODO: Add person to store/backend
      console.log('Creating person:', {
        name: newPersonName,
        role: newPersonRole,
        companyId: person?.companyId
      })
      setShowCreatePersonDialog(false)
      setPersonToCreate('')
      setNewPersonName('')
      setNewPersonRole('')
    }
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
                  Na empresa há {formatDistanceToNow(person.startDate, { 
                    locale: ptBR,
                    addSuffix: false 
                  })}
                </div>
                <div className="flex items-center gap-1">
                  <Clock className="h-4 w-4" />
                  Último 1:1: {Math.floor(Math.random() * 30) + 1} dias atrás
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
        <Tabs value={activeTab} onValueChange={setActiveTab}>
          <TabsList className="grid w-full grid-cols-4">
            <TabsTrigger value="info">Informações</TabsTrigger>
            <TabsTrigger value="history">Histórico</TabsTrigger>
            <TabsTrigger value="feedbacks">Feedbacks</TabsTrigger>
            <TabsTrigger value="chat">Chat IA</TabsTrigger>
          </TabsList>
          
          {/* Informações Tab */}
          <TabsContent value="info" className="space-y-6">
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
                    <Label className="text-sm font-medium text-muted-foreground">Email</Label>
                    <p className="text-sm">{person.email}</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Localização</Label>
                    <p className="text-sm">{person.personalInfo.location || 'Não informado'}</p>
                  </div>
                  <div>
                    <Label className="text-sm font-medium text-muted-foreground">Situação Familiar</Label>
                    <div className="flex gap-2 mt-1">
                      {person.personalInfo.hasChildren && (
                        <Badge variant="outline">Tem filhos</Badge>
                      )}
                      {person.personalInfo.hasPets && (
                        <Badge variant="outline">Tem pets</Badge>
                      )}
                      {!person.personalInfo.hasChildren && !person.personalInfo.hasPets && (
                        <p className="text-sm text-muted-foreground">Não informado</p>
                      )}
                    </div>
                  </div>
                  {person.personalInfo.pets && person.personalInfo.pets.length > 0 && (
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Pets</Label>
                      <div className="flex flex-wrap gap-1 mt-1">
                        {person.personalInfo.pets.map((pet, index) => (
                          <Badge key={index} variant="secondary" className="text-xs">
                            {pet}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                  {person.personalInfo.hobbies && person.personalInfo.hobbies.length > 0 && (
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Hobbies</Label>
                      <div className="flex flex-wrap gap-1 mt-1">
                        {person.personalInfo.hobbies.map((hobby, index) => (
                          <Badge key={index} variant="secondary" className="text-xs">
                            {hobby}
                          </Badge>
                        ))}
                      </div>
                    </div>
                  )}
                  {person.personalInfo.personalNotes && (
                    <div>
                      <Label className="text-sm font-medium text-muted-foreground">Notas Pessoais</Label>
                      <p className="text-sm text-muted-foreground">{person.personalInfo.personalNotes}</p>
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
                    
                    {/* Mention Suggestions */}
                    {showMentionSuggestions && filteredPeople.length > 0 && (
                      <div className="absolute z-10 w-full mt-1 bg-background border rounded-md shadow-lg max-h-32 overflow-y-auto">
                        {filteredPeople.slice(0, 5).map((suggestedPerson) => (
                          <button
                            key={suggestedPerson.id}
                            onClick={() => handleSelectMention(suggestedPerson)}
                            className="w-full px-3 py-2 text-left hover:bg-muted flex items-center gap-2"
                          >
                            <Avatar className="h-6 w-6">
                              <AvatarImage src={suggestedPerson.avatar} alt={suggestedPerson.name} />
                              <AvatarFallback className="text-xs">
                                {getInitials(suggestedPerson.name)}
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <p className="text-sm font-medium">{suggestedPerson.name}</p>
                              <p className="text-xs text-muted-foreground">{suggestedPerson.role}</p>
                            </div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>
                  <Button onClick={handleAddNote} className="w-full">
                    {recordType === 'oneOnOne' ? 'Registrar 1:1' :
                     recordType === 'feedback' ? 'Registrar Feedback' :
                     'Registrar Anotação'}
                  </Button>
                </CardContent>
              </Card>
            </div>
          </TabsContent>
          
          {/* Histórico Tab */}
          <TabsContent value="history" className="space-y-6">
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
                              {formatDistanceToNow(session.date, { 
                                locale: ptBR,
                                addSuffix: true 
                              })}
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
                                    p.companyId === person?.companyId  // Only same company
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
                                          setPersonToCreate(mention.mentionedPersonName)
                                          setNewPersonName(mention.mentionedPersonName)
                                          setShowCreatePersonDialog(true)
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
                              {format(session.date, 'dd/MM/yyyy', { locale: ptBR })}
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
          </TabsContent>
          
          {/* Feedbacks Tab */}
          <TabsContent value="feedbacks" className="space-y-6">
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
                                  {sourcePerson ? sourcePerson.name.split(' ').map(n => n[0]).join('').slice(0, 2) : 'LI'}
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
                              {formatDistanceToNow(feedback.createdAt, { 
                                locale: ptBR,
                                addSuffix: true 
                              })}
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
          </TabsContent>
          
          {/* Chat IA Tab */}
          <TabsContent value="chat" className="space-y-6">
            <Card>
              <CardHeader>
                <CardTitle>Chat com IA - Insights sobre {person.name}</CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {/* Mock AI conversation */}
                  <div className="space-y-3">
                    <div className="bg-blue-50 dark:bg-blue-950/30 rounded-lg p-3">
                      <p className="text-sm">
                        <strong>IA:</strong> Com base no histórico, {person.name} tem demonstrado interesse em crescimento profissional. 
                        Algumas sugestões para a próxima 1:1:
                      </p>
                      <ul className="text-sm mt-2 space-y-1 ml-4">
                        <li>• Discutir oportunidades de liderança</li>
                        <li>• Explorar áreas de interesse específicas</li>
                        <li>• Definir metas de desenvolvimento</li>
                      </ul>
                    </div>
                  </div>
                  
                  <div className="space-y-2">
                    <Label htmlFor="chat-message">Pergunte algo sobre {person.name}</Label>
                    <div className="flex gap-2">
                      <Input 
                        id="chat-message"
                        placeholder="Como posso ajudar esta pessoa a crescer profissionalmente?"
                        value={chatMessage}
                        onChange={(e) => setChatMessage(e.target.value)}
                        onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
                      />
                      <Button onClick={handleSendMessage} size="sm">
                        <Send className="h-4 w-4" />
                      </Button>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </TabsContent>
        </Tabs>
      </main>

      {/* Create Person Dialog */}
      <Dialog open={showCreatePersonDialog} onOpenChange={setShowCreatePersonDialog}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Pessoa não encontrada</DialogTitle>
            <DialogDescription>
              A pessoa "@{personToCreate}" foi mencionada mas não existe no sistema. 
              Deseja adicioná-la agora?
            </DialogDescription>
          </DialogHeader>
          
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="person-name">Nome completo</Label>
              <Input
                id="person-name"
                value={newPersonName}
                onChange={(e) => setNewPersonName(e.target.value)}
                placeholder="Digite o nome completo"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="person-role">Cargo/Função</Label>
              <Input
                id="person-role"
                value={newPersonRole}
                onChange={(e) => setNewPersonRole(e.target.value)}
                placeholder="Ex: Analista, Coordenador, Gerente..."
              />
            </div>
          </div>

          <DialogFooter>
            <Button variant="outline" onClick={() => setShowCreatePersonDialog(false)}>
              Cancelar
            </Button>
            <Button onClick={handleCreatePerson}>
              Adicionar Pessoa
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  )
}