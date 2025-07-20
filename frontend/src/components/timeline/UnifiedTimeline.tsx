'use client'

import React, { useState, useEffect, useMemo } from 'react'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { ErrorMessage } from '@/components/ui/ErrorMessage'
import { Button } from '@/components/ui/button'
import { Pagination } from '@/components/ui/pagination'
import { RefreshCw, InboxIcon } from 'lucide-react'
import { SimpleActivityCard, TimelineActivity } from './SimpleActivityCard'
import { FilterBar, FilterOptions } from './FilterBar'
import { EditNoteModal } from './EditNoteModal'
import { Person } from '@/lib/types'
import { apiClient } from '@/lib/stores/authStore'
import { useNotificationStore } from '@/lib/stores/notificationStore'

interface UnifiedTimelineProps {
  person: Person
  companyId: string
  allPeople: Person[]
  className?: string
}

export function UnifiedTimeline({ 
  person, 
  companyId, 
  allPeople, 
  className = '' 
}: UnifiedTimelineProps) {
  const [activities, setActivities] = useState<TimelineActivity[]>([])
  const [mentions, setMentions] = useState<TimelineActivity[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const { showSuccess, showError } = useNotificationStore()
  const [currentPage, setCurrentPage] = useState(1)
  
  // Edit modal state
  const [isEditModalOpen, setIsEditModalOpen] = useState(false)
  const [editingActivity, setEditingActivity] = useState<TimelineActivity | null>(null)
  const [totalRecords, setTotalRecords] = useState(0)
  const [itemsPerPage, setItemsPerPage] = useState(10)
  const [totalPages, setTotalPages] = useState(0)
  
  const [filters, setFilters] = useState<FilterOptions>({
    searchQuery: '',
    quickView: null,
    types: [],
    period: 'all',
    direction: 'all',
    sentiment: []
  })

  // Fetch timeline data with server-side filtering
  const fetchTimelineData = async (page: number = 1) => {
    setIsLoading(true)
    setError(null)
    
    try {
      // Build query parameters for server-side filtering
      const params = new URLSearchParams()
      params.append('page', page.toString())
      params.append('quantity', itemsPerPage.toString())
      
      // Add filters to query parameters - use debouncedSearchQuery for search
      if (debouncedSearchQuery) {
        params.append('search_query', debouncedSearchQuery)
      }
      
      if (filters.types.length > 0) {
        params.append('types', filters.types.join(','))
      }
      
      if (filters.sentiment.length > 0) {
        params.append('feedback_types', filters.sentiment.join(','))
      }
      
      if (filters.direction !== 'all') {
        params.append('direction', filters.direction)
      }
      
      if (filters.period !== 'all') {
        params.append('period', filters.period)
      }
      
      // Fetch unified timeline with filters
      const response = await apiClient.authGet(
        `/companies/${companyId}/people/${person.uuid}/timeline?${params.toString()}`
      )
      
      // Handle unified timeline response
      let timelineData = []
      let paginationData = { total_records: 0 }
      
      if (response?.data && Array.isArray(response.data)) {
        timelineData = response.data
      }
      
      if (response?.pagination) {
        paginationData = response.pagination
      }
      
      setTotalRecords(paginationData.total_records)
      setTotalPages(Math.ceil(paginationData.total_records / itemsPerPage))
      setCurrentPage(page)
      
      // Since we now have unified data, we set both activities and mentions to the same data
      // but separate them in the filtering logic below for UI compatibility
      const directActivities = timelineData.filter((item: any) => item.entry_source === 'direct')
      const mentionActivities = timelineData.filter((item: any) => item.entry_source === 'mention')
      
      // Always replace data for traditional pagination
      setActivities(directActivities)
      setMentions(mentionActivities)
      
    } catch (err) {
      console.error('Error fetching timeline data:', err)
      setError('Erro ao carregar timeline. Tente novamente.')
    } finally {
      setIsLoading(false)
    }
  }

  // Handle edit note
  const handleEditNote = (activity: TimelineActivity) => {
    setEditingActivity(activity)
    setIsEditModalOpen(true)
  }

  // Handle save edited note
  const handleSaveEditedNote = async (updatedActivity: TimelineActivity) => {
    try {
      const payload = {
        type: updatedActivity.type,
        content: updatedActivity.content,
        feedback_type: updatedActivity.feedback_type || undefined,
        feedback_category: updatedActivity.feedback_category === 'none' ? undefined : updatedActivity.feedback_category || undefined,
        // TODO: Handle mentioned_people if needed
        mentioned_people: []
      }

      await apiClient.authPut(
        `/companies/${companyId}/people/${person.uuid}/notes/${updatedActivity.uuid}`,
        payload
      )
      
      showSuccess('Sucesso', 'Anotação atualizada com sucesso')
      
      // Refresh timeline data
      await fetchTimelineData(currentPage)
      
    } catch (error) {
      console.error('Error updating note:', error)
      showError('Erro', 'Não foi possível atualizar a anotação. Tente novamente.')
      throw error // Re-throw to let modal handle loading state
    }
  }

  // Handle delete note
  const handleDeleteNote = async (activity: TimelineActivity) => {
    if (!confirm('Tem certeza que deseja excluir esta anotação?')) {
      return
    }

    try {
      await apiClient.authDelete(
        `/companies/${companyId}/people/${person.uuid}/notes/${activity.uuid}`
      )
      
      showSuccess('Sucesso', 'Anotação excluída com sucesso')
      
      // Refresh timeline data
      await fetchTimelineData(currentPage)
      
    } catch (error) {
      console.error('Error deleting note:', error)
      showError('Erro', 'Não foi possível excluir a anotação. Tente novamente.')
    }
  }
  
  // Pagination handlers
  const handlePageChange = (page: number) => {
    if (page !== currentPage && page >= 1 && page <= totalPages) {
      fetchTimelineData(page)
    }
  }
  
  const handleItemsPerPageChange = (newItemsPerPage: number) => {
    setItemsPerPage(newItemsPerPage)
    setCurrentPage(1)
    // Will trigger fetchTimelineData via useEffect
  }

  // Debounced search query
  const [debouncedSearchQuery, setDebouncedSearchQuery] = useState(filters.searchQuery)
  
  // Debounce search query changes
  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setDebouncedSearchQuery(filters.searchQuery)
    }, 800)
    
    return () => clearTimeout(timeoutId)
  }, [filters.searchQuery])
  
  // Fetch data when person, company, or debounced search changes
  useEffect(() => {
    if (person && companyId) {
      fetchTimelineData()
    }
  }, [person, companyId, debouncedSearchQuery])
  
  // Fetch data immediately when non-search filters change
  useEffect(() => {
    if (person && companyId) {
      fetchTimelineData(1) // Reset to page 1 when filters change
    }
  }, [person, companyId, filters.types, filters.sentiment, filters.direction, filters.period, filters.quickView, itemsPerPage])

  // Simply combine activities since filtering is now done server-side
  const filteredActivities = useMemo(() => {
    // Combine all activities (already filtered by server)
    const allActivities = [...activities, ...mentions]
    
    // Sort by date (newest first) - server should handle this but ensure consistency
    return allActivities.sort((a, b) => 
      new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
    )
  }, [activities, mentions])

  if (isLoading) {
    return (
      <div className={`flex items-center justify-center py-12 ${className}`}>
        <div className="text-center">
          <LoadingSpinner size="large" />
          <p className="text-muted-foreground mt-4">Carregando timeline...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className={`py-8 ${className}`}>
        <ErrorMessage 
          message={error}
          actionButton={
            <Button onClick={() => fetchTimelineData()} variant="outline" size="sm">
              <RefreshCw className="h-4 w-4 mr-2" />
              Tentar Novamente
            </Button>
          }
        />
      </div>
    )
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-xl font-semibold">
            Timeline de Atividades
          </h2>
          <p className="text-sm text-muted-foreground">
            Histórico de feedbacks, 1:1s e observações
          </p>
        </div>
      </div>

      {/* Filter Bar */}
      <FilterBar
        filters={filters}
        onFiltersChange={setFilters}
        totalItems={filteredActivities.length}
      />

      {/* Timeline */}
      <div className="space-y-4">
        {filteredActivities.length === 0 ? (
          <div className="text-center py-12">
            <InboxIcon className="h-12 w-12 text-muted-foreground mx-auto mb-4" />
            <h3 className="text-lg font-medium mb-2">Nenhuma atividade encontrada</h3>
            <p className="text-muted-foreground mb-4">
              {filters.searchQuery || filters.types.length > 0 || filters.quickView
                ? 'Tente ajustar os filtros para ver mais resultados.'
                : 'Ainda não há atividades registradas para esta pessoa.'}
            </p>
            {(filters.searchQuery || filters.types.length > 0 || filters.quickView) && (
              <Button
                onClick={() => setFilters({
                  searchQuery: '',
                  quickView: null,
                  types: [],
                  period: 'all',
                  direction: 'all',
                  sentiment: []
                })}
                variant="outline"
                size="sm"
              >
                Limpar Filtros
              </Button>
            )}
          </div>
        ) : (
          <>
            {filteredActivities.map((activity, index) => (
              <SimpleActivityCard
                key={`${activity.uuid}-${index}`}
                activity={activity}
                onEdit={handleEditNote}
                onDelete={handleDeleteNote}
              />
            ))}
            
            {/* Pagination */}
            <Pagination
              currentPage={currentPage}
              totalPages={totalPages}
              itemsPerPage={itemsPerPage}
              totalItems={totalRecords}
              onPageChange={handlePageChange}
              onItemsPerPageChange={handleItemsPerPageChange}
              className="mt-6"
            />
          </>
        )}
      </div>

      {/* Edit Note Modal */}
      <EditNoteModal
        isOpen={isEditModalOpen}
        onClose={() => {
          setIsEditModalOpen(false)
          setEditingActivity(null)
        }}
        activity={editingActivity}
        allPeople={allPeople}
        onSave={handleSaveEditedNote}
      />
    </div>
  )
}