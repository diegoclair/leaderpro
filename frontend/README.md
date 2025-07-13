# LeaderPro Frontend

## Sobre o Projeto

LeaderPro é uma plataforma de IA que amplifica a inteligência de liderança, mantendo memória perfeita de cada interação com o time e sugerindo ações contextuais.

**Tagline:** "Torne-se um líder mais inteligente"

## Arquitetura Frontend

### Stack Tecnológico

- **Framework**: Next.js 15.3.5 com App Router
- **Runtime**: React 19.0.0
- **TypeScript**: Tipagem estática completa (modo strict)
- **Styling**: TailwindCSS v4 + shadcn/ui
- **Estado**: Zustand (gerenciamento de estado global com persistência)
- **Build**: Turbopack (desenvolvimento mais rápido)
- **Linting**: ESLint configurado

### Estrutura de Diretórios

```
src/
├── app/                 # Next.js 15.3.5 App Router
│   ├── auth/            # Páginas de autenticação
│   │   ├── layout.tsx   # Layout para páginas de auth
│   │   ├── login/       # Página de login
│   │   └── register/    # Página de registro
│   ├── dashboard/       # Dashboard principal da aplicação
│   │   └── page.tsx     # Página principal com cards e métricas
│   ├── profile/[id]/    # Perfis individuais
│   │   ├── layout.tsx   # Layout para páginas de perfil
│   │   ├── page.tsx     # Página de perfil detalhado
│   │   └── static-params.ts # Parâmetros estáticos para build
│   ├── page.tsx         # Landing page (redireciona baseado em auth)
│   ├── layout.tsx       # Layout raiz da aplicação
│   ├── globals.css      # Estilos globais + react-mentions styling
│   └── middleware.ts    # Middleware Next.js para proteção de rotas
├── components/          # Componentes reutilizáveis
│   ├── auth/            # Componentes de autenticação
│   │   └── AuthGuard.tsx # Proteção de rotas autenticadas
│   ├── ui/              # Componentes shadcn/ui + compartilhados
│   │   ├── button.tsx   # shadcn/ui components
│   │   ├── card.tsx
│   │   ├── input.tsx
│   │   ├── loading-spinner.tsx  # ✅ Componente compartilhado (10+ usos)
│   │   ├── app-logo.tsx         # ✅ Componente compartilhado (4+ usos)
│   │   ├── error-message.tsx    # ✅ Componente compartilhado
│   │   ├── password-input.tsx   # ✅ Componente compartilhado
│   │   ├── phone-input.tsx      # ✅ Componente compartilhado
│   │   ├── submit-button.tsx    # ✅ Componente compartilhado
│   │   ├── mentions-input.tsx   # ✅ Sistema de mentions (react-mentions)
│   │   └── notifications.tsx
│   ├── company/         # Componentes relacionados a empresa
│   │   └── CompanySelector.tsx
│   ├── onboarding/      # Componentes de onboarding
│   │   └── OnboardingWizard.tsx # Wizard inicial para novos usuários
│   ├── person/          # Componentes de perfil de pessoas
│   │   └── PersonCard.tsx
│   ├── profile/         # Componentes específicos do perfil
│   │   ├── PersonInfoTab.tsx
│   │   ├── PersonHistoryTab.tsx
│   │   ├── PersonFeedbackTab.tsx
│   │   ├── PersonChatTab.tsx
│   │   └── CreatePersonDialog.tsx
│   └── layout/          # Layout e navegação
│       ├── AppHeader.tsx
│       ├── ThemeProvider.tsx # Provider do tema
│       └── ThemeToggle.tsx
├── lib/                 # Utilitários e helpers
│   ├── stores/          # Zustand stores com persistência
│   │   ├── authStore.ts # Store de autenticação (integrado com API)
│   │   ├── companyStore.ts # Store de empresas
│   │   ├── peopleStore.ts  # Store de pessoas e 1:1s
│   │   └── notificationStore.ts # Store de notificações
│   ├── constants/       # ✅ Constantes centralizadas
│   │   ├── api.ts       # ✅ API endpoints centralizados
│   │   ├── company.ts   # ✅ Tamanhos e indústrias de empresa (padrão BR)
│   │   ├── messages.ts  # ✅ Mensagens de erro/sucesso
│   │   └── validation.ts # ✅ Regras de validação
│   ├── types/           # Definições TypeScript
│   │   └── index.ts     # Types principais da aplicação
│   ├── utils/           # Funções utilitárias
│   │   ├── dates.ts     # Formatação de datas
│   │   ├── names.ts     # Utilitários de nomes
│   │   ├── cn.ts        # Utility function para classes CSS
│   │   └── storageManager.ts # ✅ Storage Manager centralizado
│   └── api/             # ✅ Cliente API centralizado
│       └── client.ts    # ✅ apiClient com auth automática
└── hooks/               # Custom React hooks
    ├── useMentions.ts   # Hook para sistema de @menções (react-mentions)
    ├── useCreatePerson.ts # Hook para criação de pessoas
    └── useAuthRedirect.ts # Hook para redirecionamento de auth
```

## Funcionalidades Principais

### 1. Sistema de Autenticação Completo
- **Login/Registro**: Interface moderna com validação
- **JWT Tokens**: Access token (15min) + Refresh token (24h)
- **Proteção de Rotas**: Middleware automático + AuthGuard
- **Renovação Automática**: Tokens renovados transparentemente
- **Persistência**: Estado mantido entre sessões
- **Logout Seguro**: Invalidação local + notificação ao backend

### 2. Multi-Empresas
- Líder pode gerenciar múltiplas empresas
- Histórico separado por empresa
- Empresa padrão configurável
- Portabilidade de dados ao mudar de empresa

### 3. Sistema de Menções (@mentions) - react-mentions
- Durante 1:1s, use `@nome` para referenciar outras pessoas
- Interface visual: mostra `@Nome` no texto
- Backend: envia formato `{{person:uuid|name}}` 
- Biblioteca: react-mentions para UX profissional
- Autocomplete com pessoas da empresa
- Cria automaticamente feedback cruzado no perfil da pessoa mencionada
- Sugere criação de perfil se pessoa não existir

### 4. Coach de IA Contextual
- Sugestões de perguntas baseadas em contexto completo
- Combina dados pessoais + temporais + geográficos + histórico
- Conexões inteligentes entre eventos

## Dados e Estado

### Estrutura de Dados Principal

```typescript
// Usuário (Autenticação)
interface User {
  uuid: string
  email: string
  name: string
  phone?: string
  profilePhoto?: string
  plan: string
  trialEndsAt?: string
  subscribedAt?: string
  timezone?: string
  language?: string
  createdAt: string
  updatedAt: string
  lastLoginAt?: string
  emailVerified: boolean
}

// Empresa
type Company = {
  id: string
  name: string
  isDefault: boolean
  createdAt: Date
  updatedAt: Date
}

// Pessoa
type Person = {
  id: string
  companyId: string
  name: string
  role: string
  email?: string
  startDate: Date
  personalInfo: {
    hasChildren?: boolean
    hasPets?: boolean
    location?: string
    interests?: string[]
    hobbies?: string[]
    pets?: string[]
    personalNotes?: string
  }
  nextOneOnOne?: Date
  avatar?: string
}

// 1:1 Session
type OneOnOneSession = {
  id: string
  personId: string
  date: Date
  notes: string
  aiSuggestions: string[]
  mentions: Mention[]
  status: 'scheduled' | 'completed' | 'cancelled'
}

// Mention (@nome)
type Mention = {
  id: string
  sessionId: string
  mentionedPersonId: string
  context: string
  createdAt: Date
}
```

### Gerenciamento de Estado (Zustand)

```typescript
// Company Store
const useCompanyStore = create((set) => ({
  companies: [],
  activeCompany: null,
  setActiveCompany: (company) => set({ activeCompany: company }),
  addCompany: (company) => set((state) => ({ 
    companies: [...state.companies, company] 
  }))
}))

// People Store
const usePeopleStore = create((set) => ({
  people: [],
  addPerson: (person) => set((state) => ({ 
    people: [...state.people, person] 
  })),
  getPeopleByCompany: (companyId) => // filtrar por empresa
}))
```

## Comunicação com Backend

### API Integration Layer

O frontend implementa comunicação completa com o backend através de:

#### AuthStore - Sistema de Autenticação
```typescript
// Store principal de autenticação
const useAuthStore = create<AuthStore>()(
  persist(
    (set, get) => ({
      // State
      user: User | null,
      tokens: AuthTokens | null,
      isLoading: boolean,
      isAuthenticated: boolean,

      // Actions - Integradas com API real
      login: (email, password) => POST /auth/login
      register: (data) => POST /users + auto-login
      logout: () => POST /auth/logout + clearAuth
      refreshToken: () => POST /auth/refresh-token
      getProfile: () => GET /users/profile
    }),
    { name: 'auth-storage' } // Persistência local
  )
)
```

#### API Client Centralizado ⭐ IMPORTANTE
```typescript
// ✅ SEMPRE use o apiClient centralizado - NUNCA use fetch() diretamente
import { apiClient } from '@/lib/stores/authStore'

// Requisições públicas (sem autenticação)
const loginData = await apiClient.post('/auth/login', { email, password })
const userData = await apiClient.post('/users', registrationData)

// Requisições autenticadas (token automático + renovação)
const profile = await apiClient.authGet('/users/profile')
await apiClient.authPost('/auth/logout')
await apiClient.authPut('/users/profile', updateData)
await apiClient.authDelete('/companies/123')

// ❌ INCORRETO - Nunca use fetch() diretamente
// fetch('/api/endpoint') // NÃO FAÇA ISSO

// Funcionalidades automáticas:
// 1. Headers de autenticação (user-token) incluídos automaticamente
// 2. Renovação automática de token em caso de 401
// 3. Gerenciamento centralizado de erros
// 4. Configuração única - mudanças de header em 1 lugar só
// 5. Type safety com TypeScript
```

#### Endpoints Implementados
- **Autenticação**: `/auth/login`, `/auth/logout`, `/auth/refresh-token`
- **Usuários**: `/users` (POST), `/users/profile` (GET)
- **Backend URL**: Configurável via `NEXT_PUBLIC_API_URL` (default: http://localhost:5000)

#### Fluxo de Autenticação
1. **Login**: Email/senha → JWT access + refresh tokens
2. **Persistência**: Tokens salvos no localStorage via Zustand persist
3. **Interceptor**: Todas as requisições autenticadas passam pelo `authFetch`
4. **Renovação**: Token expirado é renovado automaticamente e transparentemente
5. **Logout**: Limpa tokens local + notifica backend

#### Próximas APIs (Pendentes)
- **Empresas**: `/companies` (CRUD)
- **Pessoas**: `/people` (CRUD)
- **1:1s**: `/one-on-ones` (CRUD)
- **AI**: `/ai/suggestions` (contextual insights)

## Componentes Principais

### Layout e Navegação
- `AppHeader` - Header com seletor de empresa, navegação e tema
- `CompanySelector` - Dropdown para trocar empresa ativa (com persistência)
- `ThemeToggle` - Alternador de tema claro/escuro

### Dashboard
- `Dashboard` (page.tsx) - Página principal com cards de pessoas e métricas
- `PersonCard` - Card clicável de cada pessoa com informações básicas

### Perfil de Pessoa
- `ProfilePage` - Layout principal da página de perfil
- `PersonInfoTab` - Aba de informações pessoais + formulário de registro
- `PersonHistoryTab` - Aba do histórico de 1:1s da pessoa
- `PersonFeedbackTab` - Aba de feedbacks recebidos (diretos + @menções)
- `PersonChatTab` - Aba de chat com IA para insights

### Sistema de @Menções
- `MentionSuggestions` - Dropdown de sugestões ao digitar @
- `CreatePersonDialog` - Dialog para criar pessoa inexistente
- `useMentions` - Hook para detectar e gerenciar @menções
- `useCreatePerson` - Hook para criação de pessoas

### Utilitários
- `dates.ts` - Formatação de datas (formatTimeAgo, formatShortDate, etc.)
- `names.ts` - Utilitários de nomes (getInitials)

## Arquitetura e Boas Práticas Implementadas

### 🗄️ Storage Manager Centralizado ⭐ CRÍTICO PARA SEGURANÇA
```typescript
// ✅ Storage Manager - Gerencia TODO o localStorage de forma segura
import { storageManager } from '@/lib/utils/storageManager'

// Uso nos stores
storageManager.set('leaderpro-active-company', companyId)
const companyId = storageManager.get<string>('leaderpro-active-company')

// Logout seguro - limpa TODOS os dados de uma vez
storageManager.clearAll()

// ❌ NUNCA use localStorage diretamente
// localStorage.setItem() // NÃO FAÇA ISSO
```

### 🔧 Componentes Compartilhados ⭐ ELIMINA DUPLICAÇÃO
```typescript
// ✅ SEMPRE procure por componentes existentes antes de criar novos
import { LoadingSpinner } from '@/components/ui/loading-spinner'  // 10+ usos
import { AppLogo } from '@/components/ui/app-logo'                // 4+ usos  
import { ErrorMessage } from '@/components/ui/error-message'      // Shared
import { PhoneInput } from '@/components/ui/phone-input'          // Shared
import { PasswordInput } from '@/components/ui/password-input'    // Shared
import { SubmitButton } from '@/components/ui/submit-button'      // Shared

// ✅ Constantes centralizadas
import { API_ENDPOINTS } from '@/lib/constants/api'
import { COMPANY_SIZES } from '@/lib/constants/company'  // Padrões brasileiros
import { MESSAGES } from '@/lib/constants/messages'
import { VALIDATION } from '@/lib/constants/validation'
```

### Separação de Responsabilidades
- **Componentes de UI** - Focados apenas na apresentação
- **Hooks personalizados** - Lógica de negócio reutilizável
- **Utils** - Funções puras sem estado
- **Stores** - Gerenciamento de estado global
- **Constants** - Valores centralizados e reutilizáveis

### Padrões Adotados
- **Composição sobre herança** - Componentes pequenos e compostos
- **Single Responsibility** - Cada componente tem uma função específica
- **Custom Hooks** - Abstração de lógica complexa (useMentions, useCreatePerson)
- **Utils centralizados** - Evita duplicação de código
- **Constants compartilhadas** - Valores únicos em local centralizado
- **TypeScript estrito** - Tipagem completa em todos os componentes

### Performance
- **Componentes otimizados** - Componentes pequenos carregam mais rápido
- **Memoização** - useMemo para cálculos pesados
- **Lazy loading** - Componentes carregados sob demanda
- **Estados locais** - Evita re-renders desnecessários

### Limpeza de Código Realizada (2025-01-13)
**Problema:** ~200 linhas de código duplicado em múltiplos componentes
**Solução:** Criação de componentes compartilhados e constantes centralizadas

**Exemplos de refatoração:**
- **LoadingSpinner**: Eliminado 10+ duplicações
- **AppLogo**: Eliminado 4+ duplicações  
- **PhoneInput**: Componente único para máscara de telefone brasileiro
- **Company Constants**: Tamanhos de empresa padronizados para o Brasil
- **Storage Manager**: Segurança total na troca de usuários

## Boas Práticas de Desenvolvimento ⭐ LEIA ANTES DE CODAR

### 🔍 ANTES DE CRIAR QUALQUER COMPONENTE
1. **Procure primeiro**: Verifique se já existe em `/components/ui/`
2. **Analise duplicação**: O novo componente pode ser compartilhado?
3. **Local correto**: Se compartilhado, crie em `/components/ui/`
4. **Constants**: Use `/lib/constants/` para valores fixos

### TypeScript
- Sempre usar tipagem estrita
- Definir interfaces para todos os dados
- Evitar `any` - usar `unknown` quando necessário
- Utilizar type guards para validação

### Componentes
- **⚠️ SEMPRE procurar componentes existentes antes de criar novos**
- **⚠️ Se for reutilizável, criar em `/components/ui/` imediatamente**
- Usar composição sobre herança
- Componentes pequenos e focados em uma responsabilidade
- Props tipadas com interfaces
- Usar React.memo() para otimização quando necessário

### Estado e Storage
- **⚠️ NUNCA usar localStorage diretamente - SEMPRE usar storageManager**
- Zustand para estado global (empresa ativa, usuário)
- useState para estado local de componentes
- Evitar prop drilling - usar context ou store quando necessário

### API Communication
- **⚠️ SEMPRE usar apiClient - NUNCA usar fetch() diretamente**
- Headers de autenticação incluídos automaticamente
- Renovação de token transparente

### Styling
- TailwindCSS para estilização
- Componentes shadcn/ui como base
- Design system consistente (cores definidas no projeto)
- Responsive design (mobile-first)

### Performance
- Lazy loading de páginas e componentes pesados
- Otimização de imagens (Next.js Image)
- Memoização de cálculos pesados
- Paginação para listas grandes

## Design System

### Cores (definidas no projeto)
- **Primárias**: Azul #2563eb, Verde #16a34a
- **Secundárias**: Cinza #475569, Branco #ffffff
- **Acentos**: Laranja #ea580c, Roxo #7c3aed

### Tipografia
- **Font**: Inter (já configurada no Next.js)
- **Tamanhos**: Usar classes Tailwind (text-sm, text-base, text-lg, etc.)

### Ícones
- **Biblioteca**: Lucide React (outline style)
- **Uso**: Consistência visual, tamanho padrão 16px ou 20px

## Comandos de Desenvolvimento

```bash
# Desenvolvimento
npm run dev              # Inicia servidor de desenvolvimento (http://localhost:3000)
npm run build           # Build de produção (static export)
npm run start           # Inicia servidor de produção
npm run lint            # Executa ESLint
npx tsc --noEmit        # Verificação de tipos TypeScript

# Componentes shadcn/ui
npx shadcn@latest add [component]  # Adiciona componente
npx shadcn@latest list             # Lista componentes disponíveis
```

### Variáveis de Ambiente

```bash
# .env.local (opcional)
NEXT_PUBLIC_API_URL=http://localhost:5000  # URL do backend (default se não definido)
```

## Dados Mockados (Temporário)

### Empresas
1. **TechCorp** (padrão) - 2 pessoas
2. **StartupXYZ** - 2 pessoas

### Pessoas por Empresa
**TechCorp:**
- Maria Santos (Analista de Sistemas)
- João Silva (Coordenador de Projetos)

**StartupXYZ:**
- Pedro Costa (Gerente de Projetos)
- Ana Lima (Analista Junior)

## Status de Implementação

### ✅ Recentemente Implementado (2025-01-13)
1. **🗄️ Storage Manager Centralizado**: Sistema seguro de localStorage que impede vazamento de dados entre usuários
2. **🔧 Componentes Compartilhados**: Eliminação de ~200 linhas duplicadas com componentes reutilizáveis
3. **📱 Sistema de Mentions**: react-mentions para UX profissional (@nome → {{person:uuid|name}})
4. **🇧🇷 Padrões Brasileiros**: Tamanhos de empresa conforme SEBRAE/IBGE (10-49, 50-99, 100-499, 500+)
5. **🌙 Tema Dark Default**: Interface escura como padrão
6. **🏢 Auto-seleção de Empresa**: Fix no onboarding para selecionar empresa automaticamente
7. **📱 Phone Input Centralizado**: Máscara brasileira compartilhada
8. **⚡ Constantes Centralizadas**: API endpoints, mensagens, validações em arquivos únicos

### ✅ Anteriormente Implementado
9. **Sistema de Autenticação**: JWT/PASETO completo com refresh automático
10. **Gerenciamento de Empresas**: Criação, listagem e seleção de empresas
11. **Associação Usuário-Empresa**: Modelo simplificado onde cada empresa pertence diretamente a um usuário
12. **Fluxo de Onboarding**: Wizard inicial para criação da primeira empresa
13. **Integração Banco de Dados**: Criação real de empresas no backend MySQL
14. **API Endpoints**: `/companies` (POST/GET), `/users` (POST), `/auth/*` implementados

### ✅ Frontend Completo
1. **Setup e Base**: Projeto Next.js 15.3.5 + TailwindCSS v4 + shadcn/ui
2. **Arquitetura**: Estrutura modular com components/hooks/stores/utils
3. **Autenticação**: Sistema completo JWT com refresh automático
4. **API Layer**: authStore + authFetch interceptor integrado com backend
5. **Stores Zustand**: authStore (real) + companyStore (real) + peopleStore (mock)
6. **Layout**: AppHeader, ThemeProvider, AuthGuard, middleware de proteção
7. **Páginas**: Landing, Auth (login/register), Dashboard, Profile
8. **Dashboard**: Cards de pessoas, métricas, seletor de empresa
9. **Perfis**: Abas completas (info, histórico, feedback, chat IA)
10. **Sistema @mentions**: Autocomplete + criação automática de pessoas
11. **Onboarding**: Wizard inicial para novos usuários integrado com backend
12. **Feedbacks cruzados**: Sistema de menções entre perfis

### ⏳ Em Desenvolvimento
1. **API Business**: Endpoints de pessoas e 1:1s no backend
2. **Migração Mock→API**: Migrar peopleStore para usar APIs reais
3. **Sistema Member Get Member**: Programa de indicações com desconto
4. **IA Contextual**: Integração OpenAI/Claude para suggestions

### 📋 Próximos Passos
1. **People API**: Implementar endpoints de pessoas no backend Go
2. **1:1 Meetings API**: Implementar endpoints de reuniões 1:1
3. **Data Migration**: Migrar peopleStore para usar APIs reais
4. **IA Integration**: Sistema de suggestions contextuais
5. **Notificações**: Sistema de notificações em tempo real
6. **Mobile**: Responsividade completa mobile-first
7. **Testes**: Suíte de testes unitários e E2E

## Notas Importantes

- **Sempre ler esta documentação** antes de fazer mudanças significativas
- Manter tipagem TypeScript rigorosa
- Seguir padrões estabelecidos de estrutura de arquivos
- Atualizar esta documentação quando adicionar novas funcionalidades
- Validar compatibilidade de dependências antes de instalar
- Testar em múltiplas empresas para garantir isolamento de dados

---

## Getting Started (Next.js)

First, run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.