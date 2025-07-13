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
  const { isAuthenticated, isLoading: authLoading, hasHydrated } = useAuthStore()
  const { companies, isLoading: companyLoading } = useCompanyStore()
  const [hasInitialized, setHasInitialized] = useState(false)

  useEffect(() => {
    // Wait for auth hydration and both stores to be loaded
    if (authLoading || companyLoading || !hasHydrated) return
    
    // Mark as initialized after first load
    if (!hasInitialized) {
      setHasInitialized(true)
      return
    }

    // Only redirect after we've confirmed the auth state
    if (requireAuth && !isAuthenticated) {
      router.push(redirectTo || '/auth/login')
      return
    }

    // Redirect if authenticated but on auth page
    if (!requireAuth && isAuthenticated) {
      router.push('/dashboard')
    }
  }, [isAuthenticated, authLoading, companyLoading, hasHydrated, requireAuth, router, redirectTo, companies.length, hasInitialized])

  // Estados simples - aguardar hidratação completa
  const needsOnboarding = isAuthenticated && companies.length === 0 && hasInitialized && hasHydrated
  const shouldRender = hasInitialized && hasHydrated && (requireAuth ? isAuthenticated : !isAuthenticated)
  const isLoading = authLoading || companyLoading || !hasInitialized || !hasHydrated

  return {
    isAuthenticated,
    isLoading,
    needsOnboarding,
    shouldRender,
    completeOnboarding: () => router.push('/dashboard')
  }
}