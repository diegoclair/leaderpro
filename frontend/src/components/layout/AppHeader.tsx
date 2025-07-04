'use client'

import { useRouter } from 'next/navigation'
import { Building2, User, Settings } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Avatar, AvatarFallback } from '@/components/ui/avatar'
import { CompanySelector } from '../company/CompanySelector'
import { ThemeToggle } from './ThemeToggle'

interface AppHeaderProps {
  className?: string
}

export function AppHeader({ className }: AppHeaderProps) {
  const router = useRouter()

  return (
    <header className={`border-b bg-background px-6 py-4 ${className || ''}`}>
      <div className="flex items-center justify-between">
        {/* Logo and Company Selector */}
        <div className="flex items-center gap-6">
          <div 
            className="flex items-center gap-2 cursor-pointer hover:opacity-80 transition-opacity"
            onClick={() => router.push('/')}
          >
            <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-gradient-to-br from-blue-600 to-blue-700 shadow-lg">
              <span className="text-lg font-bold text-white">LP</span>
            </div>
            <div className="flex flex-col">
              <span className="text-xl font-bold bg-gradient-to-r from-blue-600 to-green-600 bg-clip-text text-transparent">
                LeaderPro
              </span>
              <span className="text-xs text-muted-foreground -mt-1">
                Seu coach de liderança
              </span>
            </div>
          </div>
          
          <CompanySelector />
        </div>

        {/* Navigation and User Menu */}
        <div className="flex items-center gap-4">
          {/* Main Navigation */}
          <nav className="flex items-center gap-1">
            <Button variant="ghost" size="sm" className="gap-2">
              <Building2 className="h-4 w-4" />
              Dashboard
            </Button>
            <Button variant="ghost" size="sm" className="gap-2">
              <User className="h-4 w-4" />
              Pessoas
            </Button>
          </nav>

          {/* Theme Toggle */}
          <ThemeToggle />

          {/* User Menu */}
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="relative h-9 w-9 rounded-full">
                <Avatar className="h-9 w-9">
                  <AvatarFallback className="bg-slate-100 text-slate-600">
                    DC
                  </AvatarFallback>
                </Avatar>
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent className="w-56" align="end" forceMount>
              <DropdownMenuLabel className="font-normal">
                <div className="flex flex-col space-y-1">
                  <p className="text-sm font-medium leading-none">Diego Clair</p>
                  <p className="text-xs leading-none text-muted-foreground">
                    diego93rodrigues@gmail.com
                  </p>
                </div>
              </DropdownMenuLabel>
              <DropdownMenuSeparator />
              <DropdownMenuItem className="gap-2">
                <Settings className="h-4 w-4" />
                Configurações
              </DropdownMenuItem>
              <DropdownMenuSeparator />
              <DropdownMenuItem>
                Sair
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </div>
    </header>
  )
}