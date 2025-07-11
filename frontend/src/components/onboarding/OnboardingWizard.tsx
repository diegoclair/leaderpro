'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Building2, Users, CheckCircle } from 'lucide-react'
import { useCompanyStore } from '@/lib/stores/companyStore'

interface OnboardingWizardProps {
  onComplete: () => void
}

export function OnboardingWizard({ onComplete }: OnboardingWizardProps) {
  const router = useRouter()
  const { addCompany, setActiveCompany } = useCompanyStore()
  
  const [step, setStep] = useState(1)
  const [isLoading, setIsLoading] = useState(false)
  const [formData, setFormData] = useState({
    companyName: '',
    industry: '',
    teamSize: '',
    userRole: ''
  })

  const handleNext = () => {
    if (step < 3) {
      setStep(step + 1)
    } else {
      handleComplete()
    }
  }

  const handleComplete = async () => {
    setIsLoading(true)
    
    try {
      // Criar empresa com dados do wizard
      const newCompany = {
        id: Date.now().toString(), // Temporary ID
        name: formData.companyName,
        description: `Empresa de ${formData.industry}`,
        industry: formData.industry,
        size: formData.teamSize,
        isDefault: true,
        createdAt: new Date(),
        updatedAt: new Date()
      }

      // Adicionar ao store
      addCompany(newCompany)
      setActiveCompany(newCompany)
      
      // Marcar onboarding como completo
      localStorage.setItem('onboarding_completed', 'true')
      
      // Chamar callback
      onComplete()
      
      // Redirecionar para dashboard
      router.push('/dashboard')
      
    } catch (error) {
      console.error('Erro ao criar empresa:', error)
    } finally {
      setIsLoading(false)
    }
  }

  const handleChange = (field: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }))
  }

  const isStepValid = () => {
    switch (step) {
      case 1:
        return formData.companyName.trim().length > 0
      case 2:
        return formData.industry && formData.teamSize
      case 3:
        return formData.userRole.trim().length > 0
      default:
        return false
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 px-4">
      <div className="w-full max-w-md">
        {/* Progress indicator */}
        <div className="flex justify-center mb-8">
          <div className="flex items-center space-x-2">
            {[1, 2, 3].map((stepNumber) => (
              <div key={stepNumber} className="flex items-center">
                <div className={`w-8 h-8 rounded-full flex items-center justify-center ${
                  stepNumber < step ? 'bg-green-500 text-white' :
                  stepNumber === step ? 'bg-blue-500 text-white' :
                  'bg-gray-200 text-gray-500'
                }`}>
                  {stepNumber < step ? <CheckCircle size={16} /> : stepNumber}
                </div>
                {stepNumber < 3 && (
                  <div className={`w-8 h-1 mx-2 ${
                    stepNumber < step ? 'bg-green-500' : 'bg-gray-200'
                  }`} />
                )}
              </div>
            ))}
          </div>
        </div>

        <Card className="w-full shadow-xl border-0 bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm">
          {step === 1 && (
            <>
              <CardHeader className="text-center">
                <div className="w-12 h-12 bg-blue-100 dark:bg-blue-900 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Building2 className="w-6 h-6 text-blue-600" />
                </div>
                <CardTitle className="text-xl">Empresa onde você trabalha</CardTitle>
                <CardDescription>
                  Digite o nome da empresa onde você é líder ou gerencia uma equipe
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="companyName">Nome da empresa</Label>
                  <Input
                    id="companyName"
                    placeholder="Ex: Google, Microsoft, sua startup..."
                    value={formData.companyName}
                    onChange={(e) => handleChange('companyName', e.target.value)}
                    className="mt-2"
                  />
                </div>
              </CardContent>
            </>
          )}

          {step === 2 && (
            <>
              <CardHeader className="text-center">
                <div className="w-12 h-12 bg-green-100 dark:bg-green-900 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <Users className="w-6 h-6 text-green-600" />
                </div>
                <CardTitle className="text-xl">Conte-nos sobre sua empresa</CardTitle>
                <CardDescription>
                  Essas informações nos ajudam a personalizar sua experiência
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="industry">Setor/Indústria</Label>
                  <Select value={formData.industry} onValueChange={(value) => handleChange('industry', value)}>
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Selecione o setor" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="technology">Tecnologia</SelectItem>
                      <SelectItem value="finance">Financeiro</SelectItem>
                      <SelectItem value="healthcare">Saúde</SelectItem>
                      <SelectItem value="education">Educação</SelectItem>
                      <SelectItem value="retail">Varejo</SelectItem>
                      <SelectItem value="manufacturing">Indústria</SelectItem>
                      <SelectItem value="consulting">Consultoria</SelectItem>
                      <SelectItem value="other">Outro</SelectItem>
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label htmlFor="teamSize">Tamanho da equipe</Label>
                  <Select value={formData.teamSize} onValueChange={(value) => handleChange('teamSize', value)}>
                    <SelectTrigger className="mt-2">
                      <SelectValue placeholder="Quantas pessoas você lidera?" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="1-5">1-5 pessoas</SelectItem>
                      <SelectItem value="6-15">6-15 pessoas</SelectItem>
                      <SelectItem value="16-50">16-50 pessoas</SelectItem>
                      <SelectItem value="51-100">51-100 pessoas</SelectItem>
                      <SelectItem value="100+">Mais de 100 pessoas</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </CardContent>
            </>
          )}

          {step === 3 && (
            <>
              <CardHeader className="text-center">
                <div className="w-12 h-12 bg-purple-100 dark:bg-purple-900 rounded-lg flex items-center justify-center mx-auto mb-4">
                  <CheckCircle className="w-6 h-6 text-purple-600" />
                </div>
                <CardTitle className="text-xl">Qual é o seu cargo?</CardTitle>
                <CardDescription>
                  Última etapa! Como você se identifica na empresa?
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div>
                  <Label htmlFor="userRole">Seu cargo/função</Label>
                  <Input
                    id="userRole"
                    placeholder="Ex: Tech Lead, Gerente, CTO..."
                    value={formData.userRole}
                    onChange={(e) => handleChange('userRole', e.target.value)}
                    className="mt-2"
                  />
                </div>
              </CardContent>
            </>
          )}

          <CardContent className="pt-0">
            <div className="flex gap-3">
              {step > 1 && (
                <Button
                  type="button"
                  variant="outline"
                  onClick={() => setStep(step - 1)}
                  className="flex-1"
                >
                  Voltar
                </Button>
              )}
              <Button
                type="button"
                onClick={handleNext}
                disabled={!isStepValid() || isLoading}
                className="flex-1"
              >
                {isLoading ? 'Criando...' : step === 3 ? 'Finalizar' : 'Próximo'}
              </Button>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  )
}