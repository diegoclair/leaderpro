'use client'

import React, { useEffect, useState } from 'react'
import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react'
import { useNotifications, useRemoveNotification, useNotificationStore, type Notification, type NotificationType } from '@/lib/stores/notificationStore'
import { Button } from './button'

const NotificationIcon = ({ type }: { type: NotificationType }) => {
  const iconProps = { size: 20 }
  
  switch (type) {
    case 'success':
      return <CheckCircle {...iconProps} className="text-green-500" />
    case 'error':
      return <AlertCircle {...iconProps} className="text-red-500" />
    case 'warning':
      return <AlertTriangle {...iconProps} className="text-yellow-500" />
    case 'info':
      return <Info {...iconProps} className="text-blue-500" />
    default:
      return <Info {...iconProps} />
  }
}

const getNotificationStyles = (type: NotificationType) => {
  const baseStyles = "border-l-4 shadow-lg backdrop-blur-sm"
  
  switch (type) {
    case 'success':
      return `${baseStyles} bg-green-50/90 border-green-500 dark:bg-green-900/20 dark:border-green-400`
    case 'error':
      return `${baseStyles} bg-red-50/90 border-red-500 dark:bg-red-900/20 dark:border-red-400`
    case 'warning':
      return `${baseStyles} bg-yellow-50/90 border-yellow-500 dark:bg-yellow-900/20 dark:border-yellow-400`
    case 'info':
      return `${baseStyles} bg-blue-50/90 border-blue-500 dark:bg-blue-900/20 dark:border-blue-400`
    default:
      return `${baseStyles} bg-gray-50/90 border-gray-500 dark:bg-gray-900/20 dark:border-gray-400`
  }
}

interface NotificationItemProps {
  notification: Notification
  onRemove: (id: string) => void
}

const NotificationItem: React.FC<NotificationItemProps> = ({ notification, onRemove }) => {
  const [isVisible, setIsVisible] = useState(false)
  const [isLeaving, setIsLeaving] = useState(false)

  useEffect(() => {
    // Trigger entrada animation
    const timer = setTimeout(() => setIsVisible(true), 10)
    return () => clearTimeout(timer)
  }, [])

  const handleRemove = () => {
    setIsLeaving(true)
    setTimeout(() => {
      onRemove(notification.id)
    }, 300) // Tempo da animação de saída
  }

  return (
    <div
      className={`
        transform transition-all duration-300 ease-in-out
        ${isVisible && !isLeaving 
          ? 'translate-x-0 opacity-100' 
          : 'translate-x-full opacity-0'
        }
      `}
    >
      <div
        className={`
          ${getNotificationStyles(notification.type)}
          rounded-lg p-4 mb-3 max-w-sm w-full
          flex items-start gap-3
        `}
      >
        <NotificationIcon type={notification.type} />
        
        <div className="flex-1 min-w-0">
          <div className="flex items-start justify-between gap-2">
            <div className="flex-1">
              <h4 className="text-sm font-semibold text-gray-900 dark:text-gray-100">
                {notification.title}
              </h4>
              {notification.message && (
                <p className="text-sm text-gray-600 dark:text-gray-300 mt-1">
                  {notification.message}
                </p>
              )}
            </div>
            
            <Button
              variant="ghost"
              size="sm"
              onClick={handleRemove}
              className="h-6 w-6 p-0 hover:bg-black/10 dark:hover:bg-white/10"
            >
              <X size={14} />
            </Button>
          </div>
          
          {notification.action && (
            <div className="mt-3">
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  notification.action?.onClick()
                  handleRemove()
                }}
                className="text-xs"
              >
                {notification.action.label}
              </Button>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export const NotificationContainer: React.FC = () => {
  const notifications = useNotifications()
  const removeNotification = useRemoveNotification()

  if (notifications.length === 0) {
    return null
  }

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {notifications.map((notification) => (
        <NotificationItem
          key={notification.id}
          notification={notification}
          onRemove={removeNotification}
        />
      ))}
    </div>
  )
}

// Componente de exemplo para demonstração (pode ser removido)
export const NotificationDemo: React.FC = () => {
  const { showSuccess, showError } = useNotificationStore()
  
  return (
    <div className="space-y-2">
      <Button onClick={() => showSuccess('Sucesso!', 'Operação realizada com sucesso')}>
        Success
      </Button>
      <Button onClick={() => showError('Erro!', 'Algo deu errado')} variant="destructive">
        Error
      </Button>
    </div>
  )
}