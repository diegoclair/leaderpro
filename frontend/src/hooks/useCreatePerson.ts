import { useState } from 'react'

export function useCreatePerson() {
  const [showCreatePersonDialog, setShowCreatePersonDialog] = useState(false)
  const [personToCreate, setPersonToCreate] = useState('')
  const [newPersonName, setNewPersonName] = useState('')
  const [newPersonRole, setNewPersonRole] = useState('')

  const openCreateDialog = (mentionName: string) => {
    setPersonToCreate(mentionName)
    setNewPersonName(mentionName)
    setShowCreatePersonDialog(true)
  }

  const closeCreateDialog = () => {
    setShowCreatePersonDialog(false)
    setPersonToCreate('')
    setNewPersonName('')
    setNewPersonRole('')
  }

  const handleCreatePerson = (companyId?: string) => {
    if (newPersonName.trim() && newPersonRole.trim()) {
      // TODO: Add person to store/backend
      console.log('Creating person:', {
        name: newPersonName,
        role: newPersonRole,
        companyId: companyId
      })
      closeCreateDialog()
      return true
    }
    return false
  }

  return {
    showCreatePersonDialog,
    personToCreate,
    newPersonName,
    newPersonRole,
    setNewPersonName,
    setNewPersonRole,
    openCreateDialog,
    closeCreateDialog,
    handleCreatePerson
  }
}