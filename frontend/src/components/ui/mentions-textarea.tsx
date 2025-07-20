'use client'

import React, { useState, useRef, useEffect } from 'react'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Person } from '@/lib/types'
import { getInitials } from '@/lib/utils/names'
import { cn } from '@/lib/utils/cn'

interface MentionsTextareaProps {
  value: string
  onChange: (value: string) => void
  people: Person[]
  placeholder?: string
  className?: string
  disabled?: boolean
  minHeight?: number
}

interface MentionSuggestion {
  person: Person
  startIndex: number
  query: string
  textNode: Node
  endIndex: number
}

export function MentionsTextarea({
  value,
  onChange,
  people,
  placeholder = 'Digite sua anotação... Use @nome para mencionar pessoas da equipe',
  className,
  disabled = false,
  minHeight = 100
}: MentionsTextareaProps) {
  const [showSuggestions, setShowSuggestions] = useState(false)
  const [mentionSuggestion, setMentionSuggestion] = useState<MentionSuggestion | null>(null)
  const [filteredPeople, setFilteredPeople] = useState<Person[]>([])
  const [selectedIndex, setSelectedIndex] = useState(0)
  const editorRef = useRef<HTMLDivElement>(null)
  const [isTyping, setIsTyping] = useState(false)
  const [lastInputType, setLastInputType] = useState<string | null>(null)
  const [lastExternalValue, setLastExternalValue] = useState('')

  // Convert backend format to display HTML
  const convertToDisplayHTML = (backendValue: string): string => {
    const mentionRegex = /\{\{person:([^|]+)\|([^}]+)\}\}/g
    return backendValue.replace(mentionRegex, (match, uuid, name) => {
      return `<span class="mention" data-uuid="${uuid}" data-name="${name}" contenteditable="false">@${name}</span>`
    })
  }

  // Convert display HTML back to backend format
  const convertToBackendValue = (html: string): string => {
    const div = document.createElement('div')
    div.innerHTML = html
    
    const mentions = div.querySelectorAll('.mention')
    mentions.forEach(mention => {
      const uuid = mention.getAttribute('data-uuid')
      const name = mention.getAttribute('data-name')
      const backendFormat = `{{person:${uuid}|${name}}}`
      mention.replaceWith(document.createTextNode(backendFormat))
    })
    
    return div.textContent || ''
  }

  // Initialize editor content when component mounts or value changes externally  
  useEffect(() => {
    if (editorRef.current && !isTyping) {
      const currentText = editorRef.current.textContent || ''
      const expectedContent = convertToDisplayHTML(value)
      
      // Only update if this is a genuine external change
      const isExternalChange = value !== lastExternalValue
      const isInitialLoad = !currentText && !lastExternalValue
      const isClearingForm = !value && currentText && lastExternalValue // Form was cleared externally
      
      if (isInitialLoad || (isExternalChange && isClearingForm)) {
        editorRef.current.innerHTML = expectedContent
        setLastExternalValue(value)
      } else if (isExternalChange && !currentText) {
        // Only update if editor is empty but we have new value
        editorRef.current.innerHTML = expectedContent
        setLastExternalValue(value)
      }
    }
  }, [value, isTyping, lastExternalValue])

  // Detecta @ e extrai query
  const detectMention = (node: Node, offset: number): MentionSuggestion | null => {
    if (node.nodeType !== Node.TEXT_NODE) return null
    
    const text = node.textContent || ''
    const beforeCursor = text.substring(0, offset)
    const lastAtIndex = beforeCursor.lastIndexOf('@')
    
    if (lastAtIndex === -1) return null
    
    // Verifica se há espaço entre @ e cursor
    const afterAt = beforeCursor.substring(lastAtIndex + 1)
    if (afterAt.includes(' ') || afterAt.includes('\n')) return null
    
    // Verifica se @ está no início ou após espaço
    const charBeforeAt = lastAtIndex > 0 ? beforeCursor[lastAtIndex - 1] : ' '
    if (charBeforeAt !== ' ' && charBeforeAt !== '\n') return null
    
    return {
      person: people[0], // Placeholder
      startIndex: lastAtIndex,
      query: afterAt.toLowerCase(),
      textNode: document.createTextNode(''), // Placeholder
      endIndex: lastAtIndex + afterAt.length + 1
    }
  }

  // Filtra pessoas baseado na query
  const filterPeople = (query: string): Person[] => {
    if (!query) return people.slice(0, 5)
    
    return people
      .filter(person => 
        person.name.toLowerCase().includes(query) ||
        (person.position && person.position.toLowerCase().includes(query))
      )
      .slice(0, 5)
  }

  // Handler para mudanças no editor
  const handleInput = (e: React.FormEvent<HTMLDivElement>) => {
    if (!editorRef.current) return
    
    setIsTyping(true)
    
    // Captura o tipo de input (insertText, insertCompositionText, etc)
    const nativeEvent = e.nativeEvent as InputEvent
    const inputType = nativeEvent?.inputType
    setLastInputType(inputType || null)
    
    const selection = window.getSelection()
    if (!selection || selection.rangeCount === 0) return
    
    const range = selection.getRangeAt(0)
    const { startContainer, startOffset } = range
    
    // Detecta mention APENAS se foi digitação real de texto (não clique, não paste, etc)
    const isRealTyping = inputType === 'insertText' || inputType === 'insertCompositionText'
    const mention = detectMention(startContainer, startOffset)
    
    if (mention && isRealTyping) {
      // Verifica se o último caractere digitado foi @ ou se está completando uma mention
      const text = startContainer.textContent || ''
      const lastChar = text[startOffset - 1]
      
      if (lastChar === '@' || mention.query.length > 0) {
        const filtered = filterPeople(mention.query)
        setMentionSuggestion({
          ...mention, 
          startIndex: mention.startIndex,
          textNode: startContainer,
          endIndex: startOffset
        })
        setFilteredPeople(filtered)
        setShowSuggestions(filtered.length > 0)
        setSelectedIndex(0)
      } else {
        setShowSuggestions(false)
        setMentionSuggestion(null)
      }
    } else {
      setShowSuggestions(false)
      setMentionSuggestion(null)
    }
    
    // Don't update backend immediately during typing to avoid cursor resets
    setTimeout(() => {
      setIsTyping(false)
      // Update backend value only after typing stops
      if (editorRef.current) {
        const backendValue = convertToBackendValue(editorRef.current.innerHTML)
        setLastExternalValue(backendValue) // Track that this change came from us
        onChange(backendValue)
      }
    }, 300)
  }

  // Handler para seleção de pessoa
  const handlePersonSelect = (person: Person) => {
    if (!mentionSuggestion || !editorRef.current) return
    
    const { textNode, startIndex, endIndex } = mentionSuggestion
    
    if (textNode.nodeType !== Node.TEXT_NODE) return
    
    const text = textNode.textContent || ''
    const beforeMention = text.substring(0, startIndex)
    const afterCursor = text.substring(endIndex)
    
    // Create mention element
    const mentionElement = document.createElement('span')
    mentionElement.className = 'mention bg-blue-100 dark:bg-blue-900/30 text-blue-700 dark:text-blue-300 px-1 py-0.5 rounded font-medium'
    mentionElement.setAttribute('data-uuid', person.uuid || person.id.toString())
    mentionElement.setAttribute('data-name', person.name)
    mentionElement.contentEditable = 'false'
    mentionElement.textContent = `@${person.name}`
    
    // Create nodes with non-breaking space for better visibility
    const beforeNode = beforeMention ? document.createTextNode(beforeMention) : null
    const spaceNode = document.createTextNode('\u00A0') // Non-breaking space
    const afterNode = afterCursor ? document.createTextNode(afterCursor) : null
    
    const parent = textNode.parentNode
    if (parent) {
      // Insert all elements
      if (beforeNode) {
        parent.insertBefore(beforeNode, textNode)
      }
      parent.insertBefore(mentionElement, textNode)
      parent.insertBefore(spaceNode, textNode)
      if (afterNode) {
        parent.insertBefore(afterNode, textNode)
      }
      
      // Remove original text node
      parent.removeChild(textNode)
      
      // Focus editor and position cursor after space
      editorRef.current.focus()
      
      setTimeout(() => {
        const selection = window.getSelection()
        if (selection) {
          const range = document.createRange()
          range.setStartAfter(spaceNode)
          range.collapse(true)
          selection.removeAllRanges()
          selection.addRange(range)
        }
      }, 0)
    }
    
    setShowSuggestions(false)
    setMentionSuggestion(null)
    
    // Update backend value immediately for mentions
    setTimeout(() => {
      if (editorRef.current) {
        const backendValue = convertToBackendValue(editorRef.current.innerHTML)
        setLastExternalValue(backendValue) // Track that this change came from us
        onChange(backendValue)
      }
    }, 0)
  }

  // Handler para teclas
  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (showSuggestions && filteredPeople.length > 0) {
      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault()
          setSelectedIndex(prev => Math.min(prev + 1, filteredPeople.length - 1))
          break
        case 'ArrowUp':
          e.preventDefault()
          setSelectedIndex(prev => Math.max(prev - 1, 0))
          break
        case 'Enter':
        case 'Tab':
          e.preventDefault()
          handlePersonSelect(filteredPeople[selectedIndex])
          break
        case 'Escape':
          setShowSuggestions(false)
          setMentionSuggestion(null)
          break
      }
      return
    }

    // Handle mention deletion
    if (e.key === 'Backspace') {
      const selection = window.getSelection()
      if (selection && selection.rangeCount > 0) {
        const range = selection.getRangeAt(0)
        const { startContainer, startOffset } = range
        
        if (startOffset === 0 && startContainer.previousSibling) {
          const prevSibling = startContainer.previousSibling
          if (prevSibling.nodeType === Node.ELEMENT_NODE && 
              (prevSibling as Element).classList.contains('mention')) {
            e.preventDefault()
            
            // Store current position reference
            const nextSibling = startContainer
            
            // Remove mention
            prevSibling.remove()
            
            // Maintain cursor position after the removal
            const newRange = document.createRange()
            newRange.setStart(nextSibling, 0)
            newRange.collapse(true)
            selection.removeAllRanges()
            selection.addRange(newRange)
            
            // Update backend value immediately for mention deletion
            if (editorRef.current) {
              const backendValue = convertToBackendValue(editorRef.current.innerHTML)
              setLastExternalValue(backendValue) // Track that this change came from us
              onChange(backendValue)
            }
          }
        }
      }
    }
  }


  // Handler para clique fora
  useEffect(() => {
    const handleClickOutside = () => {
      setShowSuggestions(false)
      setMentionSuggestion(null)
    }
    
    if (showSuggestions) {
      document.addEventListener('click', handleClickOutside)
      return () => document.removeEventListener('click', handleClickOutside)
    }
  }, [showSuggestions])

  return (
    <div className={cn('relative', className)}>
      <div className="relative">
        <div
          ref={editorRef}
          contentEditable={!disabled}
          onInput={handleInput}
          onKeyDown={handleKeyDown}
          className={cn(
            'min-h-[100px] w-full rounded-md border border-input bg-background px-3 py-2 text-sm',
            'focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2',
            'disabled:cursor-not-allowed disabled:opacity-50',
            '[&_.mention]:bg-blue-100 [&_.mention]:dark:bg-blue-900/30',
            '[&_.mention]:text-blue-700 [&_.mention]:dark:text-blue-300',
            '[&_.mention]:px-1 [&_.mention]:py-0.5 [&_.mention]:rounded [&_.mention]:font-medium',
            className
          )}
          style={{ minHeight: `${minHeight}px` }}
          suppressContentEditableWarning={true}
          data-placeholder={placeholder}
        />
        
        {/* Placeholder quando vazio */}
        <div 
          className={cn(
            'absolute left-3 top-2 text-sm text-muted-foreground pointer-events-none',
            editorRef.current?.textContent ? 'hidden' : 'block'
          )}
        >
          {placeholder}
        </div>
        
        {/* Suggestions dropdown */}
        {showSuggestions && (
          <div className="absolute top-full left-0 mt-1 w-80 p-2 bg-popover border border-border rounded-md shadow-md z-50">
            <div className="max-h-64 overflow-y-auto">
              {filteredPeople.length === 0 ? (
                <div className="py-6 text-center text-sm text-muted-foreground">
                  Nenhuma pessoa encontrada
                </div>
              ) : (
                <div className="space-y-1">
                  {filteredPeople.map((person, index) => (
                    <div
                      key={person.id}
                      onClick={(e) => {
                        e.preventDefault()
                        e.stopPropagation()
                        handlePersonSelect(person)
                      }}
                      className={cn(
                        'flex items-center gap-3 px-3 py-3 cursor-pointer rounded-md transition-colors',
                        index === selectedIndex 
                          ? 'bg-accent text-accent-foreground' 
                          : 'hover:bg-accent/50'
                      )}
                    >
                      <Avatar className="h-8 w-8 flex-shrink-0">
                        <AvatarImage src={person.avatar} alt={person.name} />
                        <AvatarFallback className="bg-blue-100 dark:bg-blue-900 text-blue-700 dark:text-blue-300 text-xs">
                          {getInitials(person.name)}
                        </AvatarFallback>
                      </Avatar>
                      
                      <div className="flex-1 min-w-0">
                        <div className="font-medium text-sm truncate">
                          {person.name}
                        </div>
                        {person.position && (
                          <div className="text-xs text-muted-foreground truncate mt-0.5">
                            {person.position}
                          </div>
                        )}
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        )}
      </div>
    </div>
  )
}