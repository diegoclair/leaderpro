import { formatDistanceToNow, format } from 'date-fns'
import { ptBR } from 'date-fns/locale'

export const formatTimeAgo = (date: Date): string => {
  return formatDistanceToNow(date, { 
    locale: ptBR,
    addSuffix: true 
  })
}

export const formatTimeAgoWithoutSuffix = (date: Date): string => {
  return formatDistanceToNow(date, { 
    locale: ptBR,
    addSuffix: false 
  })
}

export const formatShortDate = (date: Date): string => {
  return format(date, 'dd/MM/yyyy', { locale: ptBR })
}

export const formatDateTime = (date: Date): string => {
  return format(date, 'dd/MM/yyyy HH:mm', { locale: ptBR })
}

export const getMockDaysAgo = (): number => {
  // Mock data - in real app this would come from actual sessions
  return Math.floor(Math.random() * 30) + 1
}

export const getMockAverageDays = (): number => {
  // Mock calculation - in real app this would be based on actual meeting history
  return Math.floor(Math.random() * 14) + 7 // 7-21 days average
}