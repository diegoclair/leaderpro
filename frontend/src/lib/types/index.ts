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

export type Address = {
  id: string
  uuid: string
  personId: string
  city?: string
  state?: string
  country?: string
  isPrimary: boolean
  createdAt: Date
  updatedAt: Date
  active: boolean
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
  gender?: 'male' | 'female' | 'other'
  interests?: string
  personality?: string
  lastOneOnOneDate?: Date
  createdAt: Date
  age?: number
  tenure?: number
  // Address information
  primaryAddress?: Address
  // Legacy fields for compatibility
  role?: string
  personalInfo?: {
    hasChildren?: boolean
    hasPets?: boolean
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

// Re-export API types for convenience
export * from './api'