'use client'

import React, { useState, useEffect } from 'react'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { FilterToggleButton } from '@/components/ui/FilterToggleButton'
import { 
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { 
  Popover,
  PopoverContent,
  PopoverTrigger,
} from '@/components/ui/popover'
import { 
  Search, 
  Filter, 
  Calendar,
  Star,
  MessageSquare,
  Eye,
  Clock,
  X
} from 'lucide-react'

export interface FilterOptions {
  searchQuery: string
  quickView: string | null
  types: string[]
  period: string
  direction: string
  sentiment: string[]
}

interface FilterBarProps {
  filters: FilterOptions
  onFiltersChange: (filters: FilterOptions) => void
  totalItems: number
  className?: string
}

const quickViews = [
  { 
    id: 'notes', 
    name: '📋 Notas', 
    icon: Eye,
    filters: { 
      types: ['observation']
    }
  },
  { 
    id: 'feedbacks', 
    name: '⭐ Feedbacks', 
    icon: Star,
    filters: { 
      types: ['feedback']
    }
  },
  { 
    id: 'meetings', 
    name: '💬 1:1s', 
    icon: MessageSquare,
    filters: { 
      types: ['one_on_one'] 
    }
  },
  { 
    id: 'recent', 
    name: '🕒 Últimos 30 dias', 
    icon: Clock,
    filters: { period: '30d' }
  }
]

const activityTypes = [
  { value: 'observation', label: '📋 Notas', icon: Eye },
  { value: 'feedback', label: '⭐ Feedbacks', icon: Star },
  { value: 'one_on_one', label: '💬 1:1s', icon: MessageSquare }
]

const periodOptions = [
  { value: '7d', label: 'Últimos 7 dias' },
  { value: '30d', label: 'Últimos 30 dias' },
  { value: '3m', label: 'Últimos 3 meses' },
  { value: '6m', label: 'Últimos 6 meses' },
  { value: '1y', label: 'Último ano' },
  { value: 'all', label: 'Todos os períodos' }
]


const sentimentOptions = [
  { value: 'positive', label: '✨ Positivos', color: 'bg-green-500' },
  { value: 'constructive', label: '🔨 Construtivos', color: 'bg-yellow-500' },
  { value: 'neutral', label: '⚪ Neutros', color: 'bg-gray-500' }
]

export function FilterBar({ 
  filters, 
  onFiltersChange, 
  totalItems, 
  className = '' 
}: FilterBarProps) {
  const [showAdvanced, setShowAdvanced] = useState(false)
  const [tempFilters, setTempFilters] = useState({
    types: filters.types,
    sentiment: filters.sentiment
  })
  
  const handleTypeFilter = (types: string[]) => {
    // Toggle types - if all types are already selected, remove them; otherwise add them
    const currentTypes = new Set(filters.types)
    const allTypesSelected = types.every(type => currentTypes.has(type))
    
    let newTypes: string[]
    if (allTypesSelected) {
      // Remove these types
      newTypes = filters.types.filter(type => !types.includes(type))
    } else {
      // Add missing types
      newTypes = [...new Set([...filters.types, ...types])]
    }
    
    onFiltersChange({
      ...filters,
      types: newTypes,
      quickView: null
    })
  }

  const handlePeriodFilter = (period: string) => {
    const isActive = filters.period === period
    
    onFiltersChange({
      ...filters,
      period: isActive ? 'all' : period,
      quickView: isActive ? null : 'recent'
    })
  }
  
  
  const toggleTempType = (type: string) => {
    const newTypes = tempFilters.types.includes(type)
      ? tempFilters.types.filter(t => t !== type)
      : [...tempFilters.types, type]
    
    setTempFilters({
      ...tempFilters,
      types: newTypes
    })
  }
  
  const toggleTempSentiment = (sentiment: string) => {
    const newSentiments = tempFilters.sentiment.includes(sentiment)
      ? tempFilters.sentiment.filter(s => s !== sentiment)
      : [...tempFilters.sentiment, sentiment]
    
    setTempFilters({
      ...tempFilters,
      sentiment: newSentiments
    })
  }

  const applyAdvancedFilters = () => {
    onFiltersChange({
      ...filters,
      types: tempFilters.types,
      sentiment: tempFilters.sentiment,
      quickView: null
    })
    setShowAdvanced(false)
  }

  const resetAdvancedFilters = () => {
    setTempFilters({
      types: [],
      sentiment: []
    })
  }

  // Update temp filters when external filters change (e.g., from quick views)
  useEffect(() => {
    setTempFilters({
      types: filters.types,
      sentiment: filters.sentiment
    })
  }, [filters.types, filters.sentiment])
  
  const clearAllFilters = () => {
    onFiltersChange({
      searchQuery: '',
      quickView: null,
      types: [],
      period: 'all',
      direction: 'all',
      sentiment: []
    })
  }
  
  const hasActiveFilters = filters.quickView || 
    filters.searchQuery || 
    filters.types.length > 0 || 
    filters.period !== 'all' || 
    filters.sentiment.length > 0

  // Count advanced filters (sentiment filters are exclusive to advanced panel)
  const advancedFilterCount = filters.sentiment.length
  

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Search Bar */}
      <div className="flex items-center gap-3">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
          <Input
            placeholder="Buscar 'feedback positivo performance'..."
            value={filters.searchQuery}
            onChange={(e) => onFiltersChange({ ...filters, searchQuery: e.target.value })}
            className="pl-10"
          />
        </div>
        
        <div className="flex items-center gap-2 text-sm text-muted-foreground">
          <Calendar className="h-4 w-4" />
          <Select 
            value={filters.period} 
            onValueChange={(period) => onFiltersChange({ ...filters, period, quickView: null })}
          >
            <SelectTrigger className="w-40">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {periodOptions.map((option) => (
                <SelectItem key={option.value} value={option.value}>
                  {option.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {/* Quick Views */}
      <div className="space-y-2">
        <div className="flex items-center justify-between">
          <h3 className="text-sm font-medium">Quick Views:</h3>
          <div className="flex items-center gap-2">
            <span className="text-xs text-muted-foreground">
              {totalItems > 0 ? `${totalItems} atividades encontradas` : 'Nenhuma atividade'}
            </span>
            {hasActiveFilters && (
              <Button
                variant="ghost"
                size="sm"
                onClick={clearAllFilters}
                className="h-6 px-2 text-xs"
              >
                <X className="h-3 w-3 mr-1" />
                Limpar
              </Button>
            )}
          </div>
        </div>
        
        <div className="flex flex-wrap gap-2">
          {quickViews.map((quickView) => {
            // Check if this filter is active
            let isActive = false
            let handleClick: () => void
            
            if (quickView.filters.types) {
              // Type filters (Notes, Feedbacks, 1:1s)
              isActive = quickView.filters.types.every(type => filters.types.includes(type))
              handleClick = () => handleTypeFilter(quickView.filters.types!)
            } else if (quickView.filters.period) {
              // Period filters (Recent)
              isActive = filters.period === quickView.filters.period
              handleClick = () => handlePeriodFilter(quickView.filters.period!)
            } else {
              // Fallback
              isActive = false
              handleClick = () => {}
            }
            
            return (
              <FilterToggleButton
                key={quickView.id}
                id={quickView.id}
                name={quickView.name}
                icon={quickView.icon}
                isActive={isActive}
                onClick={handleClick}
              />
            )
          })}
          
          <Popover open={showAdvanced} onOpenChange={setShowAdvanced}>
            <PopoverTrigger asChild>
              <Button 
                variant={advancedFilterCount > 0 ? "default" : "outline"} 
                size="sm" 
                className="h-8 gap-2"
              >
                <Filter className="h-3 w-3" />
                Filtros Avançados
                {advancedFilterCount > 0 && (
                  <span className="bg-background text-foreground px-1.5 py-0.5 rounded-full text-xs font-medium">
                    {advancedFilterCount}
                  </span>
                )}
              </Button>
            </PopoverTrigger>
            <PopoverContent className="w-80" align="start">
              <div className="space-y-4">
                <h4 className="font-medium">Filtros Avançados</h4>
                
                {/* Activity Types */}
                <div>
                  <label className="text-sm font-medium mb-2 block">Tipos de Atividade</label>
                  <div className="flex flex-wrap gap-1">
                    {activityTypes.map((type) => {
                      const Icon = type.icon
                      const isActive = tempFilters.types.includes(type.value)
                      
                      return (
                        <Button
                          key={type.value}
                          type="button"
                          variant={isActive ? "default" : "outline"}
                          size="sm"
                          onClick={() => toggleTempType(type.value)}
                          className="h-7 text-xs gap-1"
                        >
                          <Icon className="h-3 w-3" />
                          {type.label}
                        </Button>
                      )
                    })}
                  </div>
                </div>
                
                {/* Sentiment */}
                <div>
                  <label className="text-sm font-medium mb-2 block">Sentimento</label>
                  <div className="flex flex-wrap gap-1">
                    {sentimentOptions.map((sentiment) => {
                      const isActive = tempFilters.sentiment.includes(sentiment.value)
                      
                      return (
                        <Button
                          key={sentiment.value}
                          type="button"
                          variant={isActive ? "default" : "outline"}
                          size="sm"
                          onClick={() => toggleTempSentiment(sentiment.value)}
                          className="h-7 text-xs gap-1"
                        >
                          <div className={`w-2 h-2 rounded-full ${sentiment.color}`} />
                          {sentiment.label}
                        </Button>
                      )
                    })}
                  </div>
                </div>
                
                {/* Action Buttons */}
                <div className="flex items-center justify-between pt-4 border-t border-border">
                  <Button
                    type="button"
                    variant="ghost"
                    size="sm"
                    onClick={resetAdvancedFilters}
                    className="text-xs"
                  >
                    Limpar Tudo
                  </Button>
                  
                  <div className="flex gap-2">
                    <Button
                      type="button"
                      variant="outline"
                      size="sm"
                      onClick={() => setShowAdvanced(false)}
                      className="text-xs"
                    >
                      Cancelar
                    </Button>
                    <Button
                      type="button"
                      variant="default"
                      size="sm"
                      onClick={applyAdvancedFilters}
                      className="text-xs"
                    >
                      Aplicar Filtros
                    </Button>
                  </div>
                </div>
                
              </div>
            </PopoverContent>
          </Popover>
        </div>
      </div>

    </div>
  )
}