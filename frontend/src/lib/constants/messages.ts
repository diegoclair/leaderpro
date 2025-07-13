// Error messages
export const ERROR_MESSAGES = {
  // Authentication
  LOGIN_FAILED: 'Erro ao fazer login. Verifique suas credenciais.',
  REGISTER_FAILED: 'Erro ao criar conta. Tente novamente.',
  SESSION_EXPIRED: 'Sessão expirada, faça login novamente.',
  UNAUTHORIZED: 'Você não tem permissão para acessar este recurso.',
  
  // Validation
  REQUIRED_FIELD: 'Este campo é obrigatório.',
  INVALID_EMAIL: 'Email inválido.',
  PASSWORD_TOO_SHORT: 'A senha deve ter pelo menos 8 caracteres.',
  PASSWORDS_DONT_MATCH: 'As senhas não coincidem.',
  INVALID_PHONE: 'Formato de telefone inválido.',
  INVALID_DATE: 'Data inválida.',
  
  // Company
  COMPANY_CREATE_FAILED: 'Erro ao criar empresa.',
  COMPANY_UPDATE_FAILED: 'Erro ao atualizar empresa.',
  COMPANY_DELETE_FAILED: 'Erro ao excluir empresa.',
  COMPANY_NOT_FOUND: 'Empresa não encontrada.',
  
  // Person
  PERSON_CREATE_FAILED: 'Erro ao criar pessoa.',
  PERSON_UPDATE_FAILED: 'Erro ao atualizar pessoa.',
  PERSON_DELETE_FAILED: 'Erro ao excluir pessoa.',
  PERSON_NOT_FOUND: 'Pessoa não encontrada.',
  
  // Notes
  NOTE_CREATE_FAILED: 'Erro ao criar anotação.',
  NOTE_UPDATE_FAILED: 'Erro ao atualizar anotação.',
  NOTE_DELETE_FAILED: 'Erro ao excluir anotação.',
  
  // Generic
  NETWORK_ERROR: 'Erro de conexão. Verifique sua internet.',
  UNKNOWN_ERROR: 'Erro desconhecido. Tente novamente.',
  SERVER_ERROR: 'Erro interno do servidor.',
} as const

// Success messages
export const SUCCESS_MESSAGES = {
  // Authentication
  LOGIN_SUCCESS: 'Login realizado com sucesso!',
  REGISTER_SUCCESS: 'Conta criada com sucesso!',
  LOGOUT_SUCCESS: 'Logout realizado com sucesso!',
  
  // Company
  COMPANY_CREATED: 'Empresa criada com sucesso!',
  COMPANY_UPDATED: 'Empresa atualizada com sucesso!',
  COMPANY_DELETED: 'Empresa excluída com sucesso!',
  
  // Person
  PERSON_CREATED: 'Pessoa criada com sucesso!',
  PERSON_UPDATED: 'Pessoa atualizada com sucesso!',
  PERSON_DELETED: 'Pessoa excluída com sucesso!',
  
  // Notes
  NOTE_CREATED: 'Anotação criada com sucesso!',
  NOTE_UPDATED: 'Anotação atualizada com sucesso!',
  NOTE_DELETED: 'Anotação excluída com sucesso!',
} as const

// Loading messages
export const LOADING_MESSAGES = {
  LOADING: 'Carregando...',
  SAVING: 'Salvando...',
  DELETING: 'Excluindo...',
  UPDATING: 'Atualizando...',
  PROCESSING: 'Processando...',
  AUTHENTICATING: 'Autenticando...',
  CREATING: 'Criando...',
} as const

// Confirmation messages
export const CONFIRMATION_MESSAGES = {
  DELETE_COMPANY: 'Tem certeza que deseja excluir esta empresa?',
  DELETE_PERSON: 'Tem certeza que deseja excluir esta pessoa?',
  DELETE_NOTE: 'Tem certeza que deseja excluir esta anotação?',
  LOGOUT: 'Tem certeza que deseja sair?',
  UNSAVED_CHANGES: 'Há alterações não salvas. Deseja sair mesmo assim?',
} as const

// All messages combined
export const MESSAGES = {
  ERROR: ERROR_MESSAGES,
  SUCCESS: SUCCESS_MESSAGES,
  LOADING: LOADING_MESSAGES,
  CONFIRMATION: CONFIRMATION_MESSAGES,
} as const