import { create } from 'zustand'

export type NotificationType = 'success' | 'error' | 'warning' | 'info'

export interface Notification {
  id: string
  type: NotificationType
  title: string
  message?: string
  duration?: number // em milissegundos, 0 = não remove automaticamente
  action?: {
    label: string
    onClick: () => void
  }
}

interface NotificationState {
  notifications: Notification[]
  
  // Actions
  addNotification: (notification: Omit<Notification, 'id'>) => string
  removeNotification: (id: string) => void
  clearAllNotifications: () => void
  
  // Métodos de conveniência
  showSuccess: (title: string, message?: string, duration?: number) => string
  showError: (title: string, message?: string, duration?: number) => string
  showWarning: (title: string, message?: string, duration?: number) => string
  showInfo: (title: string, message?: string, duration?: number) => string
}

const DEFAULT_DURATION = {
  success: 4000,
  info: 4000, 
  warning: 6000,
  error: 8000, // Erros ficam mais tempo
}

let notificationId = 0

export const useNotificationStore = create<NotificationState>((set, get) => ({
  notifications: [],

  addNotification: (notification) => {
    const id = `notification-${++notificationId}`
    const duration = notification.duration ?? DEFAULT_DURATION[notification.type]
    
    const newNotification: Notification = {
      ...notification,
      id,
      duration
    }
    
    set((state) => ({
      notifications: [...state.notifications, newNotification]
    }))
    
    // Auto-remover após duração especificada (se > 0)
    if (duration > 0) {
      setTimeout(() => {
        get().removeNotification(id)
      }, duration)
    }
    
    return id
  },

  removeNotification: (id) => {
    set((state) => ({
      notifications: state.notifications.filter(n => n.id !== id)
    }))
  },

  clearAllNotifications: () => {
    set({ notifications: [] })
  },

  // Métodos de conveniência
  showSuccess: (title, message, duration) => {
    return get().addNotification({ type: 'success', title, message, duration })
  },

  showError: (title, message, duration) => {
    return get().addNotification({ type: 'error', title, message, duration })
  },

  showWarning: (title, message, duration) => {
    return get().addNotification({ type: 'warning', title, message, duration })
  },

  showInfo: (title, message, duration) => {
    return get().addNotification({ type: 'info', title, message, duration })
  },
}))

// Hooks para facilitar o uso
export const useNotifications = () => useNotificationStore(state => state.notifications)
export const useAddNotification = () => useNotificationStore(state => state.addNotification)
export const useRemoveNotification = () => useNotificationStore(state => state.removeNotification)

// Hooks para métodos específicos
export const useShowSuccess = () => useNotificationStore(state => state.showSuccess)
export const useShowError = () => useNotificationStore(state => state.showError)
export const useShowWarning = () => useNotificationStore(state => state.showWarning)
export const useShowInfo = () => useNotificationStore(state => state.showInfo)