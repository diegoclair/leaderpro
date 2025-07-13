# 🗄️ Storage Manager - Sistema Centralizado de localStorage

## Problema Resolvido

**Antes (❌ Problemático):**
- Cada store gerenciava seu próprio localStorage
- No logout, precisava chamar múltiplas funções de clear
- Empresa de usuário anterior ficava persistida
- Risco de dados vazarem entre usuários

**Depois (✅ Seguro):**
```typescript
// Uma única chamada limpa TUDO
storageManager.clearAll()
```

## Como Usar

### 1. Salvando Dados
```typescript
import { storageManager } from '@/lib/utils/storageManager'

// Ao invés de localStorage.setItem()
storageManager.set('leaderpro-active-company', companyId)
```

### 2. Lendo Dados
```typescript
// Ao invés de localStorage.getItem()
const companyId = storageManager.get<string>('leaderpro-active-company')
```

### 3. Removendo Dados Específicos
```typescript
// Ao invés de localStorage.removeItem()
storageManager.remove('leaderpro-active-company')
```

### 4. Limpeza Completa (Logout)
```typescript
// Limpa TODOS os dados do LeaderPro de uma vez
storageManager.clearAll()
```

## Chaves Gerenciadas

O Storage Manager controla todas estas chaves:
- `auth-storage` - Dados de autenticação (Zustand persist)
- `leaderpro-active-company` - Empresa ativa
- `leaderpro-theme` - Tema do usuário  
- `leaderpro-onboarding` - Status do onboarding

## Segurança

✅ **Antes**: Bug crítico - empresa de outro usuário ficava ativa
✅ **Depois**: Logout limpa 100% dos dados, impossível vazar dados

## Debugging

```typescript
// Para debug - mostra todos os dados salvos
storageManager.debug()
```

## Implementação nos Stores

### CompanyStore
```typescript
// Antes
localStorage.setItem('leaderpro-active-company', company.id)

// Depois  
storageManager.set('leaderpro-active-company', company.id)
```

### AuthStore (Logout)
```typescript
// Antes - múltiplas chamadas
get().clearAuth()
useCompanyStore.getState().clearCompanyData()
usePeopleStore.getState().clearPeopleData()

// Depois - uma única chamada
storageManager.clearAll()
```

## Vantagens

1. **Centralizado**: Um lugar para todas as operações de storage
2. **Seguro**: Impossível esquecer de limpar algum dado no logout
3. **Debugável**: Fácil de ver todos os dados salvos
4. **Extensível**: Adicionar novas chaves é simples
5. **Type-safe**: TypeScript garante uso de chaves corretas