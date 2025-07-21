'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { ErrorMessage } from '@/components/ui/ErrorMessage'
import { ThemeToggle } from '@/components/ui/ThemeToggle'
import { Settings, Palette, User } from 'lucide-react'
import { useTheme } from 'next-themes'
import { useThemePreferences } from '@/hooks/useThemePreferences'
import { AppHeader } from '@/components/layout/AppHeader'
import { useAuthRedirect } from '@/hooks/useAuthRedirect'

export default function SettingsPage() {
  const { isLoading: authLoading, shouldRender } = useAuthRedirect({ requireAuth: true })
  const { theme, setTheme } = useTheme()
  const { preferences, isLoading, isSaving, fetchPreferences } = useThemePreferences()

  // Handle theme change
  const handleThemeChange = (newTheme: string) => {
    setTheme(newTheme)
  }

  // Don't render until auth is verified
  if (!shouldRender) {
    return null
  }

  if (authLoading || isLoading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-center">
          <LoadingSpinner size="large" />
          <p className="text-muted-foreground mt-4">Carregando configurações...</p>
        </div>
      </div>
    )
  }

  // Show error state if preferences failed to load
  if (!isLoading && !preferences) {
    return (
      <div className="py-8">
        <ErrorMessage 
          message="Erro ao carregar preferências. Tente novamente."
          actionButton={
            <Button onClick={fetchPreferences} variant="outline" size="sm">
              Tentar Novamente
            </Button>
          }
        />
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-background">
      <AppHeader />
      <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-2">
          <Settings className="h-8 w-8 text-primary" />
          <h1 className="text-3xl font-bold">Configurações</h1>
        </div>
        <p className="text-muted-foreground">
          Gerencie suas preferências e configurações da conta
        </p>
      </div>

      <div className="space-y-6">
        {/* Appearance Settings */}
        <Card>
          <CardHeader>
            <CardTitle>Aparência</CardTitle>
            <CardDescription>
              Personalize como a aplicação será exibida
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label className="text-base">Tema</Label>
                <p className="text-sm text-muted-foreground">
                  Escolha entre o tema claro ou escuro
                </p>
              </div>
              <ThemeToggle
                theme={theme || 'light'}
                onThemeChange={handleThemeChange}
                disabled={isSaving}
              />
            </div>
          </CardContent>
        </Card>

        {/* Profile Settings Placeholder */}
        <Card>
          <CardHeader>
            <CardTitle>Perfil</CardTitle>
            <CardDescription>
              Gerencie as informações do seu perfil
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="flex items-center justify-between py-2">
              <div className="space-y-0.5">
                <Label className="text-base">Editar perfil</Label>
                <p className="text-sm text-muted-foreground">
                  Altere suas informações pessoais e de contato
                </p>
              </div>
              <Button variant="outline" disabled>
                Em breve
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
      </div>
    </div>
  )
}