'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Send, Bot, Lightbulb, MessageSquare, Sparkles } from 'lucide-react'
import { Person } from '@/lib/types'

interface AIAssistantSidebarProps {
  person: Person
  className?: string
}

export function AIAssistantSidebar({ person, className }: AIAssistantSidebarProps) {
  const [chatMessage, setChatMessage] = useState('')
  const [isThinking, setIsThinking] = useState(false)

  const handleSendMessage = async () => {
    if (chatMessage.trim() && !isThinking) {
      setIsThinking(true)
      // TODO: Send message to AI chat
      console.log('Sending message to AI:', chatMessage)
      
      // Simulate AI response delay
      setTimeout(() => {
        setIsThinking(false)
        setChatMessage('')
      }, 1500)
    }
  }

  const handleQuickSuggestion = (suggestion: string) => {
    setChatMessage(suggestion)
  }

  return (
    <Card className={className}>
      <CardHeader className="pb-3">
        <CardTitle className="flex items-center gap-2 text-base">
          <Bot className="h-4 w-4 text-blue-600" />
          IA Assistant
        </CardTitle>
      </CardHeader>
      <CardContent className="space-y-4">
        {/* AI Insights */}
        <div className="space-y-3">
          <div className="bg-blue-50 dark:bg-blue-950/30 border-l-4 border-blue-500 rounded-r-lg p-3">
            <div className="flex items-start gap-2 mb-2">
              <Lightbulb className="h-4 w-4 text-blue-600 mt-0.5 flex-shrink-0" />
              <div>
                <p className="text-sm font-medium text-blue-900 dark:text-blue-100 mb-1">
                  Insights sobre {person.name}
                </p>
                <p className="text-xs text-blue-700 dark:text-blue-200">
                  Com base no histórico, sugiro focar em desenvolvimento de liderança na próxima 1:1.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Quick Suggestions */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 mb-2">
            <Sparkles className="h-3 w-3 text-muted-foreground" />
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
              Sugestões Rápidas
            </span>
          </div>
          <div className="space-y-1">
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-start text-xs h-7 px-2"
              onClick={() => handleQuickSuggestion('Que perguntas fazer na próxima 1:1?')}
            >
              💬 Perguntas para 1:1
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-start text-xs h-7 px-2"
              onClick={() => handleQuickSuggestion('Como dar feedback construtivo?')}
            >
              📝 Dicas de feedback
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="w-full justify-start text-xs h-7 px-2"
              onClick={() => handleQuickSuggestion('Sugestões de desenvolvimento profissional')}
            >
              🎯 Desenvolvimento
            </Button>
          </div>
        </div>

        {/* Chat Input */}
        <div className="space-y-2">
          <div className="flex items-center gap-2 mb-2">
            <MessageSquare className="h-3 w-3 text-muted-foreground" />
            <span className="text-xs font-medium text-muted-foreground uppercase tracking-wide">
              Pergunte à IA
            </span>
          </div>
          <div className="space-y-2">
            <Input 
              placeholder={`Pergunte sobre ${person.name}...`}
              value={chatMessage}
              onChange={(e) => setChatMessage(e.target.value)}
              onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
              disabled={isThinking}
              className="text-sm"
            />
            <Button 
              onClick={handleSendMessage} 
              size="sm"
              className="w-full"
              disabled={isThinking || !chatMessage.trim()}
            >
              {isThinking ? (
                <>
                  <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-white mr-2"></div>
                  Pensando...
                </>
              ) : (
                <>
                  <Send className="h-3 w-3 mr-2" />
                  Enviar
                </>
              )}
            </Button>
          </div>
        </div>

        {/* Quick Status */}
        <div className="pt-2 border-t">
          <div className="flex items-center justify-between">
            <Badge variant="secondary" className="text-xs">
              <Bot className="h-3 w-3 mr-1" />
              Online
            </Badge>
            <span className="text-xs text-muted-foreground">
              Contexto: {person.name}
            </span>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}