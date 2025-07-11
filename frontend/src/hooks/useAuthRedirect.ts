import { useEffect, useState } from 'react'
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
  const { isAuthenticated, tokens, isLoading } = useAuthStore()
  const { companies } = useCompanyStore()
  const [needsOnboarding, setNeedsOnboarding] = useState(false)

  useEffect(() => {
    // Se ainda está carregando, não faz nada
    if (isLoading) return

    // Se requer autenticação mas não está autenticado
    if (requireAuth && !isAuthenticated) {
      router.push(redirectTo || '/auth/login')
      return
    }

    // Se não requer autenticação mas está autenticado (páginas de auth)
    if (!requireAuth && isAuthenticated) {
      // Verificar se precisa de onboarding
      const hasCompletedOnboarding = localStorage.getItem('onboarding_completed') === 'true'
      const hasCompanies = companies.length > 0
      
      if (!hasCompletedOnboarding || !hasCompanies) {
        setNeedsOnboarding(true)
      } else {
        router.push('/dashboard')
      }
      return
    }

    // Se está autenticado e requer auth, verificar onboarding
    if (requireAuth && isAuthenticated) {
      const hasCompletedOnboarding = localStorage.getItem('onboarding_completed') === 'true'
      const hasCompanies = companies.length > 0
      
      if (!hasCompletedOnboarding || !hasCompanies) {
        setNeedsOnboarding(true)
      }
    }
  }, [isAuthenticated, requireAuth, redirectTo, router, isLoading, companies])

  return {
    isAuthenticated,
    isLoading,
    needsOnboarding,
    shouldRender: requireAuth ? isAuthenticated : !isAuthenticated,
    completeOnboarding: () => setNeedsOnboarding(false)
  }
}