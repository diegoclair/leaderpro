/**
 * Storage Manager - Centralized localStorage management for LeaderPro
 * 
 * All localStorage operations should go through this manager to ensure:
 * - Consistent key naming
 * - Easy cleanup on logout
 * - Better debugging and tracking
 */

export type StorageKey = 
  | 'auth-storage'           // Zustand auth persist storage
  | 'leaderpro-active-company'  // Currently selected company
  | 'leaderpro-theme'        // User theme preference
  | 'leaderpro-onboarding'   // Onboarding completion status

class StorageManager {
  private readonly prefix = 'leaderpro-'
  
  // Lista de todas as chaves gerenciadas pelo app
  private readonly managedKeys: StorageKey[] = [
    'auth-storage',
    'leaderpro-active-company',
    'leaderpro-theme',
    'leaderpro-onboarding'
  ]

  /**
   * Get item from localStorage
   */
  get<T = string>(key: StorageKey): T | null {
    if (typeof window === 'undefined') return null
    
    try {
      const value = localStorage.getItem(key)
      if (value === null) return null
      
      // Try to parse as JSON, fallback to string
      try {
        return JSON.parse(value) as T
      } catch {
        return value as unknown as T
      }
    } catch (error) {
      console.error(`❌ Error reading from localStorage key "${key}":`, error)
      return null
    }
  }

  /**
   * Set item in localStorage
   */
  set<T>(key: StorageKey, value: T): void {
    if (typeof window === 'undefined') return
    
    try {
      const serializedValue = typeof value === 'string' ? value : JSON.stringify(value)
      localStorage.setItem(key, serializedValue)
      console.log(`📦 Stored in localStorage: ${key}`)
    } catch (error) {
      console.error(`❌ Error storing to localStorage key "${key}":`, error)
    }
  }

  /**
   * Remove specific item from localStorage
   */
  remove(key: StorageKey): void {
    if (typeof window === 'undefined') return
    
    try {
      localStorage.removeItem(key)
      console.log(`🗑️ Removed from localStorage: ${key}`)
    } catch (error) {
      console.error(`❌ Error removing localStorage key "${key}":`, error)
    }
  }

  /**
   * Clear ALL LeaderPro data from localStorage
   * This is called on logout to ensure no user data persists
   */
  clearAll(): void {
    if (typeof window === 'undefined') return
    
    console.log('🧹 Clearing all LeaderPro data from localStorage...')
    
    try {
      // Remove all managed keys
      this.managedKeys.forEach(key => {
        if (localStorage.getItem(key) !== null) {
          localStorage.removeItem(key)
          console.log(`🗑️ Cleared: ${key}`)
        }
      })
      
      // Also scan for any orphaned keys that start with our prefix
      const orphanedKeys: string[] = []
      for (let i = 0; i < localStorage.length; i++) {
        const key = localStorage.key(i)
        if (key && key.startsWith(this.prefix) && !this.managedKeys.includes(key as StorageKey)) {
          orphanedKeys.push(key)
        }
      }
      
      // Remove orphaned keys
      orphanedKeys.forEach(key => {
        localStorage.removeItem(key)
        console.log(`🗑️ Cleared orphaned key: ${key}`)
      })
      
      console.log('✅ All LeaderPro data cleared from localStorage')
      
    } catch (error) {
      console.error('❌ Error clearing localStorage:', error)
    }
  }

  /**
   * Check if a key exists in localStorage
   */
  has(key: StorageKey): boolean {
    if (typeof window === 'undefined') return false
    return localStorage.getItem(key) !== null
  }

  /**
   * Get all managed keys and their values (for debugging)
   */
  debug(): Record<string, any> {
    if (typeof window === 'undefined') return {}
    
    const data: Record<string, any> = {}
    this.managedKeys.forEach(key => {
      data[key] = this.get(key)
    })
    
    console.table(data)
    return data
  }
}

// Export singleton instance
export const storageManager = new StorageManager()

// Convenience functions for common operations
export const getStoredData = <T>(key: StorageKey): T | null => storageManager.get<T>(key)
export const setStoredData = <T>(key: StorageKey, value: T): void => storageManager.set(key, value)
export const removeStoredData = (key: StorageKey): void => storageManager.remove(key)
export const clearAllStoredData = (): void => storageManager.clearAll()
export const hasStoredData = (key: StorageKey): boolean => storageManager.has(key)