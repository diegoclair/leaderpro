'use client'

import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'

interface CreatePersonDialogProps {
  open: boolean
  onClose: () => void
  personName: string
  newPersonName: string
  newPersonRole: string
  setNewPersonName: (name: string) => void
  setNewPersonRole: (role: string) => void
  onCreatePerson: () => boolean
}

export default function CreatePersonDialog({
  open,
  onClose,
  personName,
  newPersonName,
  newPersonRole,
  setNewPersonName,
  setNewPersonRole,
  onCreatePerson
}: CreatePersonDialogProps) {
  const handleCreate = () => {
    const success = onCreatePerson()
    if (success) {
      onClose()
    }
  }

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Pessoa não encontrada</DialogTitle>
          <DialogDescription>
            A pessoa "@{personName}" foi mencionada mas não existe no sistema. 
            Deseja adicioná-la agora?
          </DialogDescription>
        </DialogHeader>
        
        <div className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="person-name">Nome completo</Label>
            <Input
              id="person-name"
              value={newPersonName}
              onChange={(e) => setNewPersonName(e.target.value)}
              placeholder="Digite o nome completo"
            />
          </div>
          
          <div className="space-y-2">
            <Label htmlFor="person-role">Cargo/Função</Label>
            <Input
              id="person-role"
              value={newPersonRole}
              onChange={(e) => setNewPersonRole(e.target.value)}
              placeholder="Ex: Analista, Coordenador, Gerente..."
            />
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={onClose}>
            Cancelar
          </Button>
          <Button onClick={handleCreate}>
            Adicionar Pessoa
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  )
}