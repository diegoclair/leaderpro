// Core Types for LeaderPro

export type Company = {
  id: string
  uuid: string
  name: string
  industry: string
  size: string
  role: string
  isDefault: boolean
  createdAt: Date
  updatedAt: Date
}

export type Person = {
  id: string
  uuid: string
  companyId: string
  name: string
  email?: string
  position?: string
  department?: string
  phone?: string
  birthday?: Date
  startDate?: Date
  isManager: boolean
  managerUUID?: string
  notes?: string
  hasKids: boolean
  interests?: string
  personality?: string
  createdAt: Date
  age?: number
  tenure?: number
  // Legacy fields for compatibility
  role?: string
  personalInfo?: {
    hasChildren?: boolean
    hasPets?: boolean
    location?: string
    interests?: string[]
    hobbies?: string[]
    pets?: string[]
    personalNotes?: string
  }
  nextOneOnOne?: Date
  avatar?: string
}

export type OneOnOneSession = {
  id: string
  personId: string
  date: Date
  notes: string
  aiSuggestions: string[]
  mentions: Mention[]
  status: 'scheduled' | 'completed' | 'cancelled'
}

export type Mention = {
  id: string
  sessionId: string
  mentionedPersonId: string
  mentionedPersonName: string
  context: string
  createdAt: Date
}

export type Feedback = {
  id: string
  personId: string
  content: string
  type: 'positive' | 'negative' | 'neutral'
  source: 'direct' | 'mention' // direct = líder escreveu, mention = veio de @mention
  sourcePersonId?: string // se source = mention, quem mencionou
  sessionId?: string // se veio de uma sessão 1:1
  createdAt: Date
}

export type AISuggestion = {
  id: string
  personId: string
  type: 'question' | 'reminder' | 'insight'
  content: string
  context: string
  priority: 'low' | 'medium' | 'high'
  createdAt: Date
  isUsed: boolean
}