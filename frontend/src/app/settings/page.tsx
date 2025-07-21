'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { ErrorMessage } from '@/components/ui/ErrorMessage'
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
            <div className="flex items-center gap-2">
              <Palette className="h-5 w-5" />
              <CardTitle>Aparência</CardTitle>
            </div>
            <CardDescription>
              Personalize a aparência da aplicação
            </CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="theme">Tema</Label>
              <Select
                value={theme || 'light'}
                onValueChange={handleThemeChange}
                disabled={isSaving}
              >
                <SelectTrigger id="theme" className="w-full">
                  <SelectValue placeholder="Selecione um tema" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="light">Claro</SelectItem>
                  <SelectItem value="dark">Escuro</SelectItem>
                </SelectContent>
              </Select>
              <p className="text-sm text-muted-foreground">
                {(theme || 'light') === 'light' && 'Tema claro selecionado'}
                {(theme || 'light') === 'dark' && 'Tema escuro selecionado'}
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Profile Settings Placeholder */}
        <Card>
          <CardHeader>
            <div className="flex items-center gap-2">
              <User className="h-5 w-5" />
              <CardTitle>Perfil</CardTitle>
            </div>
            <CardDescription>
              Gerencie as informações do seu perfil
            </CardDescription>
          </CardHeader>
          <CardContent>
            <p className="text-muted-foreground">
              Edição de perfil será implementada em breve...
            </p>
          </CardContent>
        </Card>
      </div>
      </div>
    </div>
  )
}