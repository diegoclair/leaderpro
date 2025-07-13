import { cn } from '@/lib/utils/cn'
import { AlertCircle } from 'lucide-react'

interface ErrorMessageProps {
  children: React.ReactNode
  variant?: 'default' | 'destructive' | 'warning'
  className?: string
  showIcon?: boolean
}

export function ErrorMessage({ 
  children, 
  variant = 'destructive', 
  className,
  showIcon = true 
}: ErrorMessageProps) {
  const variants = {
    default: 'text-gray-600 bg-gray-50 border-gray-200',
    destructive: 'text-red-600 bg-red-50 border-red-200',
    warning: 'text-orange-600 bg-orange-50 border-orange-200'
  }

  if (!children) return null

  return (
    <div className={cn(
      'p-3 text-sm border rounded-lg flex items-start space-x-2',
      variants[variant],
      className
    )}>
      {showIcon && (
        <AlertCircle className="h-4 w-4 mt-0.5 flex-shrink-0" />
      )}
      <div className="flex-1">{children}</div>
    </div>
  )
}