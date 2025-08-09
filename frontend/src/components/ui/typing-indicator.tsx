'use client'

import React from 'react'
import { cn } from '@/lib/utils'

interface TypingIndicatorProps {
  className?: string
  size?: 'sm' | 'md' | 'lg'
}

export function TypingIndicator({ className, size = 'md' }: TypingIndicatorProps) {
  const sizeClasses = {
    sm: 'w-2 h-2',
    md: 'w-2.5 h-2.5', 
    lg: 'w-3 h-3'
  }

  return (
    <div className={cn("flex items-center gap-1", className)}>
      <div className={cn(
        "rounded-full bg-current opacity-70 animate-bounce",
        sizeClasses[size]
      )} style={{ animationDelay: '0ms' }} />
      <div className={cn(
        "rounded-full bg-current opacity-70 animate-bounce",
        sizeClasses[size]
      )} style={{ animationDelay: '150ms' }} />
      <div className={cn(
        "rounded-full bg-current opacity-70 animate-bounce", 
        sizeClasses[size]
      )} style={{ animationDelay: '300ms' }} />
    </div>
  )
}