import { create } from 'zustand'
import { persist } from 'zustand/middleware'

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

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000'

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
          const response = await fetch(`${API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ email, password })
          })

          if (!response.ok) {
            const errorData = await response.json()
            throw new Error(errorData.message || 'Erro ao fazer login')
          }

          const authResponse = await response.json()
          
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
          throw error
        }
      },

      register: async (data: RegisterData) => {
        set({ isLoading: true })
        
        try {
          // Criar usuário - agora retorna AuthResponse direto
          const registerResponse = await fetch(`${API_BASE_URL}/users`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify(data)
          })

          if (!registerResponse.ok) {
            const errorData = await registerResponse.json()
            throw new Error(errorData.message || 'Erro ao criar conta')
          }

          // Backend já faz login automático e retorna user + tokens
          const authResponse = await registerResponse.json()
          
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
          throw error
        }
      },

      logout: async () => {
        const { tokens } = get()
        
        try {
          if (tokens?.accessToken) {
            await fetch(`${API_BASE_URL}/auth/logout`, {
              method: 'POST',
              headers: {
                'Authorization': `Bearer ${tokens.accessToken}`,
                'Content-Type': 'application/json',
              }
            })
          }
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
          const response = await fetch(`${API_BASE_URL}/auth/refresh-token`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ refresh_token: tokens.refreshToken })
          })

          if (!response.ok) {
            get().clearAuth()
            throw new Error('Token expirado, faça login novamente')
          }

          const newTokens = await response.json()
          
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
        const { tokens } = get()
        
        if (!tokens?.accessToken) {
          throw new Error('Token de acesso não encontrado')
        }

        try {
          const response = await fetch(`${API_BASE_URL}/users/profile`, {
            headers: {
              'Authorization': `Bearer ${tokens.accessToken}`,
              'Content-Type': 'application/json',
            }
          })

          if (!response.ok) {
            if (response.status === 401) {
              // Token expirado, tentar renovar
              await get().refreshToken()
              return get().getProfile()
            }
            throw new Error('Erro ao buscar perfil')
          }

          const user: User = await response.json()
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

// Interceptor para renovação automática de token
export const authFetch = async (url: string, options: RequestInit = {}) => {
  const { tokens, refreshToken, clearAuth } = useAuthStore.getState()
  
  if (!tokens?.accessToken) {
    throw new Error('Usuário não autenticado')
  }

  // Adicionar token de autorização
  const headers = {
    ...options.headers,
    'Authorization': `Bearer ${tokens.accessToken}`,
    'Content-Type': 'application/json',
  }

  let response = await fetch(url, {
    ...options,
    headers
  })

  // Se token expirou, tentar renovar
  if (response.status === 401) {
    try {
      await refreshToken()
      const newTokens = useAuthStore.getState().tokens
      
      // Repetir requisição com novo token
      response = await fetch(url, {
        ...options,
        headers: {
          ...options.headers,
          'Authorization': `Bearer ${newTokens?.accessToken}`,
          'Content-Type': 'application/json',
        }
      })
      
    } catch (error) {
      clearAuth()
      throw new Error('Sessão expirada, faça login novamente')
    }
  }

  return response
}