'use client'

import React from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { LucideIcon } from 'lucide-react'
import { cn } from '@/lib/utils'

interface InfoCardProps {
  title: string
  icon: LucideIcon
  children: React.ReactNode
  className?: string
  headerAction?: React.ReactNode
  contentClassName?: string
}

export function InfoCard({ 
  title, 
  icon: Icon, 
  children, 
  className,
  headerAction,
  contentClassName
}: InfoCardProps) {
  return (
    <Card className={cn("flex flex-col transition-all duration-300", className)}>
      <CardHeader className="px-4 pt-1 pb-3">
        <div className="flex items-center justify-between">
          <CardTitle className="flex items-center gap-2 text-base font-semibold">
            <Icon className="h-5 w-5 text-muted-foreground" />
            {title}
          </CardTitle>
          {headerAction && (
            <div className="flex items-center gap-2">
              {headerAction}
            </div>
          )}
        </div>
      </CardHeader>
      <CardContent className={cn("flex-1", contentClassName)}>
        {children}
      </CardContent>
    </Card>
  )
}