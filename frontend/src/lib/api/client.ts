// Constantes centralizadas de autenticação
const AUTH_CONSTANTS = {
  TOKEN_HEADER: 'user-token', // Header que o backend espera
  API_BASE_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:5000',
} as const

// Type for notification store to avoid circular dependency
type NotificationStoreState = {
  showError: (title: string, message?: string, duration?: number) => string
  showInfo: (title: string, message?: string, duration?: number) => string
}

// Import notification store (lazy para evitar ciclo de dependência)
let getNotificationStore: (() => Promise<NotificationStoreState>) | undefined

const setNotificationStore = async () => {
  if (!getNotificationStore) {
    getNotificationStore = async () => {
      // Import dinâmico para evitar problemas de circular dependency
      const { useNotificationStore } = await import('../stores/notificationStore')
      return useNotificationStore.getState()
    }
  }
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

    try {
      return await fetch(fullUrl, {
        ...options,
        headers: {
          ...defaultHeaders,
          ...options.headers,
        },
      })
    } catch (networkError) {
      console.error('❌ Erro de rede:', { fullUrl, error: networkError })
      throw new Error(`Erro de conexão: ${networkError instanceof Error ? networkError.message : 'Verifique se o backend está rodando'}`)
    }
  }

  // Método público para requisições sem autenticação
  async request<T>(url: string, options: RequestInit = {}): Promise<T> {
    const response = await this.baseRequest(url, options)
    
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Erro na requisição' }))
      const errorMessage = errorData.message || `Erro ${response.status}`
      
      // Mostrar notificação de erro (exceto para login/register que já tratam)
      await setNotificationStore()
      if (getNotificationStore && !url.includes('/auth/login') && !url.includes('/users')) {
        try {
          const store = await getNotificationStore()
          store.showError('Erro na API', errorMessage)
        } catch {
          // Se falhar ao mostrar notificação, continuar mesmo assim
        }
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
      return {} as T
    }
  }

  // Método para requisições autenticadas com renovação automática
  async authenticatedRequest<T>(url: string, options: RequestInit = {}): Promise<T> {
    if (!getAuthState) {
      throw new Error('AuthState getter não foi configurado')
    }

    const { tokens, clearAuth } = getAuthState()
    
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
        const { refreshToken, clearAuth: clearAuthFromState } = getAuthState()
        await refreshToken()
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
        
      } catch {
        const { clearAuth: clearAuthFromState } = getAuthState()
        clearAuthFromState()
        
        // Mostrar mensagem amigável ao invés de erro técnico
        await setNotificationStore()
        if (getNotificationStore) {
          const store = await getNotificationStore()
          store.showInfo(
            'Sessão encerrada',
            'Por segurança, faça login novamente para continuar',
            6000 // Duração um pouco maior para dar tempo de ler
          )
        }
        
        // Lançar erro especial para identificar sessão expirada
        throw new Error('SESSION_EXPIRED')
      }
    }

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Erro na requisição' }))
      const errorMessage = errorData.message || `Erro ${response.status}`
      
      // Mostrar notificação de erro para requisições autenticadas
      await setNotificationStore()
      if (getNotificationStore && !url.includes('/auth/logout')) {
        try {
          const store = await getNotificationStore()
          store.showError('Erro na requisição', errorMessage)
        } catch {
          // Se falhar ao mostrar notificação, continuar mesmo assim
        }
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
      return {} as T
    }
  }

  // Métodos de conveniência para diferentes tipos de requisição
  async get<T>(url: string, options?: RequestInit): Promise<T> {
    return this.request<T>(url, { ...options, method: 'GET' })
  }

  async post<T>(url: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(url, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async put<T>(url: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.request<T>(url, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async delete<T>(url: string, options?: RequestInit): Promise<T> {
    return this.request<T>(url, { ...options, method: 'DELETE' })
  }

  // Métodos autenticados
  async authGet<T>(url: string, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, { ...options, method: 'GET' })
  }

  async authPost<T>(url: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, {
      ...options,
      method: 'POST',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async authPut<T>(url: string, data?: unknown, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, {
      ...options,
      method: 'PUT',
      body: data ? JSON.stringify(data) : undefined,
    })
  }

  async authDelete<T>(url: string, options?: RequestInit): Promise<T> {
    return this.authenticatedRequest<T>(url, { ...options, method: 'DELETE' })
  }
}

// Instância singleton do cliente API
export const apiClient = new ApiClient()

// Exportar constantes para uso em outros lugares se necessário
export { AUTH_CONSTANTS }