// Typed API Client - Provides strongly typed methods for common API calls
// This layer sits on top of the base apiClient to provide type safety

import { apiClient } from './client'
import type {
  // Auth types
  LoginResponse,
  RegisterResponse,
  RefreshTokenResponse,
  UserProfileResponse,
  
  // Company types
  CompaniesResponse,
  CompanyCreateResponse,
  CompanyResponse,
  
  // People types
  PeopleResponse,
  PersonResponse,
  PersonCreateResponse,
  
  // Dashboard types
  DashboardResponse,
  
  // Timeline types
  TimelineResponse,
  
  // Notes types
  NoteCreateData,
  NoteResponse,
  NoteUpdateData,
  
  // Generic types
  VoidResponse,
  EmptyResponse
} from '../types/api'

import type { User } from '../stores/authStore'
import type { Person, Company } from '../types'

export class TypedApiClient {
  // Authentication endpoints
  static async login(email: string, password: string): Promise<LoginResponse> {
    return apiClient.post<LoginResponse>('/auth/login', { email, password })
  }

  static async register(data: {
    name: string
    email: string
    password: string
    phone?: string
  }): Promise<RegisterResponse> {
    return apiClient.post<RegisterResponse>('/users', data)
  }

  static async refreshToken(refreshToken: string): Promise<RefreshTokenResponse> {
    return apiClient.post<RefreshTokenResponse>('/auth/refresh-token', { 
      refresh_token: refreshToken 
    })
  }

  static async logout(): Promise<VoidResponse> {
    return apiClient.authPost<VoidResponse>('/auth/logout')
  }

  // User endpoints
  static async getUserProfile(): Promise<UserProfileResponse> {
    return apiClient.authGet<UserProfileResponse>('/users/profile')
  }

  // Company endpoints
  static async getCompanies(): Promise<Company[]> {
    return apiClient.authGet<Company[]>('/companies')
  }

  static async createCompany(companyData: {
    name: string
    industry: string
    size: string
    role: string
  }): Promise<CompanyCreateResponse> {
    return apiClient.authPost<CompanyCreateResponse>('/companies', companyData)
  }

  static async getCompany(companyId: string): Promise<CompanyResponse> {
    return apiClient.authGet<CompanyResponse>(`/companies/${companyId}`)
  }

  // People endpoints
  static async getPeople(companyUuid: string): Promise<Person[]> {
    return apiClient.authGet<Person[]>(`/companies/${companyUuid}/people`)
  }

  static async createPerson(companyUuid: string, personData: Partial<Person>): Promise<PersonCreateResponse> {
    return apiClient.authPost<PersonCreateResponse>(`/companies/${companyUuid}/people`, personData)
  }

  static async updatePerson(companyUuid: string, personUuid: string, personData: Partial<Person>): Promise<VoidResponse> {
    return apiClient.authPut<VoidResponse>(`/companies/${companyUuid}/people/${personUuid}`, personData)
  }

  static async deletePerson(companyUuid: string, personUuid: string): Promise<VoidResponse> {
    return apiClient.authDelete<VoidResponse>(`/companies/${companyUuid}/people/${personUuid}`)
  }

  // Dashboard endpoints
  static async getDashboard(companyUuid: string): Promise<DashboardResponse> {
    return apiClient.authGet<DashboardResponse>(`/dashboard?company_uuid=${companyUuid}`)
  }

  // Timeline endpoints
  static async getTimeline(
    companyId: string, 
    personUuid: string, 
    params?: URLSearchParams
  ): Promise<TimelineResponse> {
    const url = `/companies/${companyId}/people/${personUuid}/timeline${params ? `?${params.toString()}` : ''}`
    return apiClient.authGet<TimelineResponse>(url)
  }

  // Notes endpoints
  static async createNote(
    companyId: string, 
    personId: string, 
    noteData: NoteCreateData
  ): Promise<NoteResponse> {
    return apiClient.authPost<NoteResponse>(`/companies/${companyId}/people/${personId}/notes`, noteData)
  }

  static async updateNote(
    companyId: string, 
    personUuid: string, 
    noteUuid: string, 
    noteData: NoteUpdateData
  ): Promise<VoidResponse> {
    return apiClient.authPut<VoidResponse>(`/companies/${companyId}/people/${personUuid}/notes/${noteUuid}`, noteData)
  }

  static async deleteNote(
    companyId: string, 
    personUuid: string, 
    noteUuid: string
  ): Promise<VoidResponse> {
    return apiClient.authDelete<VoidResponse>(`/companies/${companyId}/people/${personUuid}/notes/${noteUuid}`)
  }
}

// Export a default instance for convenience
export const typedApi = TypedApiClient