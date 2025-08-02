'use client'

import React, { useState, useRef, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { 
  Send, 
  Bot, 
  Sparkles, 
  ThumbsUp, 
  ThumbsDown, 
  RefreshCw,
  Trash2,
  User,
  AlertCircle
} from 'lucide-react'
import { Person } from '@/lib/types'
import { useAIChat } from '@/hooks/useAIChat'
import { cn } from '@/lib/utils'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

interface AIAssistantSidebarProps {
  person: Person
  className?: string
}

export function AIAssistantSidebar({ person, className }: AIAssistantSidebarProps) {
  const [inputMessage, setInputMessage] = useState('')
  const inputRef = useRef<HTMLInputElement>(null)

  const {
    messages,
    isLoading,
    error,
    sendMessage,
    sendFeedback,
    clearMessages,
    retryLastMessage,
    messagesEndRef
  } = useAIChat({ personUuid: person.uuid })

  // Auto-focus input after sending
  useEffect(() => {
    if (!isLoading) {
      inputRef.current?.focus()
    }
  }, [isLoading])

  const handleSendMessage = () => {
    if (inputMessage.trim() && !isLoading) {
      sendMessage(inputMessage)
      setInputMessage('')
    }
  }

  const handleQuickSuggestion = (suggestion: string) => {
    setInputMessage(suggestion)
    inputRef.current?.focus()
  }

  const handleFeedback = (messageId: string, feedback: 'helpful' | 'not_helpful') => {
    sendFeedback(messageId, feedback)
  }

  const showQuickSuggestions = messages.length === 0

  return (
    <Card className={cn("flex flex-col h-[600px]", className)}>
      <CardContent className="flex-1 flex flex-col p-0">
        
        <ScrollArea className="flex-1 p-4">
          {showQuickSuggestions ? (
            <div className="space-y-4">
              {/* Welcome message */}
              <div className="text-center py-6">
                <Bot className="h-12 w-12 text-blue-600 mx-auto mb-4" />
                <h3 className="text-lg font-semibold mb-2">
                  Como posso ajudar com {person.name}?
                </h3>
                <p className="text-sm text-muted-foreground">
                  Use as sugestões abaixo ou faça sua própria pergunta
                </p>
              </div>

              {/* Quick Suggestions */}
              <div className="space-y-2">
                <div className="flex items-center gap-2 mb-3">
                  <Sparkles className="h-4 w-4 text-muted-foreground" />
                  <span className="text-sm font-medium text-muted-foreground">
                    Sugestões Rápidas
                  </span>
                </div>
                <div className="grid gap-2">
                  <Button
                    variant="outline"
                    className="justify-start text-left h-auto py-2.5 px-3"
                    onClick={() => handleQuickSuggestion(`Que perguntas fazer na próxima 1:1 com ${person.name}?`)}
                  >
                    <div>
                      <div className="font-medium text-sm mb-0.5">💬 Perguntas para 1:1</div>
                      <div className="text-xs text-muted-foreground">
                        Sugestões personalizadas para sua próxima reunião
                      </div>
                    </div>
                  </Button>
                  <Button
                    variant="outline"
                    className="justify-start text-left h-auto py-2.5 px-3"
                    onClick={() => handleQuickSuggestion(`Como dar feedback construtivo para ${person.name} sobre pontualidade?`)}
                  >
                    <div>
                      <div className="font-medium text-sm mb-0.5">📝 Dicas de feedback</div>
                      <div className="text-xs text-muted-foreground">
                        Como abordar tópicos sensíveis de forma construtiva
                      </div>
                    </div>
                  </Button>
                  <Button
                    variant="outline"
                    className="justify-start text-left h-auto py-2.5 px-3"
                    onClick={() => handleQuickSuggestion(`Quais são as melhores formas de desenvolver ${person.name} profissionalmente?`)}
                  >
                    <div>
                      <div className="font-medium text-sm mb-0.5">🎯 Desenvolvimento</div>
                      <div className="text-xs text-muted-foreground">
                        Plano de crescimento e desenvolvimento de carreira
                      </div>
                    </div>
                  </Button>
                </div>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              {messages.map((message) => (
                <div
                  key={message.id}
                  className={cn(
                    "flex gap-3",
                    message.role === 'user' && "justify-end"
                  )}
                >
                  {message.role === 'assistant' && (
                    <div className="flex-shrink-0">
                      <div className="h-8 w-8 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                        <Bot className="h-5 w-5 text-blue-600 dark:text-blue-400" />
                      </div>
                    </div>
                  )}

                  <div className={cn(
                    "flex flex-col gap-1 max-w-[80%]",
                    message.role === 'user' && "items-end"
                  )}>
                    <div
                      className={cn(
                        "rounded-lg px-4 py-2",
                        message.role === 'user'
                          ? "bg-blue-600 text-white"
                          : "bg-muted",
                        message.error && "bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800"
                      )}
                    >
                      {message.isStreaming ? (
                        <div className="flex items-center gap-2">
                          <span className="whitespace-pre-wrap">{message.content}</span>
                          <span className="inline-flex">
                            <span className="animate-pulse">▊</span>
                          </span>
                        </div>
                      ) : (
                        <div className="whitespace-pre-wrap">{message.content}</div>
                      )}
                    </div>

                    {/* Error state */}
                    {message.error && (
                      <div className="flex items-center gap-2 text-xs text-red-600 dark:text-red-400">
                        <AlertCircle className="h-3 w-3" />
                        {message.error}
                      </div>
                    )}

                    {/* Timestamp and feedback */}
                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span>
                        {format(message.timestamp, 'HH:mm', { locale: ptBR })}
                      </span>
                      
                      {/* Feedback buttons for AI messages */}
                      {message.role === 'assistant' && message.usageId && !message.error && !message.isStreaming && (
                        <div className="flex items-center gap-1 ml-2">
                          {message.feedback ? (
                            <span className="text-xs">
                              {message.feedback === 'helpful' ? '👍' : '👎'} Obrigado!
                            </span>
                          ) : (
                            <>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-6 w-6 p-0"
                                onClick={() => handleFeedback(message.id, 'helpful')}
                                title="Útil"
                              >
                                <ThumbsUp className="h-3 w-3" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-6 w-6 p-0"
                                onClick={() => handleFeedback(message.id, 'not_helpful')}
                                title="Não útil"
                              >
                                <ThumbsDown className="h-3 w-3" />
                              </Button>
                            </>
                          )}
                        </div>
                      )}
                    </div>
                  </div>

                  {message.role === 'user' && (
                    <div className="flex-shrink-0">
                      <div className="h-8 w-8 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                        <User className="h-5 w-5 text-gray-600 dark:text-gray-400" />
                      </div>
                    </div>
                  )}
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
          )}
        </ScrollArea>

        {/* Input area */}
        <div className="border-t p-4">
          {error && (
            <div className="flex items-center justify-between mb-3 text-sm text-red-600 dark:text-red-400">
              <span>Erro ao enviar mensagem</span>
              <Button
                variant="ghost"
                size="sm"
                onClick={retryLastMessage}
                className="h-7 px-2"
              >
                <RefreshCw className="h-3 w-3 mr-1" />
                Tentar novamente
              </Button>
            </div>
          )}

          <div className="relative">
            <Input
              ref={inputRef}
              placeholder={`Pergunte sobre ${person.name}...`}
              value={inputMessage}
              onChange={(e) => setInputMessage(e.target.value)}
              onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && handleSendMessage()}
              disabled={isLoading}
              className="min-h-[56px] py-4 pl-4 pr-14 text-base resize-none rounded-2xl border-2 focus:border-blue-500 transition-colors"
            />
            <Button
              onClick={handleSendMessage}
              disabled={isLoading || !inputMessage.trim()}
              size="icon"
              className={cn(
                "absolute right-2 top-1/2 -translate-y-1/2 h-10 w-10 rounded-full transition-all",
                (isLoading || !inputMessage.trim())
                  ? "bg-gray-200 dark:bg-gray-700 text-gray-400 cursor-not-allowed"
                  : "bg-blue-600 hover:bg-blue-700 text-white shadow-md hover:shadow-lg"
              )}
            >
              {isLoading ? (
                <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-current" />
              ) : (
                <Send className="h-4 w-4" />
              )}
            </Button>
          </div>

          <div className="flex items-center justify-between mt-3">
            <Badge variant="secondary" className="text-xs">
              <Bot className="h-3 w-3 mr-1" />
              {isLoading ? 'Processando...' : 'Online'}
            </Badge>
            <div className="flex items-center gap-2">
              {messages.length > 0 && (
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={clearMessages}
                  className="h-7 px-2 text-xs text-muted-foreground hover:text-red-600"
                  title="Limpar conversa"
                >
                  <Trash2 className="h-3 w-3 mr-1" />
                  Limpar
                </Button>
              )}
              <span className="text-xs text-muted-foreground">
                Contexto: {person.name}
              </span>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}