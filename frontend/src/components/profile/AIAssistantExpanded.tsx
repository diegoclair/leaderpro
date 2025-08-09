'use client'

import React, { useState, useRef, useEffect } from 'react'
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Textarea } from '@/components/ui/textarea'
import { MarkdownRenderer } from '@/components/ui/markdown-renderer'
import { 
  Send, 
  Bot, 
  User,
  Sparkles, 
  ThumbsUp, 
  ThumbsDown, 
  RefreshCw,
  Trash2,
  AlertCircle,
  Copy,
  Check
} from 'lucide-react'
import { Person } from '@/lib/types'
import { AIChatState } from '@/hooks/useAIChat'
import { cn } from '@/lib/utils'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

interface AIAssistantExpandedProps {
  person: Person
  open: boolean
  onOpenChange: (open: boolean) => void
  chatState: AIChatState
}

export function AIAssistantExpanded({ person, open, onOpenChange, chatState }: AIAssistantExpandedProps) {
  const [inputMessage, setInputMessage] = useState('')
  const [copiedMessageId, setCopiedMessageId] = useState<string | null>(null)
  const inputRef = useRef<HTMLTextAreaElement>(null)

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
    if (!isLoading && open) {
      inputRef.current?.focus()
    }
  }, [isLoading, open])

  const handleSendMessage = () => {
    if (inputMessage.trim() && !isLoading) {
      sendMessage(inputMessage)
      setInputMessage('')
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      handleSendMessage()
    }
  }

  const handleQuickSuggestion = (suggestion: string) => {
    setInputMessage(suggestion)
    inputRef.current?.focus()
  }

  const handleFeedback = (messageId: string, feedback: 'helpful' | 'not_helpful') => {
    sendFeedback(messageId, feedback)
  }

  const handleCopyMessage = async (content: string, messageId: string) => {
    try {
      await navigator.clipboard.writeText(content)
      setCopiedMessageId(messageId)
      setTimeout(() => setCopiedMessageId(null), 2000)
    } catch (err) {
      console.error('Failed to copy:', err)
    }
  }

  const showQuickSuggestions = messages.length === 0

  const quickSuggestions = [
    {
      icon: '💬',
      title: 'Perguntas para 1:1',
      description: 'Sugestões personalizadas para sua próxima reunião',
      text: `Que perguntas fazer na próxima 1:1 com ${person.name}?`
    },
    {
      icon: '📝',
      title: 'Dicas de feedback',
      description: 'Como abordar tópicos sensíveis de forma construtiva',
      text: `Como dar feedback construtivo para ${person.name} sobre pontualidade?`
    },
    {
      icon: '🎯',
      title: 'Desenvolvimento',
      description: 'Plano de crescimento e desenvolvimento de carreira',
      text: `Quais são as melhores formas de desenvolver ${person.name} profissionalmente?`
    },
    {
      icon: '⭐',
      title: 'Motivação',
      description: 'Estratégias para manter engajamento alto',
      text: `Como posso motivar melhor ${person.name} no trabalho?`
    },
    {
      icon: '🤝',
      title: 'Relacionamento',
      description: 'Construir um relacionamento de trabalho mais forte',
      text: `Como melhorar minha comunicação com ${person.name}?`
    },
    {
      icon: '📈',
      title: 'Performance',
      description: 'Análise de desempenho e melhoria contínua',
      text: `Com base no histórico, como está a performance de ${person.name}?`
    }
  ]

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-7xl w-[95vw] h-[90vh] flex flex-col p-0">
        <DialogHeader className="px-6 py-4 border-b">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <Bot className="h-6 w-6 text-blue-600" />
              <div>
                <DialogTitle className="text-xl">
                  Coach de Liderança - {person.name}
                </DialogTitle>
                <p className="text-sm text-muted-foreground mt-1">
                  Assistente inteligente para desenvolvimento de liderança
                </p>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">
                <Bot className="h-3 w-3 mr-1" />
                {isLoading ? 'Processando...' : 'Online'}
              </Badge>
            </div>
          </div>
        </DialogHeader>

        <div className="flex-1 flex flex-col min-h-0">
          <ScrollArea className="flex-1 px-6">
            {showQuickSuggestions ? (
              <div className="py-8">
                {/* Welcome message */}
                <div className="text-center mb-12">
                  <div className="h-20 w-20 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center mx-auto mb-6">
                    <Sparkles className="h-10 w-10 text-white" />
                  </div>
                  <h2 className="text-2xl font-bold mb-3">
                    Como posso ajudar você com {person.name}?
                  </h2>
                  <p className="text-muted-foreground max-w-md mx-auto">
                    Sou seu coach de liderança pessoal. Tenho contexto completo sobre sua equipe 
                    e posso oferecer conselhos personalizados baseados no histórico de interações.
                  </p>
                </div>

                {/* Quick Suggestions Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-3 max-w-6xl mx-auto">
                  {quickSuggestions.map((suggestion, index) => (
                    <Button
                      key={index}
                      variant="outline"
                      className="h-auto p-4 text-left justify-start hover:bg-muted/50 transition-all duration-200 hover:scale-[1.02] min-h-[80px]"
                      onClick={() => handleQuickSuggestion(suggestion.text)}
                    >
                      <div className="flex items-start gap-3 w-full">
                        <span className="text-xl flex-shrink-0 mt-0.5">{suggestion.icon}</span>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium text-sm mb-1.5 leading-tight">{suggestion.title}</div>
                          <div className="text-xs text-muted-foreground leading-relaxed line-clamp-2">
                            {suggestion.description}
                          </div>
                        </div>
                      </div>
                    </Button>
                  ))}
                </div>
              </div>
            ) : (
              <div className="space-y-6 py-6">
                {messages.map((message) => (
                  <div
                    key={message.id}
                    className={cn(
                      "flex gap-4",
                      message.role === 'user' && "justify-end"
                    )}
                  >
                    {message.role === 'assistant' && (
                      <div className="flex-shrink-0 mt-2">
                        <div className="h-8 w-8 rounded-full bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center">
                          <Bot className="h-5 w-5 text-white" />
                        </div>
                      </div>
                    )}

                    <div className={cn(
                      "flex flex-col gap-2 max-w-[90%] lg:max-w-[80%]",
                      message.role === 'user' && "items-end"
                    )}>
                      <div
                        className={cn(
                          "rounded-2xl px-4 py-3 relative group",
                          message.role === 'user'
                            ? "bg-blue-600 text-white"
                            : "bg-muted border",
                          message.error && "bg-red-50 dark:bg-red-950 border-red-200 dark:border-red-800"
                        )}
                      >
                        {message.isStreaming ? (
                          <div className="flex items-start gap-2">
                            <MarkdownRenderer 
                              content={message.content} 
                              className="flex-1"
                            />
                            <div className="flex-shrink-0 mt-1">
                              <span className="inline-flex animate-pulse text-blue-500">▊</span>
                            </div>
                          </div>
                        ) : (
                          <>
                            {message.role === 'user' ? (
                              <div className="whitespace-pre-wrap">{message.content}</div>
                            ) : (
                              <MarkdownRenderer content={message.content} />
                            )}
                            
                            {/* Copy button for assistant messages */}
                            {message.role === 'assistant' && !message.error && (
                              <Button
                                variant="ghost"
                                size="sm"
                                className="absolute top-2 right-2 h-8 w-8 p-0 opacity-0 group-hover:opacity-100 transition-opacity"
                                onClick={() => handleCopyMessage(message.content, message.id)}
                              >
                                {copiedMessageId === message.id ? (
                                  <Check className="h-3 w-3" />
                                ) : (
                                  <Copy className="h-3 w-3" />
                                )}
                              </Button>
                            )}
                          </>
                        )}
                      </div>

                      {/* Error state */}
                      {message.error && (
                        <div className="flex items-center gap-2 text-sm text-red-600 dark:text-red-400 px-2">
                          <AlertCircle className="h-4 w-4" />
                          {message.error}
                        </div>
                      )}

                      {/* Timestamp and feedback */}
                      <div className="flex items-center gap-3 text-xs text-muted-foreground px-2">
                        <span>
                          {format(message.timestamp, 'HH:mm', { locale: ptBR })}
                        </span>
                        
                        {/* Feedback buttons for AI messages */}
                        {message.role === 'assistant' && message.usageId && !message.error && !message.isStreaming && (
                          <div className="flex items-center gap-1">
                            {message.feedback ? (
                              <span className="text-xs">
                                {message.feedback === 'helpful' ? '👍 Útil' : '👎 Não útil'}
                              </span>
                            ) : (
                              <div className="flex items-center gap-1">
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="h-6 w-6 p-0 hover:bg-green-100 dark:hover:bg-green-900"
                                  onClick={() => handleFeedback(message.id, 'helpful')}
                                  title="Útil"
                                >
                                  <ThumbsUp className="h-3 w-3" />
                                </Button>
                                <Button
                                  variant="ghost"
                                  size="sm"
                                  className="h-6 w-6 p-0 hover:bg-red-100 dark:hover:bg-red-900"
                                  onClick={() => handleFeedback(message.id, 'not_helpful')}
                                  title="Não útil"
                                >
                                  <ThumbsDown className="h-3 w-3" />
                                </Button>
                              </div>
                            )}
                          </div>
                        )}
                      </div>
                    </div>

                    {message.role === 'user' && (
                      <div className="flex-shrink-0 mt-2">
                        <div className="h-8 w-8 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center">
                          <User className="h-5 w-5 text-gray-600 dark:text-gray-300" />
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
          <div className="border-t px-6 py-4 bg-background">
            {error && (
              <div className="flex items-center justify-between mb-4 p-3 bg-red-50 dark:bg-red-950 border border-red-200 dark:border-red-800 rounded-lg">
                <div className="flex items-center gap-2 text-red-600 dark:text-red-400">
                  <AlertCircle className="h-4 w-4" />
                  <span className="text-sm">Erro ao enviar mensagem</span>
                </div>
                <Button
                  variant="ghost"
                  size="sm"
                  onClick={retryLastMessage}
                  className="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300"
                >
                  <RefreshCw className="h-3 w-3 mr-1" />
                  Tentar novamente
                </Button>
              </div>
            )}

            <div className="relative">
              <Textarea
                ref={inputRef}
                placeholder={`Faça uma pergunta sobre ${person.name}...`}
                value={inputMessage}
                onChange={(e) => setInputMessage(e.target.value)}
                onKeyDown={handleKeyDown}
                disabled={isLoading}
                className="min-h-[60px] max-h-[200px] py-4 pl-4 pr-14 text-base resize-none rounded-2xl border-2 focus:border-blue-500 transition-colors"
                rows={1}
              />
              <Button
                onClick={handleSendMessage}
                disabled={isLoading || !inputMessage.trim()}
                size="icon"
                className={cn(
                  "absolute right-2 bottom-2 h-10 w-10 rounded-full transition-all",
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
              <div className="flex items-center gap-2 text-xs text-muted-foreground">
                <span>Pressione Enter para enviar, Shift+Enter para nova linha</span>
              </div>
              <div className="flex items-center gap-2">
                {messages.length > 0 && (
                  <Button
                    variant="ghost"
                    size="sm"
                    onClick={clearMessages}
                    className="h-8 px-3 text-xs text-muted-foreground hover:text-red-600"
                  >
                    <Trash2 className="h-3 w-3 mr-1" />
                    Limpar conversa
                  </Button>
                )}
              </div>
            </div>
          </div>
        </div>
      </DialogContent>
    </Dialog>
  )
}