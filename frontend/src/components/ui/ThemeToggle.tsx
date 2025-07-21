'use client'

import React from 'react'
import { cn } from '@/lib/utils'

interface ThemeToggleProps {
  theme: string
  onThemeChange: (theme: string) => void
  disabled?: boolean
  className?: string
}

export function ThemeToggle({ theme, onThemeChange, disabled = false, className }: ThemeToggleProps) {
  const isDark = theme === 'dark'

  const handleToggle = () => {
    if (!disabled) {
      onThemeChange(isDark ? 'light' : 'dark')
    }
  }

  return (
    <div className={cn("flex items-center justify-between", className)}>
      {/* Theme Options */}
      <div className="flex items-center gap-1 bg-muted rounded-lg p-1">
        {/* Light Theme Button */}
        <button
          onClick={() => !disabled && onThemeChange('light')}
          disabled={disabled}
          className={cn(
            "px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200",
            "disabled:cursor-not-allowed disabled:opacity-50",
            !isDark
              ? "bg-background text-foreground shadow-sm"
              : "text-muted-foreground hover:text-foreground"
          )}
        >
          Claro
        </button>
        
        {/* Dark Theme Button */}
        <button
          onClick={() => !disabled && onThemeChange('dark')}
          disabled={disabled}
          className={cn(
            "px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200",
            "disabled:cursor-not-allowed disabled:opacity-50",
            isDark
              ? "bg-background text-foreground shadow-sm"
              : "text-muted-foreground hover:text-foreground"
          )}
        >
          Escuro
        </button>
      </div>
    </div>
  )
}