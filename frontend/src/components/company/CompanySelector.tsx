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
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from '@/components/ui/tooltip'
import { useActiveCompany, useCompanies, useSetActiveCompany, useAddCompany, useLoadCompanies } from '@/lib/stores/companyStore'
import CompanyModal, { CompanyFormData } from './CompanyModal'
import { apiClient } from '@/lib/stores/authStore'
import { useNotificationStore } from '@/lib/stores/notificationStore'
import type { CompanyCreateResponse } from '@/lib/types/api'

export function CompanySelector() {
  const activeCompany = useActiveCompany()
  const companies = useCompanies()
  const setActiveCompany = useSetActiveCompany()
  const addCompany = useAddCompany()
  const loadCompanies = useLoadCompanies()
  const { showSuccess, showError } = useNotificationStore()
  
  const [isCompanyModalOpen, setIsCompanyModalOpen] = useState(false)
  const [isCreatingCompany, setIsCreatingCompany] = useState(false)

  // Helper function to truncate company name intelligently
  const truncateCompanyName = (name: string, maxLength: number = 30) => {
    if (name.length <= maxLength) return name
    return name.substring(0, maxLength) + '...'
  }

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
      const response = await apiClient.authPost<CompanyCreateResponse>('/companies', companyData)
      
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
        updatedAt: new Date(response.updated_at)
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
    <TooltipProvider>
      <DropdownMenu>
        {activeCompany?.name && activeCompany.name.length > 28 ? (
          <Tooltip>
            <TooltipTrigger asChild>
              <DropdownMenuTrigger asChild>
                <Button
                  variant="outline"
                  role="combobox"
                  className="w-64 justify-between"
                >
                  <div className="flex items-center gap-2 min-w-0 flex-1">
                    <span className="truncate text-left">
                      {truncateCompanyName(activeCompany.name, 28)}
                    </span>
                    {activeCompany?.isDefault && (
                      <Badge variant="secondary" className="text-xs shrink-0">
                        Atual
                      </Badge>
                    )}
                  </div>
                  <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
                </Button>
              </DropdownMenuTrigger>
            </TooltipTrigger>
            <TooltipContent>
              <p>{activeCompany.name}</p>
            </TooltipContent>
          </Tooltip>
        ) : (
          <DropdownMenuTrigger asChild>
            <Button
              variant="outline"
              role="combobox"
              className="w-64 justify-between"
            >
              <div className="flex items-center gap-2 min-w-0 flex-1">
                <span className="truncate text-left">
                  {activeCompany?.name || 'Selecionar empresa...'}
                </span>
                {activeCompany?.isDefault && (
                  <Badge variant="secondary" className="text-xs shrink-0">
                    Atual
                  </Badge>
                )}
              </div>
              <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </DropdownMenuTrigger>
        )}
        <DropdownMenuContent className="w-64 p-0">
          <DropdownMenuLabel className="px-2 py-1.5 text-sm font-normal">
            Suas empresas
          </DropdownMenuLabel>
          <DropdownMenuSeparator />
          
          {companies.map((company) => (
            <Tooltip key={company.id}>
              <TooltipTrigger asChild>
                <DropdownMenuItem
                  className="flex items-center justify-between px-2 py-2 cursor-pointer"
                  onClick={() => handleCompanySelect(company.id)}
                >
                  <div className="flex items-center gap-2 min-w-0 flex-1">
                    <Check
                      className={`h-4 w-4 shrink-0 ${
                        activeCompany?.id === company.id ? 'opacity-100' : 'opacity-0'
                      }`}
                    />
                    <span className="truncate text-left">
                      {truncateCompanyName(company.name, 32)}
                    </span>
                  </div>
                  {company.isDefault && (
                    <Badge variant="secondary" className="text-xs shrink-0">
                      Atual
                    </Badge>
                  )}
                </DropdownMenuItem>
              </TooltipTrigger>
              {company.name.length > 32 && (
                <TooltipContent>
                  <p>{company.name}</p>
                </TooltipContent>
              )}
            </Tooltip>
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
    </TooltipProvider>
  )
}