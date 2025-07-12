import { create } from 'zustand'
import { persist } from 'zustand/middleware'
import { apiClient, setAuthStateGetter } from '../api/client'
import { useNotificationStore } from './notificationStore'

export interface User {
  uuid: string
  email: string
  name: string
  phone?: string
  profilePhoto?: string
  plan: string
  trialEndsAt?: string
  subscribedAt?: string
  timezone?: string
  language?: string
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
  emailVerified: boolean
}

interface AuthTokens {
  accessToken: string
  accessTokenExpiresAt: string
  refreshToken: string
  refreshTokenExpiresAt: string
}

interface AuthState {
  user: User | null
  tokens: AuthTokens | null
  isLoading: boolean
  isAuthenticated: boolean
}

interface AuthActions {
  login: (email: string, password: string) => Promise<void>
  register: (data: RegisterData) => Promise<void>
  logout: () => Promise<void>
  refreshToken: () => Promise<void>
  getProfile: () => Promise<void>
  setUser: (user: User) => void
  setTokens: (tokens: AuthTokens) => void
  clearAuth: () => void
}

export interface RegisterData {
  name: string
  email: string
  password: string
  phone?: string
  timezone?: string
  language?: string
}

type AuthStore = AuthState & AuthActions

export const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // State
      user: null,
      tokens: null,
      isLoading: false,
      isAuthenticated: false,

      // Actions
      login: async (email: string, password: string) => {
        set({ isLoading: true })
        
        try {
          const authResponse = await apiClient.post('/auth/login', { email, password })
          
          set({ 
            user: authResponse.user,
            tokens: {
              accessToken: authResponse.auth.access_token,
              accessTokenExpiresAt: authResponse.auth.access_token_expires_at,
              refreshToken: authResponse.auth.refresh_token,
              refreshTokenExpiresAt: authResponse.auth.refresh_token_expires_at,
            },
            isAuthenticated: true,
            isLoading: false 
          })
          
        } catch (error) {
          set({ isLoading: false })
          
          // Mostrar notificação de erro
          const errorMessage = error instanceof Error ? error.message : 'Erro ao fazer login'
          useNotificationStore.getState().showError('Erro no Login', errorMessage)
          
          throw error
        }
      },

      register: async (data: RegisterData) => {
        set({ isLoading: true })
        
        try {
          // Backend já faz login automático e retorna user + tokens
          const authResponse = await apiClient.post('/users', data)
          
          set({ 
            user: authResponse.user,
            tokens: {
              accessToken: authResponse.auth.access_token,
              accessTokenExpiresAt: authResponse.auth.access_token_expires_at,
              refreshToken: authResponse.auth.refresh_token,
              refreshTokenExpiresAt: authResponse.auth.refresh_token_expires_at,
            },
            isAuthenticated: true,
            isLoading: false 
          })
          
        } catch (error) {
          set({ isLoading: false })
          
          // Mostrar notificação de erro
          const errorMessage = error instanceof Error ? error.message : 'Erro ao criar conta'
          useNotificationStore.getState().showError('Erro no Registro', errorMessage)
          
          throw error
        }
      },

      logout: async () => {
        try {
          // Usar authPost que já inclui o token automaticamente
          await apiClient.authPost('/auth/logout')
        } catch (error) {
          console.error('Erro ao fazer logout no servidor:', error)
        } finally {
          get().clearAuth()
        }
      },

      refreshToken: async () => {
        const { tokens } = get()
        
        if (!tokens?.refreshToken) {
          get().clearAuth()
          throw new Error('Token de refresh não encontrado')
        }

        try {
          const newTokens = await apiClient.post('/auth/refresh-token', { 
            refresh_token: tokens.refreshToken 
          })
          
          set({ 
            tokens: {
              ...tokens,
              accessToken: newTokens.access_token,
              accessTokenExpiresAt: newTokens.access_token_expires_at
            }
          })
          
        } catch (error) {
          get().clearAuth()
          throw error
        }
      },

      getProfile: async () => {
        try {
          // authGet já gerencia token automaticamente e renovação em caso de 401
          const user: User = await apiClient.authGet('/users/profile')
          set({ user })
          
        } catch (error) {
          console.error('Erro ao buscar perfil:', error)
          throw error
        }
      },

      setUser: (user: User) => {
        set({ user })
      },

      setTokens: (tokens: AuthTokens) => {
        set({ tokens, isAuthenticated: true })
      },

      clearAuth: () => {
        set({ 
          user: null, 
          tokens: null, 
          isAuthenticated: false,
          isLoading: false 
        })
      }
    }),
    {
      name: 'auth-storage',
      partialize: (state) => ({
        user: state.user,
        tokens: state.tokens,
        isAuthenticated: state.isAuthenticated
      })
    }
  )
)

// Configurar o getter de autenticação para o API client
setAuthStateGetter(() => {
  const { tokens, refreshToken, clearAuth } = useAuthStore.getState()
  return { tokens, refreshToken, clearAuth }
})

// Re-exportar o apiClient para facilitar o uso
export { apiClient }