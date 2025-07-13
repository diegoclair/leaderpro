'use client'

import React, { useState, useRef, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Calendar } from 'lucide-react'
import { Person } from '@/lib/types'

interface PersonModalProps {
  open: boolean
  onClose: () => void
  mode: 'create' | 'edit'
  person?: Person // Only needed for edit mode
  onSubmit: (personData: PersonFormData) => Promise<boolean>
  isLoading?: boolean
}

export interface PersonFormData {
  name: string
  email?: string
  position?: string
  department?: string
  phone?: string
  start_date?: string
  notes?: string
}

export default function PersonModal({
  open,
  onClose,
  mode,
  person,
  onSubmit,
  isLoading = false
}: PersonModalProps) {
  const [formData, setFormData] = useState<PersonFormData>({
    name: '',
    email: '',
    position: '',
    department: '',
    phone: '',
    start_date: '',
    notes: ''
  })
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [hasTypedName, setHasTypedName] = useState(false)
  const dateInputRef = useRef<HTMLInputElement>(null)

  // Load person data into form when modal opens (for edit mode)
  useEffect(() => {
    if (open) {
      if (mode === 'edit' && person) {
        setFormData({
          name: person.name || '',
          email: person.email || '',
          position: person.position || '',
          department: person.department || '',
          phone: person.phone || '',
          start_date: person.startDate ? formatDateForDisplay(person.startDate instanceof Date ? person.startDate.toISOString().split('T')[0] : person.startDate) : '',
          notes: person.notes || ''
        })
      } else {
        // Reset form for create mode
        setFormData({
          name: '',
          email: '',
          position: '',
          department: '',
          phone: '',
          start_date: '',
          notes: ''
        })
      }
      setHasTypedName(false)
    }
  }, [open, mode, person])

  // Check if date is valid (dd/mm/yyyy format)
  const isDateValid = (dateString: string): boolean => {
    if (!dateString) return true // Empty is valid (optional field)
    const cleanDate = dateString.replace(/\D/g, '')
    return cleanDate.length === 8 && formatDateForBackend(dateString) !== ''
  }

  const handleInputChange = (field: keyof PersonFormData, value: string) => {
    if (field === 'name') {
      setHasTypedName(true)
    }
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
  }

  // Convert date from YYYY-MM-DD (backend) to DD/MM/YYYY (Brazilian display)
  const formatDateForDisplay = (dateString: string): string => {
    if (!dateString) return ''
    const [year, month, day] = dateString.split('-')
    return `${day}/${month}/${year}`
  }

  // Convert date from DD/MM/YYYY (Brazilian input) to YYYY-MM-DD (backend)
  const formatDateForBackend = (dateString: string): string => {
    if (!dateString) return ''
    const cleanDate = dateString.replace(/\D/g, '')
    if (cleanDate.length !== 8) return '' // Invalid date
    
    const day = cleanDate.slice(0, 2)
    const month = cleanDate.slice(2, 4)
    const year = cleanDate.slice(4, 8)
    
    // Basic validation
    const dayNum = parseInt(day)
    const monthNum = parseInt(month)
    const yearNum = parseInt(year)
    
    if (dayNum < 1 || dayNum > 31 || monthNum < 1 || monthNum > 12 || yearNum < 1900) {
      return ''
    }
    
    return `${year}-${month}-${day}`
  }

  const handleDateChange = (value: string) => {
    // If user is typing manually, format as they type
    let formattedValue = value.replace(/\D/g, '') // Remove non-digits
    
    if (formattedValue.length >= 2) {
      formattedValue = formattedValue.slice(0, 2) + '/' + formattedValue.slice(2)
    }
    if (formattedValue.length >= 5) {
      formattedValue = formattedValue.slice(0, 5) + '/' + formattedValue.slice(5, 9)
    }
    
    // Limit to 10 characters (dd/mm/yyyy)
    if (formattedValue.length > 10) {
      formattedValue = formattedValue.slice(0, 10)
    }
    
    setFormData(prev => ({
      ...prev,
      start_date: formattedValue
    }))
  }

  const handlePhoneChange = (value: string) => {
    // Remove non-digits
    const digits = value.replace(/\D/g, '')
    
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
    
    setFormData(prev => ({
      ...prev,
      phone: formattedValue
    }))
  }

  // Handle calendar picker (YYYY-MM-DD format from date input)
  const handleCalendarChange = (value: string) => {
    if (value) {
      // Convert from YYYY-MM-DD to DD/MM/YYYY
      const [year, month, day] = value.split('-')
      const brDate = `${day}/${month}/${year}`
      setFormData(prev => ({
        ...prev,
        start_date: brDate
      }))
    } else {
      setFormData(prev => ({
        ...prev,
        start_date: ''
      }))
    }
  }

  // Open calendar picker
  const openCalendar = () => {
    if (dateInputRef.current) {
      dateInputRef.current.showPicker?.()
    }
  }

  // Convert Brazilian date to YYYY-MM-DD for the hidden date input
  const getDateInputValue = (): string => {
    if (!formData.start_date) return ''
    const backendDate = formatDateForBackend(formData.start_date)
    return backendDate || ''
  }

  const handleSubmit = async () => {
    if (!formData.name.trim() || isSubmitting) {
      return
    }

    setIsSubmitting(true)
    try {
      // Convert date to backend format before sending
      let startDateForBackend: string | undefined = undefined
      if (formData.start_date) {
        const backendDate = formatDateForBackend(formData.start_date)
        if (backendDate) {
          // Convert YYYY-MM-DD to RFC3339 format (ISO 8601) that Go expects
          startDateForBackend = `${backendDate}T00:00:00Z`
        }
      }
      
      const dataToSend = {
        ...formData,
        start_date: startDateForBackend
      }
      
      const success = await onSubmit(dataToSend)
      if (success) {
        if (mode === 'create') {
          // Reset form for create mode
          setFormData({
            name: '',
            email: '',
            position: '',
            department: '',
            phone: '',
            start_date: '',
            notes: ''
          })
        }
        onClose()
      }
    } catch (error) {
      console.error('Error submitting person:', error)
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleClose = () => {
    setHasTypedName(false)
    setIsSubmitting(false)
    onClose()
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey && formData.name.trim()) {
      e.preventDefault()
      handleSubmit()
    }
  }

  // Dynamic content based on mode
  const modalTitle = mode === 'create' ? 'Adicionar Nova Pessoa' : 'Editar Pessoa'
  const modalDescription = mode === 'create' 
    ? 'Adicione uma nova pessoa ao seu time. Apenas o nome é obrigatório.'
    : `Atualize as informações de ${person?.name}. Apenas o nome é obrigatório.`
  const submitButtonText = mode === 'create' ? 'Adicionar Pessoa' : 'Salvar Alterações'
  const loadingButtonText = mode === 'create' ? 'Adicionando...' : 'Salvando...'

  return (
    <Dialog open={open} onOpenChange={handleClose}>
      <DialogContent className="max-w-2xl max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>{modalTitle}</DialogTitle>
          <DialogDescription>
            {modalDescription}
          </DialogDescription>
        </DialogHeader>
        
        <div className="space-y-4" onKeyDown={handleKeyDown}>
          <div className="space-y-2">
            <Label htmlFor="name" className="text-sm font-medium">
              Nome completo <span className="text-red-500">*</span>
            </Label>
            <Input
              id="name"
              value={formData.name}
              onChange={(e) => handleInputChange('name', e.target.value)}
              placeholder="Ex: Maria Silva"
              required
              className={hasTypedName && !formData.name.trim() ? "border-red-200 focus:border-red-400" : ""}
            />
            {hasTypedName && !formData.name.trim() && (
              <p className="text-xs text-red-500">Nome é obrigatório</p>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="email" className="text-sm font-medium">Email</Label>
              <Input
                id="email"
                type="email"
                value={formData.email}
                onChange={(e) => handleInputChange('email', e.target.value)}
                placeholder="maria@empresa.com"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="phone" className="text-sm font-medium">Telefone</Label>
              <Input
                id="phone"
                type="tel"
                value={formData.phone}
                onChange={(e) => handlePhoneChange(e.target.value)}
                placeholder="(11) 99999-9999"
                maxLength={15}
              />
            </div>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="space-y-2">
              <Label htmlFor="position" className="text-sm font-medium">Cargo/Posição</Label>
              <Input
                id="position"
                value={formData.position}
                onChange={(e) => handleInputChange('position', e.target.value)}
                placeholder="Ex: Analista de Sistemas"
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="department" className="text-sm font-medium">Departamento/Squad</Label>
              <Input
                id="department"
                value={formData.department}
                onChange={(e) => handleInputChange('department', e.target.value)}
                placeholder="Ex: Tecnologia, Backend Squad"
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="start_date" className="text-sm font-medium text-muted-foreground">
              Data de início <span className="text-xs">(opcional)</span>
            </Label>
            <div className="relative">
              <Input
                id="start_date"
                type="text"
                value={formData.start_date}
                onChange={(e) => handleDateChange(e.target.value)}
                placeholder="dd/mm/aaaa"
                maxLength={10}
                className={`pr-10 ${formData.start_date && !isDateValid(formData.start_date) ? "border-red-200 focus:border-red-400" : ""}`}
              />
              <Button
                type="button"
                variant="ghost"
                size="sm"
                className="absolute right-0 top-0 h-full px-3 py-2 hover:bg-transparent"
                onClick={openCalendar}
              >
                <Calendar className="h-4 w-4 text-muted-foreground" />
              </Button>
              {/* Hidden date input for calendar picker */}
              <input
                ref={dateInputRef}
                type="date"
                value={getDateInputValue()}
                onChange={(e) => handleCalendarChange(e.target.value)}
                className="absolute inset-0 w-full h-full opacity-0 cursor-pointer"
                style={{ zIndex: -1 }}
              />
            </div>
            {formData.start_date && !isDateValid(formData.start_date) && (
              <p className="text-xs text-red-500">Data inválida. Use o formato dd/mm/aaaa</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="notes" className="text-sm font-medium text-muted-foreground">
              Notas <span className="text-xs">(opcional)</span>
            </Label>
            <Textarea
              id="notes"
              value={formData.notes}
              onChange={(e) => handleInputChange('notes', e.target.value)}
              placeholder="Informações relevantes sobre a pessoa..."
              rows={3}
              className="resize-none"
            />
          </div>
        </div>

        <DialogFooter className="gap-2">
          <Button 
            variant="outline" 
            onClick={handleClose} 
            disabled={isSubmitting || isLoading}
          >
            Cancelar
          </Button>
          <Button 
            onClick={handleSubmit} 
            disabled={!formData.name.trim() || isSubmitting || isLoading}
            className="min-w-[120px]"
          >
            {(isSubmitting || isLoading) ? loadingButtonText : submitButtonText}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}