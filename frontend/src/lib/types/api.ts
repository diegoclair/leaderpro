// API Response Types for LeaderPro
// This file centralizes all API response types to replace 'any' usage in apiClient

import { User } from '../stores/authStore'
import { Person, Company } from './index'

// Base response interfaces
export interface ApiError {
  message: string
  code?: string
  details?: any
}

export interface PaginationMeta {
  current_page: number
  records_per_page: number
  total_records: number
  total_pages: number
}

// Authentication API Responses
export interface AuthTokens {
  accessToken: string
  accessTokenExpiresAt: string
  refreshToken: string
  refreshTokenExpiresAt: string
}

// Backend format for auth tokens (snake_case)
export interface BackendAuthTokens {
  access_token: string
  access_token_expires_at: string
  refresh_token: string
  refresh_token_expires_at: string
}

export interface LoginResponse {
  user: User
  auth: BackendAuthTokens  // Backend envia 'auth', não 'tokens'
}

export interface RegisterResponse {
  user: User
  auth: BackendAuthTokens  // Backend envia 'auth', não 'tokens'
}

export interface RefreshTokenResponse {
  access_token: string
  access_token_expires_at: string
}

// User API Responses
export type UserProfileResponse = User

// Company API Responses - Backend returns snake_case
export interface ApiCompany {
  uuid: string
  name: string
  industry: string
  size: string
  role: string
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface CompaniesResponse {
  companies: ApiCompany[]
}

export interface CompanyCreateResponse {
  company: ApiCompany
}

export interface CompanyResponse {
  company: ApiCompany
}

// People API Responses - Backend returns snake_case
export interface ApiPerson {
  uuid: string
  name: string
  email?: string
  position?: string
  department?: string
  phone?: string
  birthday?: string
  start_date?: string
  is_manager: boolean
  manager_uuid?: string
  notes?: string
  has_kids: boolean
  gender?: 'male' | 'female' | 'other'
  interests?: string
  personality?: string
  last_one_on_one_date?: string
  created_at: string
  updated_at: string
  age?: number
  tenure?: number
}

export interface PeopleResponse {
  people: ApiPerson[]
}

export interface PersonResponse {
  person: ApiPerson
}

export interface PersonCreateResponse {
  person: ApiPerson
}

// Dashboard API Responses
export interface DashboardStats {
  total_people: number
  one_on_ones_this_month: number
  feedbacks_count_this_month: number
  average_frequency_days: number
  last_meeting_date?: string
}

export interface DashboardResponse {
  people: ApiPerson[]
  stats: DashboardStats
}

// Timeline API Responses
export interface TimelineActivity {
  uuid: string
  type: 'feedback' | 'one_on_one' | 'observation' | 'mention'
  content: string
  author_name: string
  author_id?: string
  created_at: string
  feedback_type?: 'positive' | 'constructive' | 'neutral'
  feedback_category?: string
  person_name?: string
  source_person_name?: string
  entry_source: 'direct' | 'mention'
}

export interface TimelineResponse {
  data: TimelineActivity[]
  pagination: PaginationMeta
}

// Notes API Responses
export interface NoteCreateData {
  content: string
  type: 'feedback' | 'one_on_one' | 'observation'
  feedback_type?: 'positive' | 'constructive' | 'neutral'
  feedback_category?: string
}

export interface NoteResponse {
  note: TimelineActivity
}

export interface NoteUpdateData {
  content: string
  type: 'feedback' | 'one_on_one' | 'observation'
  feedback_type?: 'positive' | 'constructive' | 'neutral'
}

// Generic responses for operations without specific data
export type EmptyResponse = Record<string, never>

// Helper type for void operations (like logout, delete)  
export type VoidResponse = EmptyResponse

// Union types for API responses
export type ApiResponse<T = any> = T | ApiError