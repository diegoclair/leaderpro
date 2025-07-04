# LeaderPro Frontend

## Sobre o Projeto

LeaderPro é uma plataforma de IA que amplifica a inteligência de liderança, mantendo memória perfeita de cada interação com o time e sugerindo ações contextuais.

**Tagline:** "Torne-se um líder mais inteligente"

## Arquitetura Frontend

### Stack Tecnológico

- **Framework**: Next.js 14 com App Router
- **TypeScript**: Tipagem estática completa
- **Styling**: TailwindCSS + shadcn/ui
- **Estado**: Zustand (gerenciamento de estado global)
- **Build**: Turbopack (desenvolvimento mais rápido)
- **Linting**: ESLint configurado

### Estrutura de Diretórios

```
src/
├── app/                 # Next.js 14 App Router
│   ├── (auth)/          # Rotas que requerem autenticação
│   ├── onboarding/      # Setup inicial e criação de empresa
│   ├── dashboard/       # Dashboard principal
│   ├── person/          # Perfis individuais
│   ├── settings/        # Configurações de empresa e usuário
│   └── api/             # API routes (quando necessário)
├── components/          # Componentes reutilizáveis
│   ├── ui/              # Componentes shadcn/ui
│   ├── company/         # Componentes relacionados a empresa
│   ├── person/          # Componentes de perfil de pessoas
│   ├── oneononoe/       # Componentes de 1:1
│   └── layout/          # Layout e navegação
├── lib/                 # Utilitários e helpers
│   ├── stores/          # Zustand stores
│   ├── types/           # Definições TypeScript
│   ├── utils/           # Funções utilitárias
│   └── data/            # Dados mockados (temporário)
└── hooks/               # Custom React hooks
```

## Funcionalidades Principais

### 1. Multi-Empresas
- Líder pode gerenciar múltiplas empresas
- Histórico separado por empresa
- Empresa padrão configurável
- Portabilidade de dados ao mudar de empresa

### 2. Sistema de Menções (@mentions)
- Durante 1:1s, use `@nome` para referenciar outras pessoas
- Cria automaticamente feedback cruzado no perfil da pessoa mencionada
- Sugere criação de perfil se pessoa não existir

### 3. Coach de IA Contextual
- Sugestões de perguntas baseadas em contexto completo
- Combina dados pessoais + temporais + geográficos + histórico
- Conexões inteligentes entre eventos

## Dados e Estado

### Estrutura de Dados Principal

```typescript
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
    location?: string
    interests?: string[]
  }
  nextOneOnOne?: Date
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

## Componentes Principais

### Layout e Navegação
- `AppHeader` - Header com seletor de empresa e navegação
- `CompanySelector` - Dropdown para trocar empresa ativa
- `Sidebar` - Navegação lateral (dashboard, pessoas, configurações)

### Dashboard
- `DashboardOverview` - Visão geral de próximos 1:1s e métricas
- `PersonCard` - Card compacto de cada pessoa
- `UpcomingOneOnOnes` - Lista de próximas reuniões

### Perfil de Pessoa
- `PersonProfile` - Perfil completo com timeline
- `PersonTimeline` - Timeline de eventos e feedbacks
- `OneOnOneHistory` - Histórico de 1:1s
- `CrossReferenceFeedback` - Feedbacks vindos de @mentions

### 1:1 Sessions
- `OneOnOneSession` - Interface principal para conduzir reunião
- `AISuggestions` - Sugestões contextuais da IA
- `MentionInput` - Input que detecta @mentions automaticamente
- `SessionNotes` - Área de anotações estruturadas

## Boas Práticas de Desenvolvimento

### TypeScript
- Sempre usar tipagem estrita
- Definir interfaces para todos os dados
- Evitar `any` - usar `unknown` quando necessário
- Utilizar type guards para validação

### Componentes
- Usar composição sobre herança
- Componentes pequenos e focados em uma responsabilidade
- Props tipadas com interfaces
- Usar React.memo() para otimização quando necessário

### Estado
- Zustand para estado global (empresa ativa, usuário)
- useState para estado local de componentes
- Evitar prop drilling - usar context ou store quando necessário

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
npm run dev              # Inicia servidor de desenvolvimento
npm run build           # Build de produção
npm run start           # Inicia servidor de produção
npm run lint            # Executa ESLint
npm run type-check      # Verificação de tipos TypeScript

# Componentes shadcn/ui
npx shadcn@latest add [component]  # Adiciona componente
npx shadcn@latest list             # Lista componentes disponíveis
```

## Dados Mockados (Temporário)

### Empresas
1. **TechCorp** (padrão) - 2 pessoas
2. **StartupXYZ** - 2 pessoas

### Pessoas por Empresa
**TechCorp:**
- Maria Santos (Senior Developer)
- João Silva (Mid-level Developer)

**StartupXYZ:**
- Pedro Costa (Tech Lead)
- Ana Lima (Junior Developer)

## Próximos Passos

1. ✅ Setup inicial do projeto
2. ✅ Configuração do shadcn/ui
3. 🔄 Implementação da estrutura básica
4. ⏳ Criação dos stores Zustand
5. ⏳ Componentes de layout
6. ⏳ Dashboard principal
7. ⏳ Sistema de multi-empresas
8. ⏳ Perfis de pessoas
9. ⏳ Sistema de 1:1s
10. ⏳ Sistema de @mentions

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