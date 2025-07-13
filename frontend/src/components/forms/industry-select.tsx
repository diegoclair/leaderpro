import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
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
  return (
    <Select value={value} onValueChange={onValueChange} disabled={disabled}>
      <SelectTrigger className={className}>
        <SelectValue placeholder={placeholder} />
      </SelectTrigger>
      <SelectContent>
        {INDUSTRIES.map((industry) => (
          <SelectItem key={industry.value} value={industry.value}>
            {industry.label}
          </SelectItem>
        ))}
      </SelectContent>
    </Select>
  )
}