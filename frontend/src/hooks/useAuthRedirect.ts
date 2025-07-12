import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/stores/authStore'
import { useCompanyStore } from '@/lib/stores/companyStore'

interface UseAuthRedirectOptions {
  requireAuth?: boolean
  redirectTo?: string
}

export function useAuthRedirect({ 
  requireAuth = true, 
  redirectTo 
}: UseAuthRedirectOptions = {}) {
  const router = useRouter()
  const { isAuthenticated, isLoading } = useAuthStore()
  const companies = useCompanyStore(state => state.companies)

  useEffect(() => {
    if (isLoading) return

    // Redirecionar se não autenticado e requer auth
    if (requireAuth && !isAuthenticated) {
      router.push(redirectTo || '/auth/login')
      return
    }

    // Redirecionar se autenticado mas está em página de auth
    if (!requireAuth && isAuthenticated && companies.length > 0) {
      router.push('/dashboard')
    }
  }, [isAuthenticated, isLoading, requireAuth, router, redirectTo, companies.length])

  // Estados simples
  const needsOnboarding = isAuthenticated && companies.length === 0
  const shouldRender = requireAuth ? isAuthenticated : !isAuthenticated

  return {
    isAuthenticated,
    isLoading,
    needsOnboarding,
    shouldRender,
    completeOnboarding: () => router.push('/dashboard')
  }
}