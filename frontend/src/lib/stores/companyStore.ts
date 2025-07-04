import { create } from 'zustand'
import { Company } from '../types'
import { mockCompanies } from '../data/mockData'

interface CompanyState {
  companies: Company[]
  activeCompany: Company | null
  isLoading: boolean
  
  // Actions
  setActiveCompany: (company: Company) => void
  addCompany: (company: Company) => void
  updateCompany: (id: string, updates: Partial<Company>) => void
  deleteCompany: (id: string) => void
  setDefaultCompany: (id: string) => void
  loadCompanies: () => void
}

export const useCompanyStore = create<CompanyState>((set, get) => ({
  companies: [],
  activeCompany: null,
  isLoading: false,

  setActiveCompany: (company: Company) => {
    set({ activeCompany: company })
    // Persist active company to localStorage
    if (typeof window !== 'undefined') {
      localStorage.setItem('leaderpro-active-company', company.id)
    }
  },

  addCompany: (company: Company) => {
    set((state) => ({
      companies: [...state.companies, company]
    }))
  },

  updateCompany: (id: string, updates: Partial<Company>) => {
    set((state) => ({
      companies: state.companies.map((company) =>
        company.id === id 
          ? { ...company, ...updates, updatedAt: new Date() }
          : company
      ),
      activeCompany: state.activeCompany?.id === id 
        ? { ...state.activeCompany, ...updates, updatedAt: new Date() }
        : state.activeCompany
    }))
  },

  deleteCompany: (id: string) => {
    set((state) => {
      const remainingCompanies = state.companies.filter(c => c.id !== id)
      const newActiveCompany = state.activeCompany?.id === id 
        ? remainingCompanies.find(c => c.isDefault) || remainingCompanies[0] || null
        : state.activeCompany

      return {
        companies: remainingCompanies,
        activeCompany: newActiveCompany
      }
    })
  },

  setDefaultCompany: (id: string) => {
    set((state) => ({
      companies: state.companies.map((company) => ({
        ...company,
        isDefault: company.id === id,
        updatedAt: company.id === id ? new Date() : company.updatedAt
      }))
    }))
  },

  loadCompanies: () => {
    set({ isLoading: true })
    
    // Simulate API call - in real app this would be an async API call
    setTimeout(() => {
      // Try to restore previously selected company from localStorage
      let activeCompany = mockCompanies.find(c => c.isDefault) || mockCompanies[0]
      
      if (typeof window !== 'undefined') {
        const savedCompanyId = localStorage.getItem('leaderpro-active-company')
        if (savedCompanyId) {
          const savedCompany = mockCompanies.find(c => c.id === savedCompanyId)
          if (savedCompany) {
            activeCompany = savedCompany
          }
        }
      }
      
      set({
        companies: mockCompanies,
        activeCompany: activeCompany,
        isLoading: false
      })
    }, 100)
  }
}))

// Selectors for easier component usage
export const useActiveCompany = () => useCompanyStore(state => state.activeCompany)
export const useCompanies = () => useCompanyStore(state => state.companies)

// Individual action hooks to avoid object recreation
export const useSetActiveCompany = () => useCompanyStore(state => state.setActiveCompany)
export const useAddCompany = () => useCompanyStore(state => state.addCompany)
export const useUpdateCompany = () => useCompanyStore(state => state.updateCompany)
export const useDeleteCompany = () => useCompanyStore(state => state.deleteCompany)
export const useSetDefaultCompany = () => useCompanyStore(state => state.setDefaultCompany)
export const useLoadCompanies = () => useCompanyStore(state => state.loadCompanies)