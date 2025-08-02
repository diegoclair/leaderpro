// AI Chat Types
export interface ChatMessage {
  id: string
  role: 'user' | 'assistant'
  content: string
  timestamp: Date
  usageId?: number
  feedback?: 'helpful' | 'not_helpful' | 'neutral'
  isStreaming?: boolean
  error?: string
}

export interface ChatRequest {
  message: string
}

export interface ChatResponse {
  response: string
  usage_id: number
}

export interface ChatUsageDetails {
  input_tokens: number
  output_tokens: number
  total_tokens: number
  cost_usd: number
}

// AI Feedback Types
export interface FeedbackRequest {
  feedback: 'helpful' | 'not_helpful' | 'neutral'
  comment?: string
}

// AI Usage Report Types
export type UsagePeriod = 'today' | 'week' | 'month' | 'year' | 'all'

export interface UsageStats {
  total_requests: number
  total_tokens: number
  total_cost_usd: number
  average_response_time_ms: number
  feedback_stats: {
    helpful: number
    not_helpful: number
    neutral: number
  }
}

export interface UsageReportResponse {
  period: UsagePeriod
  stats: UsageStats
  daily_breakdown?: Array<{
    date: string
    requests: number
    tokens: number
    cost_usd: number
  }>
}

// Chat History Types (for future pagination)
export interface ChatHistoryRequest {
  page?: number
  limit?: number
}

export interface ChatHistoryResponse {
  messages: ChatMessage[]
  has_more: boolean
  total_count: number
}