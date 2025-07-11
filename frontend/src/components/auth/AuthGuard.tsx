'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuthStore } from '@/lib/stores/authStore'

interface AuthGuardProps {
  children: React.ReactNode
  requireAuth?: boolean
}

export function AuthGuard({ children, requireAuth = true }: AuthGuardProps) {
  const router = useRouter()
  const { isAuthenticated, tokens, getProfile, isLoading } = useAuthStore()

  useEffect(() => {
    const checkAuth = async () => {
      // Se requer autenticação mas não está autenticado
      if (requireAuth && !isAuthenticated) {
        router.push('/auth/login')
        return
      }

      // Se não requer autenticação e está autenticado, redirecionar para dashboard
      if (!requireAuth && isAuthenticated) {
        router.push('/')
        return
      }

      // Se está autenticado mas não tem dados do usuário, buscar perfil
      if (isAuthenticated && tokens && !useAuthStore.getState().user) {
        try {
          await getProfile()
        } catch (error) {
          console.error('Erro ao buscar perfil:', error)
          router.push('/auth/login')
        }
      }
    }

    checkAuth()
  }, [isAuthenticated, tokens, requireAuth, router, getProfile])

  // Se está carregando, mostrar loading
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  // Se requer autenticação mas não está autenticado, não renderizar nada
  // (será redirecionado pelo useEffect)
  if (requireAuth && !isAuthenticated) {
    return null
  }

  // Se não requer autenticação mas está autenticado, não renderizar nada
  // (será redirecionado pelo useEffect)
  if (!requireAuth && isAuthenticated) {
    return null
  }

  return <>{children}</>
}