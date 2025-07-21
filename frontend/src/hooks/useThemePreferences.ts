'use client'

import { useEffect, useState } from 'react'
import { useTheme } from 'next-themes'
import { apiClient } from '@/lib/stores/authStore'
import { useNotificationStore } from '@/lib/stores/notificationStore'

interface UserPreferences {
  theme: string
}

export function useThemePreferences() {
  const { theme, setTheme } = useTheme()
  const [preferences, setPreferences] = useState<UserPreferences | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isSaving, setIsSaving] = useState(false)
  const { showError } = useNotificationStore()

  // Fetch user preferences from backend
  const fetchPreferences = async () => {
    try {
      setIsLoading(true)
      
      const response = await apiClient.authGet<UserPreferences>('/users/preferences')
      
      if (response?.theme) {
        setPreferences(response)
        
        // Apply the theme from backend
        setTheme(response.theme)
      } else {
        // No preferences found, use default light theme
        setTheme('light')
      }
    } catch (error) {
      console.error('Error fetching preferences:', error)
      // If error loading preferences, default to light theme
      setTheme('light')
    } finally {
      setIsLoading(false)
    }
  }

  // Update theme preference in backend
  const updateThemePreference = async (newTheme: string) => {
    if (isSaving || (preferences && preferences.theme === newTheme)) return

    try {
      setIsSaving(true)
      
      const response = await apiClient.authPut<UserPreferences>('/users/preferences', {
        theme: newTheme
      })
      
      if (response?.theme) {
        setPreferences(response)
      }
    } catch (error) {
      console.error('Error updating theme preference:', error)
      showError('Erro', 'Não foi possível salvar a preferência de tema')
      
      // Revert theme change on error
      if (preferences) {
        setTheme(preferences.theme)
      } else {
        setTheme('light') // Default fallback
      }
    } finally {
      setIsSaving(false)
    }
  }

  // Load preferences on mount
  useEffect(() => {
    fetchPreferences()
  }, [])

  // Listen to theme changes and sync with backend
  useEffect(() => {
    if (!isLoading && theme && preferences && theme !== preferences.theme) {
      updateThemePreference(theme)
    }
  }, [theme, preferences, isLoading])

  return {
    preferences,
    isLoading,
    isSaving,
    fetchPreferences,
    updateThemePreference
  }
}