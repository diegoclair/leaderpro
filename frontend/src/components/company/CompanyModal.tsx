'use client'

import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Building2 } from 'lucide-react'

interface CompanyModalProps {
  open: boolean
  onClose: () => void
  onSubmit: (companyData: CompanyFormData) => Promise<boolean>
  isLoading?: boolean
}

export interface CompanyFormData {
  name: string
  industry: string
  size: string
  role: string
  is_default: boolean
}

const INDUSTRIES = [
  { value: 'technology', label: 'Tecnologia' },
  { value: 'finance', label: 'Financeiro' },
  { value: 'healthcare', label: 'Saúde' },
  { value: 'education', label: 'Educação' },
  { value: 'retail', label: 'Varejo' },
  { value: 'manufacturing', label: 'Manufatura' },
  { value: 'consulting', label: 'Consultoria' },
  { value: 'marketing', label: 'Marketing' },
  { value: 'real_estate', label: 'Imóveis' },
  { value: 'other', label: 'Outro' }
]

const COMPANY_SIZES = [
  { value: 'small', label: 'Pequena (2-10 pessoas)' },
  { value: 'medium', label: 'Média (11-50 pessoas)' },
  { value: 'large', label: 'Grande (51-200 pessoas)' },
  { value: 'enterprise', label: 'Corporação (200+ pessoas)' }
]

export default function CompanyModal({
  open,
  onClose,
  onSubmit,
  isLoading = false
}: CompanyModalProps) {
  const [formData, setFormData] = useState<CompanyFormData>({
    name: '',
    industry: '',
    size: '',
    role: '',
    is_default: false
  })
  const [isSubmitting, setIsSubmitting] = useState(false)

  // Reset form when modal opens
  useEffect(() => {
    if (open) {
      setFormData({
        name: '',
        industry: '',
        size: '',
        role: '',
        is_default: false
      })
      setIsSubmitting(false)
    }
  }, [open])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Validation
    if (!formData.name.trim()) {
      return
    }

    setIsSubmitting(true)
    
    try {
      const success = await onSubmit(formData)
      if (success) {
        onClose()
      }
    } finally {
      setIsSubmitting(false)
    }
  }

  const handleInputChange = (field: keyof CompanyFormData, value: string | boolean) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <Building2 className="h-5 w-5" />
            Adicionar Nova Empresa
          </DialogTitle>
          <DialogDescription>
            Crie uma nova empresa para gerenciar seu time e realizar 1:1s.
          </DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="space-y-4">
          {/* Nome da Empresa */}
          <div className="space-y-2">
            <Label htmlFor="company-name">
              Nome da Empresa <span className="text-red-500">*</span>
            </Label>
            <Input
              id="company-name"
              value={formData.name}
              onChange={(e) => handleInputChange('name', e.target.value)}
              placeholder="Ex: TechCorp Brasil"
              required
              disabled={isSubmitting || isLoading}
            />
          </div>

          {/* Indústria */}
          <div className="space-y-2">
            <Label htmlFor="company-industry">Setor/Indústria</Label>
            <Select
              value={formData.industry}
              onValueChange={(value) => handleInputChange('industry', value)}
              disabled={isSubmitting || isLoading}
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione o setor" />
              </SelectTrigger>
              <SelectContent>
                {INDUSTRIES.map((industry) => (
                  <SelectItem key={industry.value} value={industry.value}>
                    {industry.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Tamanho da Empresa */}
          <div className="space-y-2">
            <Label htmlFor="company-size">Tamanho da Empresa</Label>
            <Select
              value={formData.size}
              onValueChange={(value) => handleInputChange('size', value)}
              disabled={isSubmitting || isLoading}
            >
              <SelectTrigger>
                <SelectValue placeholder="Selecione o tamanho" />
              </SelectTrigger>
              <SelectContent>
                {COMPANY_SIZES.map((size) => (
                  <SelectItem key={size.value} value={size.value}>
                    {size.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>

          {/* Seu Papel na Empresa */}
          <div className="space-y-2">
            <Label htmlFor="your-role">Seu Papel na Empresa</Label>
            <Input
              id="your-role"
              value={formData.role}
              onChange={(e) => handleInputChange('role', e.target.value)}
              placeholder="Ex: CTO, Gerente de Projetos, Team Lead"
              disabled={isSubmitting || isLoading}
            />
          </div>

          <DialogFooter className="gap-2">
            <Button
              type="button"
              variant="outline"
              onClick={onClose}
              disabled={isSubmitting || isLoading}
            >
              Cancelar
            </Button>
            <Button
              type="submit"
              disabled={!formData.name.trim() || isSubmitting || isLoading}
              className="min-w-[100px]"
            >
              {isSubmitting || isLoading ? (
                <div className="flex items-center gap-2">
                  <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
                  Criando...
                </div>
              ) : (
                'Criar Empresa'
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}