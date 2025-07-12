import { useEffect, useState } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/stores/authStore'
import { useCompanyStore, useLoadCompanies } from '@/lib/stores/companyStore'

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
  const loadCompanies = useLoadCompanies()
  const [needsOnboarding, setNeedsOnboarding] = useState(false)
  const [companiesLoaded, setCompaniesLoaded] = useState(false)

  // Carregar empresas quando autenticado
  useEffect(() => {
    const loadUserCompanies = async () => {
      if (isAuthenticated && !companiesLoaded) {
        console.log('🔐 AuthRedirect: Usuário autenticado, carregando empresas...')
        try {
          await loadCompanies()
          console.log('✅ AuthRedirect: Empresas carregadas, total:', companies.length)
        } catch (error) {
          console.error('❌ AuthRedirect: Erro ao carregar empresas:', error)
        } finally {
          setCompaniesLoaded(true)
        }
      }
    }
    
    loadUserCompanies()
  }, [isAuthenticated, companiesLoaded, loadCompanies])

  useEffect(() => {
    // Se ainda está carregando auth, não faz nada
    if (isLoading) return
    
    // Se autenticado mas ainda não carregou empresas, aguardar
    if (isAuthenticated && !companiesLoaded) return

    // Se requer autenticação mas não está autenticado
    if (requireAuth && !isAuthenticated) {
      router.push(redirectTo || '/auth/login')
      return
    }

    // Se não requer autenticação mas está autenticado (páginas de auth)
    if (!requireAuth && isAuthenticated) {
      const hasCompanies = companies.length > 0
      
      if (!hasCompanies) {
        setNeedsOnboarding(true)
      } else {
        router.push('/dashboard')
      }
      return
    }

    // Se está autenticado e requer auth, verificar onboarding
    if (requireAuth && isAuthenticated && companiesLoaded) {
      const hasCompanies = companies.length > 0
      
      console.log('📊 Empresas carregadas:', companies.length, hasCompanies ? '- Dashboard normal' : '- Iniciando onboarding')
      
      setNeedsOnboarding(!hasCompanies)
    }
  }, [isAuthenticated, requireAuth, redirectTo, router, isLoading, companies.length, companiesLoaded])

  return {
    isAuthenticated,
    isLoading,
    needsOnboarding,
    shouldRender: requireAuth ? isAuthenticated : !isAuthenticated,
    completeOnboarding: () => setNeedsOnboarding(false)
  }
}