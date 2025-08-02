import { useNotificationStore } from '@/lib/stores/notificationStore'

export interface ShowNotificationParams {
  type: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  duration?: number
  action?: {
    label: string
    onClick: () => void
  }
}

export const showNotification = (params: ShowNotificationParams) => {
  return useNotificationStore.getState().addNotification(params)
}