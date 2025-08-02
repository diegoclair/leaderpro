// Authentication endpoints
export const AUTH_ENDPOINTS = {
  LOGIN: '/auth/login',
  LOGOUT: '/auth/logout',
  REFRESH: '/auth/refresh',
} as const

// User endpoints
export const USER_ENDPOINTS = {
  REGISTER: '/users',
  PROFILE: '/users/profile',
  UPDATE_PROFILE: '/users/profile',
} as const

// Company endpoints  
export const COMPANY_ENDPOINTS = {
  LIST: '/companies',
  CREATE: '/companies',
  GET_BY_ID: (uuid: string) => `/companies/${uuid}`,
  UPDATE: (uuid: string) => `/companies/${uuid}`,
  DELETE: (uuid: string) => `/companies/${uuid}`,
} as const

// Person endpoints
export const PERSON_ENDPOINTS = {
  LIST: (companyUuid: string) => `/companies/${companyUuid}/people`,
  CREATE: (companyUuid: string) => `/companies/${companyUuid}/people`,
  GET_BY_ID: (companyUuid: string, personUuid: string) => `/companies/${companyUuid}/people/${personUuid}`,
  UPDATE: (companyUuid: string, personUuid: string) => `/companies/${companyUuid}/people/${personUuid}`,
  DELETE: (companyUuid: string, personUuid: string) => `/companies/${companyUuid}/people/${personUuid}`,
  SEARCH: (companyUuid: string, query: string) => `/companies/${companyUuid}/people?search=${encodeURIComponent(query)}`,
} as const

// Note endpoints
export const NOTE_ENDPOINTS = {
  CREATE: (companyUuid: string, personUuid: string) => `/companies/${companyUuid}/people/${personUuid}/notes`,
  TIMELINE: (companyUuid: string, personUuid: string) => `/companies/${companyUuid}/people/${personUuid}/timeline`,
  UPDATE: (companyUuid: string, personUuid: string, noteUuid: string) => `/companies/${companyUuid}/people/${personUuid}/notes/${noteUuid}`,
  DELETE: (companyUuid: string, personUuid: string, noteUuid: string) => `/companies/${companyUuid}/people/${personUuid}/notes/${noteUuid}`,
} as const

// AI endpoints
export const AI_ENDPOINTS = {
  CHAT: (companyUuid: string, personUuid: string) => 
    `/companies/${companyUuid}/people/${personUuid}/ai/chat`,
  USAGE: (companyUuid: string, period?: string) => 
    `/companies/${companyUuid}/ai/usage${period ? `?period=${period}` : ''}`,
  FEEDBACK: (companyUuid: string, usageId: number) => 
    `/companies/${companyUuid}/ai/usage/${usageId}/feedback`,
} as const

// Utility function to get all endpoints
export const API_ENDPOINTS = {
  AUTH: AUTH_ENDPOINTS,
  USER: USER_ENDPOINTS,
  COMPANY: COMPANY_ENDPOINTS,
  PERSON: PERSON_ENDPOINTS,
  NOTE: NOTE_ENDPOINTS,
  AI: AI_ENDPOINTS,
} as const