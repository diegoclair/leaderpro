'use client'

import { forwardRef } from 'react'
import { Input } from '@/components/ui/input'
import { cn } from '@/lib/utils/cn'

interface PhoneInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'onChange'> {
  value?: string
  onChange?: (value: string) => void
  className?: string
}

export const PhoneInput = forwardRef<HTMLInputElement, PhoneInputProps>(
  ({ value = '', onChange, className, ...props }, ref) => {
    
    const formatPhoneNumber = (phoneValue: string): string => {
      // Remove all non-digits
      const digits = phoneValue.replace(/\D/g, '')
      
      // Limit to 11 digits (Brazilian mobile maximum)
      const limitedDigits = digits.slice(0, 11)
      
      let formattedValue = ''
      
      if (limitedDigits.length === 0) {
        formattedValue = ''
      } else if (limitedDigits.length <= 2) {
        // (XX
        formattedValue = `(${limitedDigits}`
      } else if (limitedDigits.length <= 6) {
        // (XX) XXXX
        formattedValue = `(${limitedDigits.slice(0, 2)}) ${limitedDigits.slice(2)}`
      } else if (limitedDigits.length <= 10) {
        // (XX) XXXX-XXXX (landline format)
        formattedValue = `(${limitedDigits.slice(0, 2)}) ${limitedDigits.slice(2, 6)}-${limitedDigits.slice(6)}`
      } else {
        // (XX) XXXXX-XXXX (mobile format with 11 digits)
        formattedValue = `(${limitedDigits.slice(0, 2)}) ${limitedDigits.slice(2, 7)}-${limitedDigits.slice(7)}`
      }
      
      return formattedValue
    }

    const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
      const rawValue = e.target.value
      const formattedValue = formatPhoneNumber(rawValue)
      
      if (onChange) {
        onChange(formattedValue)
      }
    }

    return (
      <Input
        {...props}
        ref={ref}
        type="tel"
        value={value}
        onChange={handleChange}
        placeholder="(11) 99999-9999"
        className={cn(className)}
        maxLength={15} // Maximum length for formatted phone
      />
    )
  }
)

PhoneInput.displayName = 'PhoneInput'

// Utility function to extract only digits from a formatted phone number
export const extractPhoneDigits = (formattedPhone: string): string => {
  return formattedPhone.replace(/\D/g, '')
}

// Utility function to validate Brazilian phone format
export const isValidBrazilianPhone = (phone: string): boolean => {
  const digits = extractPhoneDigits(phone)
  
  // Brazilian phone numbers:
  // - Landline: 10 digits (XX XXXX-XXXX)
  // - Mobile: 11 digits (XX XXXXX-XXXX)
  return digits.length === 10 || digits.length === 11
}