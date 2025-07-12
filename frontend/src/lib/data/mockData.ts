import { Company, Person, OneOnOneSession, Feedback, AISuggestion } from '../types'

// Mock Companies
export const mockCompanies: Company[] = [
  {
    id: '1',
    uuid: '1',
    name: 'TechCorp',
    industry: 'technology',
    size: 'medium',
    role: 'CTO',
    isDefault: false, // Não será mais default por padrão
    createdAt: new Date('2024-01-15'),
    updatedAt: new Date('2024-07-01')
  }
]

// Mock People
export const mockPeople: Person[] = [
  // TechCorp
  {
    id: '1',
    companyId: '1',
    name: 'Maria Santos',
    role: 'Analista de Sistemas',
    email: 'maria.santos@techcorp.com',
    startDate: new Date('2022-03-15'),
    personalInfo: {
      hasChildren: false,
      hasPets: true,
      location: 'São Paulo, SP',
      hobbies: ['Leitura', 'Yoga', 'Culinária'],
      pets: ['Gato - Mimi', 'Cachorro - Thor'],
      personalNotes: 'Gosta de viajar e conhecer novas culturas'
    },
    nextOneOnOne: new Date('2024-07-04T14:00:00'),
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Maria'
  },
  {
    id: '2',
    companyId: '1', 
    name: 'João Silva',
    role: 'Coordenador de Projetos',
    email: 'joao.silva@techcorp.com',
    startDate: new Date('2023-07-01'),
    personalInfo: {
      hasChildren: true,
      hasPets: false,
      location: 'Rio de Janeiro, RJ',
      hobbies: ['Futebol', 'Videogame', 'Culinária'],
      personalNotes: 'Pai de dois filhos, gosta de equilibrar trabalho e família'
    },
    nextOneOnOne: new Date('2024-07-05T10:00:00'),
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Joao'
  },
  
  // StartupXYZ
  {
    id: '3',
    companyId: '2',
    name: 'Pedro Costa',
    role: 'Gerente de Projetos',
    email: 'pedro.costa@startupxyz.com',
    startDate: new Date('2023-01-10'),
    personalInfo: {
      hasChildren: false,
      hasPets: true,
      location: 'Belo Horizonte, MG',
      hobbies: ['Ciclismo', 'Fotografia', 'Viagem'],
      pets: ['Cachorro - Buddy'],
      personalNotes: 'Apaixonado por natureza e esportes ao ar livre'
    },
    nextOneOnOne: new Date('2024-07-08T15:30:00'),
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Pedro'
  },
  {
    id: '4',
    companyId: '2',
    name: 'Ana Lima',
    role: 'Analista Junior', 
    email: 'ana.lima@startupxyz.com',
    startDate: new Date('2024-02-01'),
    personalInfo: {
      hasChildren: false,
      hasPets: false,
      location: 'Porto Alegre, RS',
      hobbies: ['Design', 'Pintura', 'Música'],
      personalNotes: 'Muito criativa e sempre busca aprender coisas novas'
    },
    nextOneOnOne: new Date('2024-07-06T09:00:00'),
    avatar: 'https://api.dicebear.com/7.x/avataaars/svg?seed=Ana'
  }
]

// Mock One-on-One Sessions
export const mockOneOnOneSessions: OneOnOneSession[] = [
  {
    id: '1',
    personId: '1', // Maria
    date: new Date('2024-06-20T14:00:00'),
    notes: 'Maria demonstrou muito interesse em liderar o próximo projeto de migração. Mencionou que @João tem ajudado muito com DevOps e que gostaria de trabalhar mais próxima dele. Sugeriu também que @Pedro poderia dar insights sobre arquitetura de microservices.',
    aiSuggestions: [
      'Pergunte sobre os desafios técnicos do projeto atual',
      'Como está o relacionamento com o time?',
      'Que tipo de projeto gostaria de liderar?'
    ],
    mentions: [
      {
        id: '1',
        sessionId: '1',
        mentionedPersonId: '2',
        mentionedPersonName: 'João',
        context: 'tem ajudado muito com DevOps',
        createdAt: new Date('2024-06-20T14:15:00')
      },
      {
        id: '2', 
        sessionId: '1',
        mentionedPersonId: '3',
        mentionedPersonName: 'Pedro',
        context: 'poderia dar insights sobre arquitetura de microservices',
        createdAt: new Date('2024-06-20T14:20:00')
      }
    ],
    status: 'completed'
  },
  {
    id: '2',
    personId: '2', // João
    date: new Date('2024-06-18T10:00:00'), 
    notes: 'João mencionou que está gostando muito de trabalhar com DevOps. Tem ajudado @Maria com deploy e automação. Falou que as crianças estão em casa devido às férias escolares, mas está conseguindo equilibrar bem trabalho e família.',
    aiSuggestions: [
      'Como está o equilíbrio trabalho-família com as crianças em casa?',
      'Que aspectos de DevOps mais te interessam?',
      'Como você vê sua evolução técnica?'
    ],
    mentions: [
      {
        id: '3',
        sessionId: '2', 
        mentionedPersonId: '1',
        mentionedPersonName: 'Maria',
        context: 'tem ajudado com deploy e automação',
        createdAt: new Date('2024-06-18T10:10:00')
      }
    ],
    status: 'completed'
  }
]

// Mock Feedbacks 
export const mockFeedbacks: Feedback[] = [
  // Feedbacks diretos
  {
    id: '1',
    personId: '1', // Maria
    content: 'Excelente liderança no projeto de refatoração. Mostrou iniciativa e organizou o time muito bem.',
    type: 'positive',
    source: 'direct',
    sessionId: '1',
    createdAt: new Date('2024-06-20T14:30:00')
  },
  {
    id: '2',
    personId: '2', // João  
    content: 'Muito proativo em DevOps. Está sempre disponível para ajudar o time com deployments.',
    type: 'positive',
    source: 'direct', 
    sessionId: '2',
    createdAt: new Date('2024-06-18T10:30:00')
  },
  
  // Feedbacks vindos de @mentions
  {
    id: '3',
    personId: '2', // João
    content: 'tem ajudado muito com DevOps',
    type: 'positive',
    source: 'mention',
    sourcePersonId: '1', // Maria mencionou
    sessionId: '1',
    createdAt: new Date('2024-06-20T14:15:00')
  },
  {
    id: '4',
    personId: '1', // Maria
    content: 'tem ajudado com deploy e automação', 
    type: 'positive',
    source: 'mention',
    sourcePersonId: '2', // João mencionou
    sessionId: '2',
    createdAt: new Date('2024-06-18T10:10:00')
  }
]

// Mock AI Suggestions
export const mockAISuggestions: AISuggestion[] = [
  {
    id: '1',
    personId: '1', // Maria
    type: 'question',
    content: 'Como você se sente sobre liderar um time maior no próximo projeto?',
    context: 'Maria mostrou interesse em liderança na última reunião',
    priority: 'high',
    createdAt: new Date('2024-07-01T09:00:00'),
    isUsed: false
  },
  {
    id: '2', 
    personId: '2', // João
    type: 'reminder',
    content: 'Pergunte sobre as férias escolares - João mencionou que as crianças estão em casa',
    context: 'Julho é período de férias escolares e João tem filhos',
    priority: 'medium',
    createdAt: new Date('2024-07-01T09:05:00'),
    isUsed: false
  },
  {
    id: '3',
    personId: '2', // João
    type: 'insight', 
    content: 'João recebeu feedback positivo de Maria sobre DevOps - considere discutir progressão nessa área',
    context: 'Feedback cruzado positivo sobre habilidades técnicas',
    priority: 'high',
    createdAt: new Date('2024-06-21T08:00:00'),
    isUsed: false
  },
  {
    id: '4',
    personId: '3', // Pedro
    type: 'question',
    content: 'Como está a adaptação da Ana no time? Ela precisa de mais mentoria?',
    context: 'Ana é junior e Pedro é tech lead, importante acompanhar mentoria',
    priority: 'medium', 
    createdAt: new Date('2024-07-02T14:00:00'),
    isUsed: false
  }
]

// Helper functions
export const getPersonsByCompany = (companyId: string): Person[] => {
  return mockPeople.filter(person => person.companyId === companyId)
}

export const getSessionsByPerson = (personId: string): OneOnOneSession[] => {
  return mockOneOnOneSessions.filter(session => session.personId === personId)
}

export const getFeedbacksByPerson = (personId: string): Feedback[] => {
  return mockFeedbacks.filter(feedback => feedback.personId === personId)
}

export const getSuggestionsByPerson = (personId: string): AISuggestion[] => {
  return mockAISuggestions.filter(suggestion => suggestion.personId === personId)
}

export const getSessionsWhereMentioned = (personId: string): Array<OneOnOneSession & { mentionContext: string }> => {
  const sessionsWithMentions: Array<OneOnOneSession & { mentionContext: string }> = []
  
  mockOneOnOneSessions.forEach(session => {
    session.mentions.forEach(mention => {
      if (mention.mentionedPersonId === personId) {
        sessionsWithMentions.push({
          ...session,
          mentionContext: mention.context
        })
      }
    })
  })
  
  return sessionsWithMentions.sort((a, b) => b.date.getTime() - a.date.getTime())
}

export const getUpcomingOneOnOnes = (companyId: string): Array<OneOnOneSession & { person: Person }> => {
  const companyPeople = getPersonsByCompany(companyId)
  const upcomingMeetings: Array<OneOnOneSession & { person: Person }> = []
  
  companyPeople.forEach(person => {
    if (person.nextOneOnOne && person.nextOneOnOne > new Date()) {
      upcomingMeetings.push({
        id: `upcoming-${person.id}`,
        personId: person.id,
        date: person.nextOneOnOne,
        notes: '',
        aiSuggestions: getSuggestionsByPerson(person.id).map(s => s.content),
        mentions: [],
        status: 'scheduled' as const,
        person
      })
    }
  })
  
  return upcomingMeetings.sort((a, b) => a.date.getTime() - b.date.getTime())
}