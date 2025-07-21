import { NextResponse } from 'next/server'
import type { NextRequest } from 'next/server'

// Rotas que NÃO precisam de autenticação
// const publicRoutes = ['/auth/login', '/auth/register']

// Rotas que SÃO de autenticação (para redirecionar se já logado)
// const authRoutes = ['/auth/login', '/auth/register']

export function middleware(_request: NextRequest) {
  // const { pathname } = request.nextUrl
  
  // Verificar se tem token (do localStorage não é acessível no middleware)
  // Então vamos apenas redirecionar com base na URL atual
  
  // Se está tentando acessar página de auth mas já está logado (verificamos no client)
  // deixamos o AuthGuard do client tratar isso
  
  // Se está tentando acessar rota protegida sem token (verificamos no client)
  // deixamos o AuthGuard do client redirecionar
  
  return NextResponse.next()
}

export const config = {
  matcher: [
    /*
     * Match all request paths except for the ones starting with:
     * - api (API routes)
     * - _next/static (static files)
     * - _next/image (image optimization files)
     * - favicon.ico (favicon file)
     */
    '/((?!api|_next/static|_next/image|favicon.ico).*)',
  ],
}