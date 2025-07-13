import { Button } from '@/components/ui/button'
import { LoadingSpinner } from '@/components/ui/loading-spinner'
import { cn } from '@/lib/utils/cn'

interface SubmitButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  isLoading?: boolean
  loadingText?: string
  children: React.ReactNode
  variant?: 'default' | 'destructive' | 'outline' | 'secondary' | 'ghost' | 'link'
  size?: 'default' | 'sm' | 'lg' | 'icon'
  className?: string
}

export function SubmitButton({
  isLoading = false,
  loadingText = 'Salvando...',
  children,
  disabled,
  className,
  variant = 'default',
  size = 'default',
  ...props
}: SubmitButtonProps) {
  return (
    <Button
      {...props}
      type="submit"
      disabled={disabled || isLoading}
      variant={variant}
      size={size}
      className={cn('relative', className)}
    >
      {isLoading && (
        <LoadingSpinner size="sm" className="mr-2" />
      )}
      {isLoading ? loadingText : children}
    </Button>
  )
}