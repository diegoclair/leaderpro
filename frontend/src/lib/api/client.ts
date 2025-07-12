// Constantes centralizadas de autenticação
const AUTH_CONSTANTS = {
  TOKEN_HEADER: 'user-token', // Header que o backend espera
  API_BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000',
} as const

// Import notification store (lazy para evitar ciclo de dependência)
let getNotificationStore: () => any

const setNotificationStore = () => {
  if (!getNotificationStore) {
    getNotificationStore = () => {
      // Import dinâmico para evitar problemas de circular dependency
      const { useNotificationStore } = require('../stores/notificationStore')
      return useNotificationStore.getState()
    }
  }
}

interface ApiResponse<T = any> {
  data?: T
  error?: string
  message?: string
}

interface AuthTokens {
  accessToken: string
  accessTokenExpiresAt: string
  refreshToken: string
  refreshTokenExpiresAt: string
}

interface AuthState {
  tokens: AuthTokens | null
  clearAuth: () => void
  refreshToken: () => Promise<void>
}

// Função para obter o estado de autenticação
// Evita problemas de importação circular
let getAuthState: () => AuthState

export const setAuthStateGetter = (getter: () => AuthState) => {
  getAuthState = getter
}

class ApiClient {
  private baseURL: string

  constructor(baseURL: string = AUTH_CONSTANTS.API_BASE_URL) {
    this.baseURL = baseURL
  }

  // Método privado para fazer requisições base
  private async baseRequest(url: string, options: RequestInit = {}): Promise<Response> {
    const fullUrl = url.startsWith('http') ? url : `${this.baseURL}${url}`
    
    const defaultHeaders = {
      'Content-Type': 'application/json',
    }

    return fetch(fullUrl, {
      ...options,
      headers: {
        ...defaultHeaders,
        ...options.headers,
      },
    })
  }

  // Método público para requisições sem autenticação
  async request<T = any>(url: string, options: RequestInit = {}): Promise<T> {
    const response = await this.baseRequest(url, options)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Erro na requisição' }))
      const errorMessage = errorData.message || `Erro ${response.status}`
      
      // Mostrar notificação de erro (exceto para login/register que já tratam)
      setNotificationStore()
      if (getNotificationStore && !url.includes('/auth/login') && !url.includes('/users')) {
        getNotificationStore().showError('Erro na API', errorMessage)
      }
      
      throw new Error(errorMessage)
    }

    // Verificar se há conteúdo para fazer parse JSON
    const contentType = response.headers.get('content-type')
    const hasJsonContent = contentType && contentType.includes('application/json')
    
    if (hasJsonContent) {
      return response.json()
    } else {
      // Para respostas 204 No Content ou outras sem JSON, retornar objeto vazio
      return {}
    }
  }

  // Método para requisições autenticadas com renovação automática
  async authenticatedRequest<T = any>(url: string, options: RequestInit = {}): Promise<T> {
    if (!getAuthState) {
      throw new Error('AuthState getter não foi configurado')
    }

    const { tokens, refreshToken: refreshAuthToken, clearAuth } = getAuthState()
    
    if (!tokens?.accessToken) {
      throw new Error('Usuário não autenticado')
    }

    // Adicionar token de autenticação
    const authHeaders = {
      [AUTH_CONSTANTS.TOKEN_HEADER]: tokens.accessToken,
    }

    let response = await this.baseRequest(url, {
      ...options,
      headers: {
        ...options.headers,
        ...authHeaders,
      },
    })

    // Se token expirou (401), tentar renovar
    if (response.status === 401) {
      try {
        await refreshAuthToken()
        const newTokens = getAuthState().tokens
        
        if (!newTokens?.accessToken) {
          throw new Error('Falha ao renovar token')
        }

        // Repetir requisição com novo token
        response = await this.baseRequest(url, {
          ...options,
          headers: {
            ...options.headers,
            [AUTH_CONSTANTS.TOKEN_HEADER]: newTokens.accessToken,
          },
        })
        
      } catch (error) {
        clearAuth()
        throw new Error('Sessão expirada, faça login novamente')
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Erro na requisição' }))
      const errorMessage = errorData.message || `Erro ${response.status}`
      
      // Mostrar notificação de erro para requisições autenticadas
      setNotificationStore()
      if (getNotificationStore && !url.includes('/auth/logout')) {
        getNotificationStore().showError('Erro na requisição', errorMessage)
      }
      
      throw new Error(errorMessage)
    }

    // Verificar se há conteúdo para fazer parse JSON
    const contentType = response.headers.get('content-type')
    const hasJsonContent = contentType && contentType.includes('application/json')
    
    if (hasJsonContent) {
      return response.json()
    } else {
      // Para respostas 204 No Content ou outras sem JSON, retornar objeto vazio
      return {}
    }
  }

  // Métodos de conveniência para diferentes tipos de requisição
  async get<T = any>(url: string, options?: RequestInit): Promise<T> {
    return this.request<T>(url, { ...options, method: 'GET' })
  }

  async post<T = any>(url: string, data?: any, options?: RequestInit): Promise<T> {
    return this.request<T>(url, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T = any>(url: string, data?: any, options?: RequestInit): Promise<T> {
    return this.request<T>(url, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T = any>(url: string, options?: RequestInit): Promise<T> {
    return this.request<T>(url, { ...options, method: 'DELETE' })
  }

  // Métodos autenticados
  async authGet<T = any>(url: string, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, { ...options, method: 'GET' })
  }

  async authPost<T = any>(url: string, data?: any, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async authPut<T = any>(url: string, data?: any, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async authDelete<T = any>(url: string, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, { ...options, method: 'DELETE' })
  }
}

// Instância singleton do cliente API
export const apiClient = new ApiClient()

// Exportar constantes para uso em outros lugares se necessário
export { AUTH_CONSTANTS }