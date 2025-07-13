/**
 * Gender utility functions for Portuguese language
 */

export type Gender = 'male' | 'female' | 'other' | undefined

/**
 * Returns the correct form of "mencionado/mencionada" based on gender
 */
export function getMentionedText(gender: Gender): string {
  switch (gender) {
    case 'female':
      return 'mencionada'
    case 'male':
      return 'mencionado'
    default:
      return 'mencionado(a)' // Use both forms for 'other' or undefined
  }
}

/**
 * Returns the correct form of "foi mencionado/foi mencionada" based on gender
 */
export function getWasMentionedText(gender: Gender): string {
  switch (gender) {
    case 'female':
      return 'foi mencionada'
    case 'male':
      return 'foi mencionado'
    default:
      return 'foi mencionado(a)' // Use both forms for 'other' or undefined
  }
}

/**
 * Returns a full sentence like "Sabrina foi mencionada" based on name and gender
 */
export function getMentionSentence(name: string, gender: Gender): string {
  return `${name} ${getWasMentionedText(gender)}`
}

/**
 * Returns the correct article (o/a) based on gender
 */
export function getArticle(gender: Gender): string {
  switch (gender) {
    case 'female':
      return 'a'
    case 'male':
      return 'o'
    default:
      return 'o(a)' // Use both forms for 'other' or undefined
  }
}