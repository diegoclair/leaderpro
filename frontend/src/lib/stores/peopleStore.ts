import { create } from 'zustand'
import { Person, OneOnOneSession, Feedback, AISuggestion } from '../types'
import { 
  mockPeople, 
  mockOneOnOneSessions, 
  mockFeedbacks, 
  mockAISuggestions,
  getPersonsByCompany,
  getSessionsByPerson,
  getFeedbacksByPerson,
  getSuggestionsByPerson,
  getUpcomingOneOnOnes
} from '../data/mockData'

interface PeopleState {
  people: Person[]
  oneOnOneSessions: OneOnOneSession[]
  feedbacks: Feedback[]
  aiSuggestions: AISuggestion[]
  isLoading: boolean

  // Actions
  addPerson: (person: Person) => void
  updatePerson: (id: string, updates: Partial<Person>) => void
  deletePerson: (id: string) => void
  
  addOneOnOneSession: (session: OneOnOneSession) => void
  updateOneOnOneSession: (id: string, updates: Partial<OneOnOneSession>) => void
  
  addFeedback: (feedback: Feedback) => void
  updateFeedback: (id: string, updates: Partial<Feedback>) => void
  
  addAISuggestion: (suggestion: AISuggestion) => void
  markSuggestionAsUsed: (id: string) => void
  
  loadPeopleData: () => void

  // Selectors
  getPeopleByCompany: (companyId: string) => Person[]
  getPersonById: (id: string) => Person | undefined
  getSessionsByPerson: (personId: string) => OneOnOneSession[]
  getFeedbacksByPerson: (personId: string) => Feedback[]
  getSuggestionsByPerson: (personId: string) => AISuggestion[]
  getUpcomingOneOnOnes: (companyId: string) => Array<OneOnOneSession & { person: Person }>
}

export const usePeopleStore = create<PeopleState>((set, get) => ({
  people: [],
  oneOnOneSessions: [],
  feedbacks: [],
  aiSuggestions: [],
  isLoading: false,

  addPerson: (person: Person) => {
    set((state) => ({
      people: [...state.people, person]
    }))
  },

  updatePerson: (id: string, updates: Partial<Person>) => {
    set((state) => ({
      people: state.people.map((person) =>
        person.id === id ? { ...person, ...updates } : person
      )
    }))
  },

  deletePerson: (id: string) => {
    set((state) => ({
      people: state.people.filter(p => p.id !== id),
      oneOnOneSessions: state.oneOnOneSessions.filter(s => s.personId !== id),
      feedbacks: state.feedbacks.filter(f => f.personId !== id),
      aiSuggestions: state.aiSuggestions.filter(s => s.personId !== id)
    }))
  },

  addOneOnOneSession: (session: OneOnOneSession) => {
    set((state) => ({
      oneOnOneSessions: [...state.oneOnOneSessions, session]
    }))
  },

  updateOneOnOneSession: (id: string, updates: Partial<OneOnOneSession>) => {
    set((state) => ({
      oneOnOneSessions: state.oneOnOneSessions.map((session) =>
        session.id === id ? { ...session, ...updates } : session
      )
    }))
  },

  addFeedback: (feedback: Feedback) => {
    set((state) => ({
      feedbacks: [...state.feedbacks, feedback]
    }))
  },

  updateFeedback: (id: string, updates: Partial<Feedback>) => {
    set((state) => ({
      feedbacks: state.feedbacks.map((feedback) =>
        feedback.id === id ? { ...feedback, ...updates } : feedback
      )
    }))
  },

  addAISuggestion: (suggestion: AISuggestion) => {
    set((state) => ({
      aiSuggestions: [...state.aiSuggestions, suggestion]
    }))
  },

  markSuggestionAsUsed: (id: string) => {
    set((state) => ({
      aiSuggestions: state.aiSuggestions.map((suggestion) =>
        suggestion.id === id ? { ...suggestion, isUsed: true } : suggestion
      )
    }))
  },

  loadPeopleData: () => {
    set({ isLoading: true })
    
    // Simulate API call
    setTimeout(() => {
      set({
        people: mockPeople,
        oneOnOneSessions: mockOneOnOneSessions,
        feedbacks: mockFeedbacks,
        aiSuggestions: mockAISuggestions,
        isLoading: false
      })
    }, 150)
  },

  // Selectors
  getPeopleByCompany: (companyId: string) => {
    return get().people.filter(person => person.companyId === companyId)
  },

  getPersonById: (id: string) => {
    return get().people.find(person => person.id === id)
  },

  getSessionsByPerson: (personId: string) => {
    return get().oneOnOneSessions.filter(session => session.personId === personId)
  },

  getFeedbacksByPerson: (personId: string) => {
    return get().feedbacks.filter(feedback => feedback.personId === personId)
  },

  getSuggestionsByPerson: (personId: string) => {
    return get().aiSuggestions.filter(suggestion => suggestion.personId === personId && !suggestion.isUsed)
  },

  getUpcomingOneOnOnes: (companyId: string) => {
    const state = get()
    const companyPeople = state.people.filter(person => person.companyId === companyId)
    const upcomingMeetings: Array<OneOnOneSession & { person: Person }> = []
    
    companyPeople.forEach(person => {
      if (person.nextOneOnOne && person.nextOneOnOne > new Date()) {
        upcomingMeetings.push({
          id: `upcoming-${person.id}`,
          personId: person.id,
          date: person.nextOneOnOne,
          notes: '',
          aiSuggestions: state.aiSuggestions
            .filter(s => s.personId === person.id && !s.isUsed)
            .map(s => s.content),
          mentions: [],
          status: 'scheduled' as const,
          person
        })
      }
    })
    
    return upcomingMeetings.sort((a, b) => a.date.getTime() - b.date.getTime())
  }
}))

// Selectors for easier component usage - get all people and filter in component
export const useAllPeopleFromStore = () => usePeopleStore(state => state.people)

// Note: Hooks with parameters removed to avoid getSnapshot caching issues
// Use useAllPeopleFromStore and filter in components with useMemo instead

// Get all data separately to avoid object recreation
export const useAllAISuggestions = () => usePeopleStore(state => state.aiSuggestions)
export const useAllSessions = () => usePeopleStore(state => state.oneOnOneSessions)
export const useAllFeedbacks = () => usePeopleStore(state => state.feedbacks)

// Individual action hooks to avoid object recreation
export const useAddPerson = () => usePeopleStore(state => state.addPerson)
export const useUpdatePerson = () => usePeopleStore(state => state.updatePerson)
export const useDeletePerson = () => usePeopleStore(state => state.deletePerson)
export const useAddOneOnOneSession = () => usePeopleStore(state => state.addOneOnOneSession)
export const useUpdateOneOnOneSession = () => usePeopleStore(state => state.updateOneOnOneSession)
export const useAddFeedback = () => usePeopleStore(state => state.addFeedback)
export const useAddAISuggestion = () => usePeopleStore(state => state.addAISuggestion)
export const useMarkSuggestionAsUsed = () => usePeopleStore(state => state.markSuggestionAsUsed)
export const useLoadPeopleData = () => usePeopleStore(state => state.loadPeopleData)