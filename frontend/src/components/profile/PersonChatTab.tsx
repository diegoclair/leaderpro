'use client'

import React, { useState } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Send } from 'lucide-react'
import { Person } from '@/lib/types'

interface PersonChatTabProps {
  person: Person
}

export function PersonChatTab({ person }: PersonChatTabProps) {
  const [chatMessage, setChatMessage] = useState('')

  const handleSendMessage = () => {
    if (chatMessage.trim()) {
      // TODO: Send message to AI chat
      console.log('Sending message to AI:', chatMessage)
      setChatMessage('')
    }
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle>Chat com IA - Insights sobre {person.name}</CardTitle>
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {/* Mock AI conversation */}
          <div className="space-y-3">
            <div className="bg-blue-50 dark:bg-blue-950/30 rounded-lg p-3">
              <p className="text-sm">
                <strong>IA:</strong> Com base no histórico, {person.name} tem demonstrado interesse em crescimento profissional. 
                Algumas sugestões para a próxima 1:1:
              </p>
              <ul className="text-sm mt-2 space-y-1 ml-4">
                <li>• Discutir oportunidades de liderança</li>
                <li>• Explorar áreas de interesse específicas</li>
                <li>• Definir metas de desenvolvimento</li>
              </ul>
            </div>
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="chat-message">Pergunte algo sobre {person.name}</Label>
            <div className="flex gap-2">
              <Input 
                id="chat-message"
                placeholder="Como posso ajudar esta pessoa a crescer profissionalmente?"
                value={chatMessage}
                onChange={(e) => setChatMessage(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleSendMessage()}
              />
              <Button onClick={handleSendMessage} size="sm">
                <Send className="h-4 w-4" />
              </Button>
            </div>
          </div>
        </div>
      </CardContent>
    </Card>
  )
}