'use client'

import React from 'react'
import { AlertCircle } from 'lucide-react'
import { cn } from '@/lib/utils'

interface ErrorMessageProps {
  message: string
  actionButton?: React.ReactNode
  className?: string
}

export function ErrorMessage({ message, actionButton, className }: ErrorMessageProps) {
  return (
    <div className={cn('text-center py-8', className)}>
      <AlertCircle className="h-12 w-12 text-destructive mx-auto mb-4" />
      <h3 className="text-lg font-medium mb-2">Erro</h3>
      <p className="text-muted-foreground mb-4">{message}</p>
      {actionButton}
    </div>
  )
}