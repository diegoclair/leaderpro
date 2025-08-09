'use client'

import { useState, useCallback, useRef, useEffect } from 'react'
import { apiClient } from '@/lib/stores/authStore'
import { useCompanyStore } from '@/lib/stores/companyStore'
import { API_ENDPOINTS } from '@/lib/constants/api-endpoints'
import { ChatMessage, ChatRequest, ChatResponse, FeedbackRequest } from '@/lib/types/ai'
import { showNotification } from '@/lib/utils/notifications'

interface UseAIChatOptions {
  personUuid: string
  onError?: (error: Error) => void
}

export interface AIChatState {
  messages: ChatMessage[]
  isLoading: boolean
  error: string | null
  sendMessage: (message: string) => Promise<void>
  sendFeedback: (messageId: string, feedback: 'helpful' | 'not_helpful' | 'neutral', comment?: string) => Promise<void>
  clearMessages: () => void
  retryLastMessage: () => void
  messagesEndRef: React.RefObject<HTMLDivElement | null>
}

export function useAIChat({ personUuid, onError }: UseAIChatOptions) {
  const activeCompany = useCompanyStore(state => state.activeCompany)
  const [messages, setMessages] = useState<ChatMessage[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const messagesEndRef = useRef<HTMLDivElement>(null)

  // Auto-scroll to bottom when new messages arrive
  const scrollToBottom = useCallback(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [])

  useEffect(() => {
    scrollToBottom()
  }, [messages, scrollToBottom])

  // Send message to AI
  const sendMessage = useCallback(async (content: string) => {
    if (!activeCompany?.uuid || !content.trim() || isLoading) return

    const userMessage: ChatMessage = {
      id: `user-${Date.now()}`,
      role: 'user',
      content: content.trim(),
      timestamp: new Date(),
    }

    // Add user message immediately
    setMessages(prev => [...prev, userMessage])
    setIsLoading(true)
    setError(null)

    // Create assistant message with streaming state
    const assistantMessage: ChatMessage = {
      id: `assistant-${Date.now()}`,
      role: 'assistant',
      content: '',
      timestamp: new Date(),
      isStreaming: true,
    }

    setMessages(prev => [...prev, assistantMessage])

    try {
      const request: ChatRequest = { message: content.trim() }
      const response = await apiClient.authPost<ChatResponse>(
        API_ENDPOINTS.AI.CHAT(activeCompany.uuid, personUuid),
        request
      )

      // Simulate streaming effect
      await simulateStreaming(response.response, assistantMessage.id)

      // Update message with usage ID for feedback
      setMessages(prev => 
        prev.map(msg => 
          msg.id === assistantMessage.id 
            ? { ...msg, usageId: response.usage_id, isStreaming: false }
            : msg
        )
      )
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro ao enviar mensagem'
      
      // Update assistant message with error
      setMessages(prev => 
        prev.map(msg => 
          msg.id === assistantMessage.id 
            ? { 
                ...msg, 
                content: 'Desculpe, ocorreu um erro ao processar sua mensagem.', 
                error: errorMessage,
                isStreaming: false 
              }
            : msg
        )
      )

      setError(errorMessage)
      onError?.(err instanceof Error ? err : new Error(errorMessage))
      
      // Não mostrar notificação adicional se for erro de sessão expirada
      // pois o apiClient já mostra "Sessão encerrada"
      if (errorMessage !== 'SESSION_EXPIRED') {
        showNotification({
          type: 'error',
          title: 'Erro ao enviar mensagem',
          message: errorMessage,
        })
      }
    } finally {
      setIsLoading(false)
    }
  }, [activeCompany?.uuid, personUuid, isLoading, onError])

  // Simulate streaming effect for better UX
  const simulateStreaming = useCallback(async (fullText: string, messageId: string) => {
    const words = fullText.split(' ')
    let currentText = ''

    for (let i = 0; i < words.length; i++) {
      currentText += (i > 0 ? ' ' : '') + words[i]
      
      setMessages(prev => 
        prev.map(msg => 
          msg.id === messageId 
            ? { ...msg, content: currentText }
            : msg
        )
      )

      // Variable delay for more natural feel
      const delay = Math.random() * 50 + 20
      await new Promise(resolve => setTimeout(resolve, delay))
    }
  }, [])

  // Send feedback for a message
  const sendFeedback = useCallback(async (
    messageId: string, 
    feedback: 'helpful' | 'not_helpful' | 'neutral',
    comment?: string
  ) => {
    if (!activeCompany?.uuid) return

    const message = messages.find(msg => msg.id === messageId)
    if (!message?.usageId || message.feedback) return

    try {
      const request: FeedbackRequest = { feedback, comment }
      
      await apiClient.authPost(
        API_ENDPOINTS.AI.FEEDBACK(activeCompany.uuid, message.usageId),
        request
      )

      // Update message with feedback
      setMessages(prev => 
        prev.map(msg => 
          msg.id === messageId 
            ? { ...msg, feedback }
            : msg
        )
      )

      showNotification({
        type: 'success',
        title: 'Feedback enviado',
        message: 'Obrigado pelo seu feedback!',
      })
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Erro desconhecido'
      
      // Não mostrar notificação adicional se for erro de sessão expirada
      if (errorMessage !== 'SESSION_EXPIRED') {
        showNotification({
          type: 'error',
          title: 'Erro ao enviar feedback',
          message: 'Não foi possível enviar seu feedback. Tente novamente.',
        })
      }
    }
  }, [activeCompany?.uuid, messages])

  // Clear all messages
  const clearMessages = useCallback(() => {
    setMessages([])
    setError(null)
  }, [])

  // Retry last message if it failed
  const retryLastMessage = useCallback(() => {
    const lastUserMessage = [...messages].reverse().find(msg => msg.role === 'user')
    if (lastUserMessage && error) {
      // Remove the error message
      setMessages(prev => prev.filter(msg => !msg.error))
      sendMessage(lastUserMessage.content)
    }
  }, [messages, error, sendMessage])

  return {
    messages,
    isLoading,
    error,
    sendMessage,
    sendFeedback,
    clearMessages,
    retryLastMessage,
    messagesEndRef,
  }
}