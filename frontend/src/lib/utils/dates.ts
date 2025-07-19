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

// Smart relative date formatting (dias/meses/anos atrás)
export const formatDateRelative = (date: Date | string | null | undefined, options?: { 
  showNeverText?: boolean 
}): string => {
  const { showNeverText = false } = options || {}
  
  if (!date) {
    return showNeverText ? 'Nunca' : 'Data não informada'
  }
  
  const dateObj = typeof date === 'string' ? new Date(date) : date
  const now = new Date()
  
  // Check if both dates are on the same calendar day
  const isToday = dateObj.toDateString() === now.toDateString()
  
  if (isToday) {
    return 'hoje'
  }
  
  // Calculate difference in calendar days (not 24-hour periods)
  const dateStart = new Date(dateObj.getFullYear(), dateObj.getMonth(), dateObj.getDate())
  const nowStart = new Date(now.getFullYear(), now.getMonth(), now.getDate())
  const diffInMs = nowStart.getTime() - dateStart.getTime()
  const diffInDays = Math.floor(diffInMs / (1000 * 60 * 60 * 24))
  
  if (diffInDays === 1) {
    return '1 dia atrás'
  } else if (diffInDays < 30) {
    return `${diffInDays} dias atrás`
  } else if (diffInDays < 365) {
    const months = Math.floor(diffInDays / 30)
    return months === 1 ? '1 mês atrás' : `${months} meses atrás`
  } else {
    const years = Math.floor(diffInDays / 365)
    return years === 1 ? '1 ano atrás' : `${years} anos atrás`
  }
}

// Get exact date string for tooltips
export const formatDateExact = (date: Date | string | null | undefined): string => {
  if (!date) return 'Data não informada'
  
  const dateObj = typeof date === 'string' ? new Date(date) : date
  return dateObj.toLocaleDateString('pt-BR', {
    weekday: 'long',
    day: '2-digit',
    month: 'long',
    year: 'numeric'
  })
}

// Specific function for last one-on-one (keeps backward compatibility)
export const formatLastOneOnOne = (lastOneOnOneDate?: Date): string => {
  return formatDateRelative(lastOneOnOneDate, { showNeverText: true })
}