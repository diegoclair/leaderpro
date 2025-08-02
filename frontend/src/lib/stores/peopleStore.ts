import { create } from 'zustand'
import { Person, OneOnOneSession, Feedback, AISuggestion } from '../types'
import { apiClient } from '../api/client'
import type { DashboardResponse, ApiPerson } from '../types/api'

// Dashboard stats interface matching backend
interface DashboardStats {
  totalPeople: number
  oneOnOnesCountThisMonth: number
  feedbacksCountThisMonth: number
  averageDaysBetweenOneOnOnes: number
  lastMeetingDate?: Date
}

interface PeopleState {
  people: Person[]
  oneOnOneSessions: OneOnOneSession[]
  feedbacks: Feedback[]
  aiSuggestions: AISuggestion[]
  dashboardStats: DashboardStats | null
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
  loadPeopleFromAPI: (companyUuid: string) => Promise<void>
  loadDashboardData: (companyUuid: string) => Promise<void>

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
  dashboardStats: null,
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
    
    // No longer using mock data - this should not be called
    console.warn('loadPeopleData called - should use loadPeopleFromAPI instead')
    set({
      people: [],
      oneOnOneSessions: [],
      feedbacks: [],
      aiSuggestions: [],
      isLoading: false
    })
  },

  loadPeopleFromAPI: async (companyUuid: string) => {
    set({ isLoading: true })
    
    try {
      // Fetch people from API
      const response = await apiClient.authGet<ApiPerson[]>(`/companies/${companyUuid}/people`)
      
      // Convert API response to frontend format
      const apiPeople = Array.isArray(response) ? response : []
      const people: Person[] = apiPeople.map((apiPerson) => ({
        id: apiPerson.uuid || '',
        uuid: apiPerson.uuid || '',
        companyId: companyUuid, // Use the company UUID as companyId
        name: apiPerson.name || '',
        email: apiPerson.email,
        position: apiPerson.position,
        department: apiPerson.department,
        phone: apiPerson.phone,
        birthday: apiPerson.birthday && apiPerson.birthday !== '' ? new Date(apiPerson.birthday) : undefined,
        startDate: apiPerson.start_date && apiPerson.start_date !== '' ? new Date(apiPerson.start_date) : undefined,
        isManager: apiPerson.is_manager || false,
        managerUUID: apiPerson.manager_uuid,
        notes: apiPerson.notes,
        hasKids: apiPerson.has_kids || false,
        gender: apiPerson.gender,
        interests: apiPerson.interests,
        personality: apiPerson.personality,
        lastOneOnOneDate: apiPerson.last_one_on_one_date && apiPerson.last_one_on_one_date !== '' ? 
          new Date(apiPerson.last_one_on_one_date) : undefined,
        createdAt: apiPerson.created_at ? new Date(apiPerson.created_at) : new Date(),
        age: apiPerson.age,
        tenure: apiPerson.tenure,
        // Legacy compatibility
        role: apiPerson.position,
        personalInfo: {
          hasChildren: apiPerson.has_kids,
          interests: apiPerson.interests ? [apiPerson.interests] : [],
          personalNotes: apiPerson.notes,
          location: apiPerson.department // Use department as location for display
        }
      }))

      set({
        people,
        // TODO: Implement API endpoints for these features
        oneOnOneSessions: [],
        feedbacks: [],
        aiSuggestions: [],
        isLoading: false
      })
    } catch (error) {
      console.error('Error loading people from API:', error)
      // Fallback to empty data on error
      set({
        people: [],
        oneOnOneSessions: [],
        feedbacks: [],
        aiSuggestions: [],
        isLoading: false
      })
    }
  },

  loadDashboardData: async (companyUuid: string) => {
    set({ isLoading: true })
    
    try {
      // Fetch dashboard data (people + stats) from unified API
      const response = await apiClient.authGet<DashboardResponse>(`/companies/${companyUuid}/dashboard`)
      
      // Convert people API response to frontend format
      const apiPeople = Array.isArray(response.people) ? response.people : []
      const people: Person[] = apiPeople.map((apiPerson) => ({
        id: apiPerson.uuid || '',
        uuid: apiPerson.uuid || '',
        companyId: companyUuid,
        name: apiPerson.name || '',
        email: apiPerson.email,
        position: apiPerson.position,
        department: apiPerson.department,
        phone: apiPerson.phone,
        birthday: apiPerson.birthday && apiPerson.birthday !== '' ? new Date(apiPerson.birthday) : undefined,
        startDate: apiPerson.start_date && apiPerson.start_date !== '' ? new Date(apiPerson.start_date) : undefined,
        isManager: apiPerson.is_manager || false,
        managerUUID: apiPerson.manager_uuid,
        notes: apiPerson.notes,
        hasKids: apiPerson.has_kids || false,
        gender: apiPerson.gender,
        interests: apiPerson.interests,
        personality: apiPerson.personality,
        lastOneOnOneDate: apiPerson.last_one_on_one_date && apiPerson.last_one_on_one_date !== '' ? 
          new Date(apiPerson.last_one_on_one_date) : undefined,
        createdAt: apiPerson.created_at ? new Date(apiPerson.created_at) : new Date(),
        age: apiPerson.age,
        tenure: apiPerson.tenure,
        // Legacy compatibility
        role: apiPerson.position,
        personalInfo: {
          hasChildren: apiPerson.has_kids,
          interests: apiPerson.interests ? [apiPerson.interests] : [],
          personalNotes: apiPerson.notes,
          location: apiPerson.department
        }
      }))

      // Extract stats from response
      const stats: DashboardStats = {
        totalPeople: response.stats?.total_people || 0,
        oneOnOnesCountThisMonth: response.stats?.one_on_ones_this_month || 0,
        feedbacksCountThisMonth: response.stats?.feedbacks_count_this_month || 0,
        averageDaysBetweenOneOnOnes: response.stats?.average_frequency_days || 0,
        lastMeetingDate: response.stats?.last_meeting_date && response.stats.last_meeting_date !== '' ? 
          new Date(response.stats.last_meeting_date) : undefined
      }

      set({
        people,
        dashboardStats: stats,
        // TODO: Implement API endpoints for these features
        oneOnOneSessions: [],
        feedbacks: [],
        aiSuggestions: [],
        isLoading: false
      })
    } catch (error) {
      console.error('Error loading dashboard data from API:', error)
      // Fallback to empty data on error
      set({
        people: [],
        dashboardStats: null,
        oneOnOneSessions: [],
        feedbacks: [],
        aiSuggestions: [],
        isLoading: false
      })
    }
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
export const useLoadPeopleFromAPI = () => usePeopleStore(state => state.loadPeopleFromAPI)
export const useLoadDashboardData = () => usePeopleStore(state => state.loadDashboardData)
export const useDashboardStats = () => usePeopleStore(state => state.dashboardStats)