'use client'

import { useState } from 'react'
import { Check, ChevronsUpDown, Search } from 'lucide-react'
import { Button } from '@/components/ui/button'
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { Input } from '@/components/ui/input'
import { ScrollArea } from '@/components/ui/scroll-area'
import { INDUSTRIES } from '@/lib/constants/company'

interface IndustrySelectProps {
  value?: string
  onValueChange: (value: string) => void
  placeholder?: string
  disabled?: boolean
  className?: string
}

export function IndustrySelect({
  value,
  onValueChange,
  placeholder = 'Selecione o setor',
  disabled = false,
  className
}: IndustrySelectProps) {
  const [open, setOpen] = useState(false)
  const [searchTerm, setSearchTerm] = useState('')

  // Filter industries based on search term
  const filteredIndustries = INDUSTRIES.filter(industry =>
    industry.label.toLowerCase().includes(searchTerm.toLowerCase())
  )

  // Get selected industry label
  const selectedIndustry = INDUSTRIES.find(industry => industry.value === value)

  const handleSelect = (industryValue: string) => {
    onValueChange(industryValue)
    setOpen(false)
    setSearchTerm('') // Clear search when selecting
  }

  return (
    <Popover open={open} onOpenChange={setOpen}>
      <PopoverTrigger asChild>
        <Button
          variant="outline"
          role="combobox"
          aria-expanded={open}
          className={`w-full justify-between ${className || ''}`}
          disabled={disabled}
        >
          <span className="truncate text-left">
            {selectedIndustry ? selectedIndustry.label : placeholder}
          </span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-full p-0" align="start">
        <div className="flex items-center border-b px-3">
          <Search className="mr-2 h-4 w-4 shrink-0 opacity-50" />
          <Input
            placeholder="Buscar setor..."
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            className="border-0 bg-transparent px-0 py-3 focus-visible:ring-0 focus-visible:ring-offset-0"
          />
        </div>
        <ScrollArea className="h-72">
          <div className="p-1">
            {filteredIndustries.length === 0 ? (
              <div className="px-2 py-3 text-sm text-muted-foreground">
                Nenhum setor encontrado.
              </div>
            ) : (
              filteredIndustries.map((industry) => (
                <div
                  key={industry.value}
                  className="flex cursor-pointer items-center rounded-sm px-2 py-2 text-sm hover:bg-accent hover:text-accent-foreground"
                  onClick={() => handleSelect(industry.value)}
                >
                  <Check
                    className={`mr-2 h-4 w-4 ${
                      value === industry.value ? 'opacity-100' : 'opacity-0'
                    }`}
                  />
                  <span className="flex-1">{industry.label}</span>
                </div>
              ))
            )}
          </div>
        </ScrollArea>
      </PopoverContent>
    </Popover>
  )
}