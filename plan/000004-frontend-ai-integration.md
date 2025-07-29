# Plano de Integração da IA no Frontend (SIMPLIFICADO)

## ✅ Situação Atual
- **Backend IA**: 100% implementado com chat API funcional
- **Frontend**: `AIAssistantSidebar` já existe e está sendo usado no profile
- **Arquitetura**: Padrões estabelecidos (Zustand, apiClient, types, constants)

## 🎯 Objetivo Ultra Simples
Conectar o `AIAssistantSidebar` existente com a API real de IA. **Não criar nada novo**, apenas fazer funcionar.

## 📋 TODO List (1-2 dias)

### 1. Adicionar Endpoints de IA (30 min)
```typescript
// src/lib/constants/api-endpoints.ts
export const AI_ENDPOINTS = {
  CHAT: (companyUuid: string, personUuid: string) => 
    `/companies/${companyUuid}/people/${personUuid}/ai/chat`,
  USAGE: (companyUuid: string, period?: string) => 
    `/companies/${companyUuid}/ai/usage${period ? `?period=${period}` : ''}`,
  FEEDBACK: (companyUuid: string, usageId: number) => 
    `/companies/${companyUuid}/ai/usage/${usageId}/feedback`,
} as const
```

### 2. Criar Types de IA (15 min)
```typescript
// src/lib/types/ai.ts
export interface ChatRequest {
  message: string
}

export interface ChatResponse {
  response: string
  usageId: number
  usage: {
    inputTokens: number
    outputTokens: number
    totalTokens: number
    costUSD: number
  }
}
```

### 3. Implementar Chat Hook (45 min)
```typescript
// src/hooks/useAIChat.ts
export const useAIChat = (personUuid: string) => {
  const { currentCompany } = useCompanyStore()
  const [isLoading, setIsLoading] = useState(false)

  const sendMessage = async (message: string): Promise<ChatResponse> => {
    setIsLoading(true)
    try {
      const response = await apiClient.authPost(
        AI_ENDPOINTS.CHAT(currentCompany.uuid, personUuid),
        { message }
      )
      return response
    } finally {
      setIsLoading(false)
    }
  }

  return { sendMessage, isLoading }
}
```

### 4. Conectar AIAssistantSidebar (30 min)
```typescript
// src/components/profile/AIAssistantSidebar.tsx
// Substituir o TODO por:
const { sendMessage, isLoading } = useAIChat(person.uuid)

const handleSendMessage = async () => {
  if (chatMessage.trim() && !isLoading) {
    const response = await sendMessage(chatMessage)
    // Exibir resposta na UI
    setChatMessage('')
  }
}
```

### 5. Adicionar Estado para Resposta (15 min)
```typescript
// No AIAssistantSidebar, adicionar:
const [lastResponse, setLastResponse] = useState<string>('')

// Exibir a resposta na UI
```

## 🎨 UI Já Existe
- ✅ Design aprovado e funcionando
- ✅ Layout responsivo
- ✅ Estados de loading
- ✅ Quick suggestions
- ✅ Input handling

## 🚀 Resultado Final
- User clica em "Perguntas para 1:1" → input popula → clica "Enviar"
- Aparece "Pensando..." (já existe)
- IA responde com sugestões contextuais baseadas no histórico real da pessoa
- **Zero mudanças na UI**, só conectar com backend real

## 📦 Entrega MVP
**1 tarde de trabalho** para ter IA funcional no LeaderPro!

### Etapas:
1. ✅ Backend (já feito)
2. 🔧 Frontend endpoints + types (1h)
3. 🔌 Conectar componente existente (30min)
4. ✅ Testar e ajustar (30min)

## 🔮 Futuras Melhorias (Opcionais)
- Histórico de conversas (se necessário)
- Insights automáticos no dashboard  
- Relatórios de usage (se solicitado)
- Feedback rating system

**Filosofia**: Fazer funcionar primeiro, otimizar depois.