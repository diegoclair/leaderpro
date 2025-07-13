import { cn } from '@/lib/utils/cn'

interface LoadingSpinnerProps {
  size?: 'sm' | 'default' | 'lg'
  className?: string
}

export function LoadingSpinner({ size = 'default', className }: LoadingSpinnerProps) {
  const sizes = {
    sm: 'h-4 w-4',
    default: 'h-8 w-8',
    lg: 'h-12 w-12'
  }

  return (
    <div 
      className={cn(
        'animate-spin rounded-full border-b-2 border-blue-600',
        sizes[size],
        className
      )}
      aria-label="Carregando..."
    />
  )
}

interface LoadingPageProps {
  children?: React.ReactNode
  size?: 'sm' | 'default' | 'lg'
}

export function LoadingPage({ children = 'Carregando...', size = 'lg' }: LoadingPageProps) {
  return (
    <div className="min-h-screen flex items-center justify-center flex-col space-y-4">
      <LoadingSpinner size={size} />
      {children && (
        <span className="text-sm text-gray-600 dark:text-gray-400">
          {children}
        </span>
      )}
    </div>
  )
}