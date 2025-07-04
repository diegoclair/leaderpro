'use client'

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
import { useActiveCompany, useCompanies, useSetActiveCompany } from '@/lib/stores/companyStore'

export function CompanySelector() {
  const activeCompany = useActiveCompany()
  const companies = useCompanies()
  const setActiveCompany = useSetActiveCompany()

  const handleCompanySelect = (companyId: string) => {
    const company = companies.find(c => c.id === companyId)
    if (company) {
      setActiveCompany(company)
    }
  }

  return (
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
        <DropdownMenuItem className="gap-2 px-2 py-2">
          <Plus className="h-4 w-4" />
          Adicionar empresa
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}