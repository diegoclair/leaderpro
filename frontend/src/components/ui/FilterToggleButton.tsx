'use client'

import React from 'react'
import { Button } from '@/components/ui/button'
import { LucideIcon, X } from 'lucide-react'

interface FilterToggleButtonProps {
  id: string
  name: string
  icon: LucideIcon
  isActive: boolean
  onClick: () => void
  showRemoveIcon?: boolean
  className?: string
}

export function FilterToggleButton({
  name,
  icon: Icon,
  isActive,
  onClick,
  showRemoveIcon = true,
  className = ''
}: FilterToggleButtonProps) {
  return (
    <Button
      type="button"
      variant={isActive ? "default" : "outline"}
      size="sm"
      onClick={onClick}
      className={`h-8 text-sm gap-2 btn-interactive ${className}`}
    >
      <Icon className="h-3 w-3" />
      {name}
      {isActive && showRemoveIcon && <X className="h-3 w-3" />}
    </Button>
  )
}