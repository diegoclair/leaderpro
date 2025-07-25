import { create } from 'zustand'
import { Company } from '../types'
import { apiClient } from './authStore'
import { storageManager } from '../utils/storageManager'
import type { ApiCompany } from '../types/api'

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
  loadCompanies: () => Promise<void>
}

export const useCompanyStore = create<CompanyState>((set, get) => ({
  companies: [],
  activeCompany: null,
  isLoading: false,

  setActiveCompany: (company: Company) => {
    set({ activeCompany: company })
    // Persist active company using storage manager
    storageManager.set('leaderpro-active-company', company.id)
  },

  addCompany: (company: Company) => {
    set((state) => {
      const newCompanies = [...state.companies, company]
      
      // Se é a primeira empresa ou é uma empresa default e não há empresa ativa,
      // defini-la como ativa automaticamente
      let newActiveCompany = state.activeCompany
      
      if (!state.activeCompany || (company.isDefault && state.companies.length === 0)) {
        newActiveCompany = company
        
        // Persist active company using storage manager
        storageManager.set('leaderpro-active-company', company.id)
      }
      
      return {
        companies: newCompanies,
        activeCompany: newActiveCompany
      }
    })
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

  loadCompanies: async () => {
    const { isLoading } = get()
    
    // Evitar múltiplas chamadas simultâneas
    if (isLoading) return
    
    set({ isLoading: true })
    
    try {
      // Buscar empresas da API
      const companiesFromAPI = await apiClient.authGet<ApiCompany[]>('/companies')
      
      // Converter para o formato esperado pelo frontend
      const companies: Company[] = companiesFromAPI.map((company) => ({
        id: company.uuid,
        uuid: company.uuid,
        name: company.name,
        industry: company.industry || '',
        size: company.size || '',
        role: company.role || '',
        isDefault: company.is_default || false,
        createdAt: new Date(company.created_at),
        updatedAt: new Date(company.created_at)
      }))

      // Ordenar: default primeiro, depois por nome
      companies.sort((a, b) => {
        if (a.isDefault && !b.isDefault) return -1
        if (!a.isDefault && b.isDefault) return 1
        return a.name.localeCompare(b.name)
      })
      
      // Definir empresa ativa
      let activeCompany = companies.find(c => c.isDefault) || companies[0] || null
      
      // Tentar restaurar empresa previamente selecionada APENAS se pertencer ao usuário atual
      const savedCompanyId = storageManager.get<string>('leaderpro-active-company')
      if (savedCompanyId) {
        const savedCompany = companies.find(c => c.id === savedCompanyId)
        if (savedCompany) {
          activeCompany = savedCompany
        } else {
          // Se a empresa salva não pertence ao usuário atual, limpar storage
          storageManager.remove('leaderpro-active-company')
        }
      }
      
      
      set({
        companies,
        activeCompany,
        isLoading: false
      })
      
    } catch (error) {
      console.error('❌ Erro ao carregar empresas:', error)
      
      // Em caso de erro, deixar vazio para acionar onboarding
      set({
        companies: [],
        activeCompany: null,
        isLoading: false
      })
    }
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