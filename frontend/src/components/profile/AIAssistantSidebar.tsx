'use client'

import React, { useState, useRef, useEffect } from 'react'
import { InfoCard } from '@/components/ui/info-card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { MarkdownRenderer } from '@/components/ui/markdown-renderer'
import { 
  Send, 
  Bot, 
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

  const chatState = useAIChat({ personUuid: person.uuid })
  const {
    messages,
    isLoading,
    error,
    sendMessage,
    sendFeedback,
    clearMessages,
    retryLastMessage,
    messagesEndRef
  } = chatState

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
    <InfoCard
      title="Assistente de Liderança IA"
      icon={Bot}
      className={cn("h-[500px]", className)}
      headerAction={null}
      contentClassName="flex-1 flex flex-col p-0"
    >
        
        <ScrollArea className="flex-1 px-3 py-4">
          {showQuickSuggestions ? (
            <div className="space-y-4">
              {/* Welcome com estilo */}
              <div className="text-center py-3">
                <div className="h-10 w-10 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center mx-auto mb-3">
                  <Bot className="h-5 w-5 text-white" />
                </div>
                <h3 className="text-sm font-medium mb-1">
                  Como posso ajudar com {person.name}?
                </h3>
                <p className="text-xs text-muted-foreground">
                  Escolha uma sugestão ou digite sua pergunta
                </p>
              </div>

              {/* Sugestões com melhor design */}
              <div className="space-y-2">
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full justify-start text-left h-auto py-3 px-3 hover:bg-blue-50 dark:hover:bg-blue-950 border-muted-foreground/20"
                  onClick={() => handleQuickSuggestion(`Que perguntas fazer na próxima 1:1 com ${person.name}?`)}
                >
                  <div className="flex items-center gap-2 w-full">
                    <span className="text-sm">💬</span>
                    <span className="text-xs font-medium">Perguntas para 1:1</span>
                  </div>
                </Button>
                
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full justify-start text-left h-auto py-3 px-3 hover:bg-green-50 dark:hover:bg-green-950 border-muted-foreground/20"
                  onClick={() => handleQuickSuggestion(`Como dar feedback construtivo para ${person.name}?`)}
                >
                  <div className="flex items-center gap-2 w-full">
                    <span className="text-sm">📝</span>
                    <span className="text-xs font-medium">Dicas de feedback</span>
                  </div>
                </Button>
                
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full justify-start text-left h-auto py-3 px-3 hover:bg-purple-50 dark:hover:bg-purple-950 border-muted-foreground/20"
                  onClick={() => handleQuickSuggestion(`Como desenvolver ${person.name} profissionalmente?`)}
                >
                  <div className="flex items-center gap-2 w-full">
                    <span className="text-sm">🎯</span>
                    <span className="text-xs font-medium">Desenvolvimento</span>
                  </div>
                </Button>
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
                      <div className="h-7 w-7 rounded-full bg-blue-100 dark:bg-blue-900 flex items-center justify-center">
                        <Bot className="h-4 w-4 text-blue-600 dark:text-blue-400" />
                      </div>
                    </div>
                  )}

                  <div className={cn(
                    "flex flex-col gap-1 max-w-[80%]",
                    message.role === 'user' && "items-end"
                  )}>
                    <div
                      className={cn(
                        "rounded-lg px-3 py-2 text-sm",
                        message.role === 'user'
                          ? "bg-blue-600 text-white"
                          : "bg-muted",
                        message.error && "bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800"
                      )}
                    >
                      {message.isStreaming ? (
                        <div className="flex items-start gap-2">
                          <MarkdownRenderer 
                            content={message.content} 
                            className="flex-1 [&>*]:my-1 [&>*:first-child]:mt-0 [&>*:last-child]:mb-0"
                          />
                          <span className="animate-pulse text-blue-500">▊</span>
                        </div>
                      ) : (
                        message.role === 'user' ? (
                          <div className="whitespace-pre-wrap">{message.content}</div>
                        ) : (
                          <MarkdownRenderer 
                            content={message.content} 
                            className="[&>*]:my-1 [&>*:first-child]:mt-0 [&>*:last-child]:mb-0"
                          />
                        )
                      )}
                    </div>

                    {message.error && (
                      <div className="flex items-center gap-2 text-xs text-red-600 dark:text-red-400">
                        <AlertCircle className="h-3 w-3" />
                        {message.error}
                      </div>
                    )}

                    <div className="flex items-center gap-2 text-xs text-muted-foreground">
                      <span>{format(message.timestamp, 'HH:mm', { locale: ptBR })}</span>
                      
                      {message.role === 'assistant' && message.usageId && !message.error && !message.isStreaming && (
                        <div className="flex items-center gap-1">
                          {message.feedback ? (
                            <span>{message.feedback === 'helpful' ? '👍' : '👎'}</span>
                          ) : (
                            <>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-5 w-5 p-0"
                                onClick={() => handleFeedback(message.id, 'helpful')}
                              >
                                <ThumbsUp className="h-3 w-3" />
                              </Button>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-5 w-5 p-0"
                                onClick={() => handleFeedback(message.id, 'not_helpful')}
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
                      <div className="h-7 w-7 rounded-full bg-gray-100 dark:bg-gray-800 flex items-center justify-center">
                        <User className="h-4 w-4 text-gray-600 dark:text-gray-400" />
                      </div>
                    </div>
                  )}
                </div>
              ))}
              <div ref={messagesEndRef} />
            </div>
          )}
        </ScrollArea>

        {/* Input area - melhor design */}
        <div className="border-t bg-muted/20 p-4">
          {error && (
            <div className="flex items-center justify-between mb-3 p-2.5 bg-red-50 dark:bg-red-950 rounded-lg border border-red-200 dark:border-red-800">
              <div className="flex items-center gap-2">
                <AlertCircle className="h-4 w-4 text-red-500" />
                <span className="text-sm text-red-600 dark:text-red-400">Erro ao enviar</span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={retryLastMessage}
                className="h-7 w-7 p-0 text-red-600 dark:text-red-400 hover:bg-red-100 dark:hover:bg-red-900"
              >
                <RefreshCw className="h-4 w-4" />
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
              className="h-11 pr-12 text-sm rounded-xl border-2 focus:border-blue-500 transition-colors shadow-sm bg-background"
            />
            <Button
              onClick={handleSendMessage}
              disabled={isLoading || !inputMessage.trim()}
              size="sm"
              className={cn(
                "absolute right-2 top-1/2 -translate-y-1/2 h-8 w-8 p-0 rounded-lg transition-all",
                (isLoading || !inputMessage.trim())
                  ? "bg-muted text-muted-foreground"
                  : "bg-blue-600 hover:bg-blue-700 text-white shadow-sm"
              )}
            >
              {isLoading ? (
                <div className="animate-spin h-4 w-4 border-2 border-current border-t-transparent rounded-full" />
              ) : (
                <Send className="h-4 w-4" />
              )}
            </Button>
          </div>

          {messages.length > 0 && (
            <div className="flex justify-center mt-3">
              <Button
                variant="ghost"
                size="sm"
                onClick={clearMessages}
                className="h-7 px-3 text-xs text-muted-foreground hover:text-red-600 hover:bg-red-50 dark:hover:bg-red-950"
              >
                <Trash2 className="h-3 w-3 mr-1" />
                Limpar conversa
              </Button>
            </div>
          )}
        </div>
      </InfoCard>
  )
}