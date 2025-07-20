'use client'

import React, { useState, useEffect } from 'react'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { MentionsTextarea } from '@/components/ui/mentions-textarea'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { Person } from '@/lib/types'
import { TimelineActivity } from './SimpleActivityCard'

// Função para traduzir categorias de feedback
const translateFeedbackCategory = (category: string): string => {
  const translations: Record<string, string> = {
    'performance': 'Performance',
    'behavior': 'Comportamento', 
    'skill': 'Habilidade',
    'collaboration': 'Colaboração'
  }
  return translations[category] || category
}

interface EditNoteModalProps {
  isOpen: boolean
  onClose: () => void
  activity: TimelineActivity | null
  allPeople: Person[]
  onSave: (updatedActivity: TimelineActivity) => Promise<void>
}

export function EditNoteModal({ 
  isOpen, 
  onClose, 
  activity, 
  allPeople, 
  onSave 
}: EditNoteModalProps) {
  const [content, setContent] = useState('')
  const [type, setType] = useState<'feedback' | 'one_on_one' | 'observation'>('observation')
  const [feedbackType, setFeedbackType] = useState<'positive' | 'constructive' | 'neutral' | ''>('')
  const [feedbackCategory, setFeedbackCategory] = useState('')
  const [isLoading, setIsLoading] = useState(false)

  // Reset form when activity changes
  useEffect(() => {
    if (activity) {
      setContent(activity.content || '')
      setType(activity.type as 'feedback' | 'one_on_one' | 'observation')
      setFeedbackType(activity.feedback_type || '')
      setFeedbackCategory(activity.feedback_category || 'none')
    }
  }, [activity])

  const handleSave = async () => {
    if (!activity || !content.trim()) return

    setIsLoading(true)
    try {
      const updatedActivity: TimelineActivity = {
        ...activity,
        content: content.trim(),
        type,
        feedback_type: feedbackType || undefined,
        feedback_category: feedbackCategory === 'none' ? undefined : feedbackCategory || undefined,
      }

      await onSave(updatedActivity)
      onClose()
    } catch (error) {
      console.error('Error saving note:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleClose = () => {
    onClose()
  }

  if (!activity) return null

  return (
    <Dialog open={isOpen} onOpenChange={handleClose}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>Editar Anotação</DialogTitle>
          <DialogDescription>
            Faça as alterações necessárias na anotação abaixo.
          </DialogDescription>
        </DialogHeader>

        <div className="space-y-4">
          {/* Tipo da Nota */}
          <div className="space-y-2">
            <Label htmlFor="type">Tipo da Anotação</Label>
            <Select value={type} onValueChange={(value) => setType(value as any)}>
              <SelectTrigger>
                <SelectValue placeholder="Selecione o tipo" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="observation">
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-green-500 rounded-full"></span>
                    Observação
                  </div>
                </SelectItem>
                <SelectItem value="feedback">
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-yellow-500 rounded-full"></span>
                    Feedback
                  </div>
                </SelectItem>
                <SelectItem value="one_on_one">
                  <div className="flex items-center gap-2">
                    <span className="w-2 h-2 bg-blue-500 rounded-full"></span>
                    1:1
                  </div>
                </SelectItem>
              </SelectContent>
            </Select>
          </div>

          {/* Tipo de Feedback - só aparece quando type é feedback */}
          {type === 'feedback' && (
            <div className="space-y-2">
              <Label htmlFor="feedback-type">Tipo de Feedback</Label>
              <Select value={feedbackType} onValueChange={(value) => setFeedbackType(value as any)}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione o tipo de feedback" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="positive">
                    <div className="flex items-center gap-2">
                      ✨ Positivo
                    </div>
                  </SelectItem>
                  <SelectItem value="constructive">
                    <div className="flex items-center gap-2">
                      🔨 Construtivo
                    </div>
                  </SelectItem>
                  <SelectItem value="neutral">
                    <div className="flex items-center gap-2">
                      ⚪ Neutro
                    </div>
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          {/* Categoria do Feedback - só aparece quando type é feedback */}
          {type === 'feedback' && (
            <div className="space-y-2">
              <Label htmlFor="feedback-category">Categoria (Opcional)</Label>
              <Select value={feedbackCategory} onValueChange={setFeedbackCategory}>
                <SelectTrigger>
                  <SelectValue placeholder="Selecione uma categoria" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="none">Nenhuma categoria</SelectItem>
                  <SelectItem value="performance">Performance</SelectItem>
                  <SelectItem value="behavior">Comportamento</SelectItem>
                  <SelectItem value="skill">Habilidade</SelectItem>
                  <SelectItem value="collaboration">Colaboração</SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          {/* Conteúdo da Nota */}
          <div className="space-y-2">
            <Label htmlFor="content">Conteúdo da Anotação</Label>
            <MentionsTextarea
              value={content}
              onChange={setContent}
              people={allPeople}
              placeholder={
                type === 'feedback' 
                  ? "Descreva o feedback. Use @nome para mencionar outras pessoas..."
                  : type === 'one_on_one'
                  ? "Registre os pontos principais da conversa. Use @nome para mencionar outras pessoas..."
                  : "Descreva sua observação. Use @nome para mencionar outras pessoas..."
              }
              className="min-h-[120px]"
            />
            <p className="text-xs text-muted-foreground">
              Digite @ para mencionar pessoas do seu time.
            </p>
          </div>

          {/* Preview das badges como na visualização */}
          <div className="flex gap-2 flex-wrap">
            {type === 'feedback' && feedbackType && (
              <Badge variant={feedbackType === 'positive' ? 'default' : 'secondary'} className="text-xs">
                {feedbackType === 'positive' ? '✨ Positivo' : 
                 feedbackType === 'constructive' ? '🔨 Construtivo' : '⚪ Neutro'}
              </Badge>
            )}
            {type === 'feedback' && feedbackCategory && feedbackCategory !== 'none' && (
              <Badge variant="outline" className="text-xs">
                {translateFeedbackCategory(feedbackCategory)}
              </Badge>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button type="button" variant="outline" onClick={handleClose}>
            Cancelar
          </Button>
          <Button 
            type="button" 
            onClick={handleSave}
            disabled={!content.trim() || isLoading}
          >
            {isLoading ? (
              <>
                <LoadingSpinner size="small" className="mr-2" />
                Salvando...
              </>
            ) : (
              'Salvar Alterações'
            )}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}