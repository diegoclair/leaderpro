import { cn } from '@/lib/utils/cn'

interface AppLogoProps {
  size?: 'sm' | 'default' | 'lg'
  showText?: boolean
  className?: string
}

export function AppLogo({ size = 'default', showText = true, className }: AppLogoProps) {
  const sizes = {
    sm: {
      icon: 'h-8 w-8 text-sm',
      text: 'text-lg'
    },
    default: {
      icon: 'h-12 w-12 text-xl',
      text: 'text-2xl'
    },
    lg: {
      icon: 'h-16 w-16 text-2xl',
      text: 'text-3xl'
    }
  }

  return (
    <div className={cn('flex items-center space-x-3', className)}>
      <div 
        className={cn(
          'inline-flex items-center justify-center rounded-xl bg-gradient-to-br from-blue-600 to-blue-700 shadow-lg',
          sizes[size].icon
        )}
      >
        <span className={cn('font-bold text-white', sizes[size].text.replace('text-', 'text-'))}>
          LP
        </span>
      </div>
      {showText && (
        <h1 className={cn(
          'font-bold bg-gradient-to-r from-blue-600 to-green-600 bg-clip-text text-transparent',
          sizes[size].text
        )}>
          LeaderPro
        </h1>
      )}
    </div>
  )
}