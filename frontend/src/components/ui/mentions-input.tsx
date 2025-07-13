'use client'

import React from 'react'
import { MentionsInput, Mention } from 'react-mentions'
import { Person } from '@/lib/types'
import { cn } from '@/lib/utils/cn'

interface MentionsInputComponentProps {
  value: string
  onChange: (value: string) => void
  people: Person[]
  placeholder?: string
  className?: string
  disabled?: boolean
  minHeight?: number
}

export function MentionsInputComponent({
  value,
  onChange,
  people,
  placeholder = 'Digite sua anotação... Use @nome para mencionar pessoas da equipe',
  className,
  disabled = false,
  minHeight = 100
}: MentionsInputComponentProps) {
  
  // Convert people to mentions format
  const mentionsData = people.map(person => ({
    id: person.id,
    display: person.name
  }))

  // Usar styling mais simples e compatível com shadcn/ui
  const mentionsInputStyle = {
    control: {
      backgroundColor: 'hsl(var(--background))',
      fontSize: 14,
      fontWeight: 'normal',
      minHeight: `${minHeight}px`,
      border: '1px solid hsl(var(--border))',
      borderRadius: 'calc(var(--radius) - 2px)',
    },
    
    '&multiLine': {
      control: {
        fontFamily: 'inherit',
        minHeight: `${minHeight}px`,
        border: '1px solid hsl(var(--border))',
        borderRadius: 'calc(var(--radius) - 2px)',
        padding: '12px',
        backgroundColor: 'hsl(var(--background))',
        fontSize: '14px',
      },
      highlighter: {
        padding: '12px',
        border: '1px solid transparent',
        borderRadius: 'calc(var(--radius) - 2px)',
        backgroundColor: 'transparent',
      },
      input: {
        padding: '12px',
        border: 'none',
        borderRadius: 'calc(var(--radius) - 2px)',
        outline: 'none',
        backgroundColor: 'transparent',
        color: 'hsl(var(--foreground))',
        fontSize: '14px',
        fontFamily: 'inherit',
        lineHeight: '1.5',
        resize: 'none' as const,
      },
    },

    suggestions: {
      list: {
        backgroundColor: 'hsl(var(--popover))',
        border: '1px solid hsl(var(--border))',
        borderRadius: 'calc(var(--radius) - 2px)',
        fontSize: 14,
        boxShadow: '0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)',
        maxHeight: '200px',
        overflowY: 'auto',
        zIndex: 50,
      },
      item: {
        padding: '8px 12px',
        borderBottom: '1px solid hsl(var(--border))',
        color: 'hsl(var(--foreground))',
        cursor: 'pointer',
        '&focused': {
          backgroundColor: 'hsl(var(--accent))',
          color: 'hsl(var(--accent-foreground))',
        },
      },
    },
  }

  const handleChange = (event: { target: { value: string } }) => {
    onChange(event.target.value)
  }

  return (
    <div className={cn('relative', className)}>
      <MentionsInput
        value={value}
        onChange={handleChange}
        placeholder={placeholder}
        disabled={disabled}
        style={mentionsInputStyle}
        allowSpaceInQuery
        className="mentions-input"
      >
        <Mention
          trigger="@"
          data={mentionsData}
          markup="{{person:__id__|__display__}}"
          displayTransform={(id: string, display: string) => `@${display}`}
          renderSuggestion={(suggestion, search, highlightedDisplay, index, focused) => (
            <div 
              className={cn(
                'px-3 py-2 cursor-pointer border-b border-border',
                focused && 'bg-accent text-accent-foreground'
              )}
            >
              <div className="font-medium">{suggestion.display}</div>
            </div>
          )}
          style={{
            backgroundColor: 'hsl(var(--primary) / 0.1)',
            color: 'hsl(var(--primary))',
            fontWeight: '500',
            borderRadius: '4px',
            padding: '2px 4px',
          }}
        />
      </MentionsInput>
    </div>
  )
}

// CSS personalizado para melhor integração com shadcn/ui
export const mentionsInputCSS = `
.mentions-input {
  width: 100%;
}

.mentions-input .mentions__control {
  min-height: 100px !important;
}

.mentions-input .mentions__input {
  font-family: inherit !important;
  font-size: 14px !important;
  line-height: 1.5 !important;
  color: hsl(var(--foreground)) !important;
}

.mentions-input .mentions__input:focus {
  outline: none !important;
  border-color: hsl(var(--ring)) !important;
  box-shadow: 0 0 0 2px hsl(var(--ring) / 0.2) !important;
}

.mentions-input .mentions__input:disabled {
  cursor: not-allowed !important;
  opacity: 0.5 !important;
}

.mentions-input .mentions__suggestions__list {
  background: hsl(var(--popover)) !important;
  border: 1px solid hsl(var(--border)) !important;
  border-radius: calc(var(--radius) - 2px) !important;
  box-shadow: 0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1) !important;
  z-index: 50 !important;
}

.mentions-input .mentions__suggestion {
  background: transparent !important;
  color: hsl(var(--foreground)) !important;
  border-bottom: 1px solid hsl(var(--border)) !important;
  padding: 8px 12px !important;
}

.mentions-input .mentions__suggestion--focused {
  background: hsl(var(--accent)) !important;
  color: hsl(var(--accent-foreground)) !important;
}

.mentions-input .mentions__mention {
  background: hsl(var(--primary) / 0.1) !important;
  color: hsl(var(--primary)) !important;
  font-weight: 500 !important;
  border-radius: 4px !important;
  padding: 2px 4px !important;
}

/* Dark mode adjustments */
.dark .mentions-input .mentions__input {
  color: hsl(var(--foreground)) !important;
}

.dark .mentions-input .mentions__suggestions__list {
  background: hsl(var(--popover)) !important;
  border-color: hsl(var(--border)) !important;
}

.dark .mentions-input .mentions__suggestion {
  color: hsl(var(--foreground)) !important;
  border-bottom-color: hsl(var(--border)) !important;
}

.dark .mentions-input .mentions__suggestion--focused {
  background: hsl(var(--accent)) !important;
  color: hsl(var(--accent-foreground)) !important;
}
`