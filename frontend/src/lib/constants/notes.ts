/**
 * Note types and mappings for consistent usage across the application
 */

// Note source types (what kind of note it is)
export const NOTE_SOURCE_TYPES = {
  ONE_ON_ONE: 'one_on_one',
  FEEDBACK: 'feedback',
  OBSERVATION: 'observation'
} as const

// Note types (from timeline API)
export const NOTE_TYPES = {
  NOTE: 'note',
  MENTION: 'mention'
} as const

// Feedback types
export const FEEDBACK_TYPES = {
  POSITIVE: 'positive',
  CONSTRUCTIVE: 'constructive',
  NEUTRAL: 'neutral'
} as const

// Feedback categories
export const FEEDBACK_CATEGORIES = {
  PERFORMANCE: 'performance',
  BEHAVIOR: 'behavior',
  SKILL: 'skill',
  COLLABORATION: 'collaboration'
} as const

// Type definitions for TypeScript
export type NoteSourceType = typeof NOTE_SOURCE_TYPES[keyof typeof NOTE_SOURCE_TYPES]
export type NoteType = typeof NOTE_TYPES[keyof typeof NOTE_TYPES]
export type FeedbackType = typeof FEEDBACK_TYPES[keyof typeof FEEDBACK_TYPES]
export type FeedbackCategory = typeof FEEDBACK_CATEGORIES[keyof typeof FEEDBACK_CATEGORIES]

// Labels for display
export const NOTE_SOURCE_TYPE_LABELS: Record<NoteSourceType, string> = {
  [NOTE_SOURCE_TYPES.ONE_ON_ONE]: 'Reunião 1:1',
  [NOTE_SOURCE_TYPES.FEEDBACK]: 'Feedback',
  [NOTE_SOURCE_TYPES.OBSERVATION]: 'Observação'
}

// Feedback type labels
export const FEEDBACK_TYPE_LABELS: Record<FeedbackType, string> = {
  [FEEDBACK_TYPES.POSITIVE]: 'Positivo',
  [FEEDBACK_TYPES.CONSTRUCTIVE]: 'Construtivo',
  [FEEDBACK_TYPES.NEUTRAL]: 'Neutro'
}

// Feedback category labels
export const FEEDBACK_CATEGORY_LABELS: Record<FeedbackCategory, string> = {
  [FEEDBACK_CATEGORIES.PERFORMANCE]: 'Performance',
  [FEEDBACK_CATEGORIES.BEHAVIOR]: 'Comportamento',
  [FEEDBACK_CATEGORIES.SKILL]: 'Habilidade',
  [FEEDBACK_CATEGORIES.COLLABORATION]: 'Colaboração'
}

// CSS classes for feedback types
export const FEEDBACK_TYPE_COLORS: Record<FeedbackType, string> = {
  [FEEDBACK_TYPES.POSITIVE]: 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-300',
  [FEEDBACK_TYPES.CONSTRUCTIVE]: 'bg-orange-100 text-orange-800 dark:bg-orange-900 dark:text-orange-300',
  [FEEDBACK_TYPES.NEUTRAL]: 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-300'
}

// Helper functions
export function getNoteSourceTypeLabel(sourceType: string): string {
  return NOTE_SOURCE_TYPE_LABELS[sourceType as NoteSourceType] || 'Anotação'
}

export function getFeedbackTypeLabel(feedbackType: string): string {
  return FEEDBACK_TYPE_LABELS[feedbackType as FeedbackType] || feedbackType
}

export function getFeedbackCategoryLabel(category: string): string {
  return FEEDBACK_CATEGORY_LABELS[category as FeedbackCategory] || category
}

export function getFeedbackTypeColor(feedbackType: string): string {
  return FEEDBACK_TYPE_COLORS[feedbackType as FeedbackType] || 'bg-gray-100 text-gray-800 dark:bg-gray-900 dark:text-gray-300'
}