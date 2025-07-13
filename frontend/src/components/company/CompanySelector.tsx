'use client'

import { useState } from 'react'
import { Check, ChevronsUpDown, Plus } from 'lucide-react'

import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { Badge } from '@/components/ui/badge'
import { useActiveCompany, useCompanies, useSetActiveCompany, useAddCompany, useLoadCompanies } from '@/lib/stores/companyStore'
import CompanyModal, { CompanyFormData } from './CompanyModal'
import { apiClient } from '@/lib/stores/authStore'
import { useNotificationStore } from '@/lib/stores/notificationStore'

export function CompanySelector() {
  const activeCompany = useActiveCompany()
  const companies = useCompanies()
  const setActiveCompany = useSetActiveCompany()
  const addCompany = useAddCompany()
  const loadCompanies = useLoadCompanies()
  const { showSuccess, showError } = useNotificationStore()
  
  const [isCompanyModalOpen, setIsCompanyModalOpen] = useState(false)
  const [isCreatingCompany, setIsCreatingCompany] = useState(false)

  const handleCompanySelect = (companyId: string) => {
    const company = companies.find(c => c.id === companyId)
    if (company) {
      setActiveCompany(company)
    }
  }

  const handleCreateCompany = async (companyData: CompanyFormData): Promise<boolean> => {
    setIsCreatingCompany(true)
    
    try {
      // Call backend API to create company
      const response = await apiClient.authPost('/companies', companyData)
      
      // Convert response to frontend format
      const newCompany = {
        id: response.uuid,
        uuid: response.uuid,
        name: response.name,
        industry: response.industry || '',
        size: response.size || '',
        role: response.role || '',
        isDefault: response.is_default || false,
        createdAt: new Date(response.created_at),
        updatedAt: new Date(response.created_at)
      }

      // Update store
      addCompany(newCompany)
      
      // Set as active company if it's the first one or marked as default
      if (companies.length === 0 || newCompany.isDefault) {
        setActiveCompany(newCompany)
      }

      showSuccess('Sucesso', `Empresa "${newCompany.name}" foi criada com sucesso!`)
      
      // Reload companies to ensure consistency
      await loadCompanies()
      
      return true
    } catch (error) {
      console.error('Error creating company:', error)
      showError('Erro', 'Não foi possível criar a empresa. Tente novamente.')
      return false
    } finally {
      setIsCreatingCompany(false)
    }
  }

  return (
    <>
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            variant="outline"
            role="combobox"
            className="w-64 justify-between"
          >
            <div className="flex items-center gap-2">
              <span className="truncate">
                {activeCompany?.name || 'Selecionar empresa...'}
              </span>
              {activeCompany?.isDefault && (
                <Badge variant="secondary" className="text-xs">
                  Atual
                </Badge>
              )}
            </div>
            <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent className="w-64 p-0">
          <DropdownMenuLabel className="px-2 py-1.5 text-sm font-normal">
            Suas empresas
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          
          {companies.map((company) => (
            <DropdownMenuItem
              key={company.id}
              className="flex items-center justify-between px-2 py-2"
              onClick={() => handleCompanySelect(company.id)}
            >
              <div className="flex items-center gap-2">
                <Check
                  className={`h-4 w-4 ${
                    activeCompany?.id === company.id ? 'opacity-100' : 'opacity-0'
                  }`}
                />
                <span className="truncate">{company.name}</span>
              </div>
              {company.isDefault && (
                <Badge variant="secondary" className="text-xs">
                  Atual
                </Badge>
              )}
            </DropdownMenuItem>
          ))}
          
          <DropdownMenuSeparator />
          <DropdownMenuItem 
            className="gap-2 px-2 py-2 cursor-pointer"
            onClick={() => setIsCompanyModalOpen(true)}
          >
            <Plus className="h-4 w-4" />
            Adicionar empresa
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
      
      {/* Company Creation Modal */}
      <CompanyModal
        open={isCompanyModalOpen}
        onClose={() => setIsCompanyModalOpen(false)}
        onSubmit={handleCreateCompany}
        isLoading={isCreatingCompany}
      />
    </>
  )
}