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
│   ├── page.tsx         # Dashboard principal
│   └── profile/[id]/    # Perfis individuais
│       └── page.tsx     # Página de perfil detalhado
├── components/          # Componentes reutilizáveis
│   ├── ui/              # Componentes shadcn/ui (button, card, etc.)
│   ├── company/         # Componentes relacionados a empresa
│   │   └── CompanySelector.tsx
│   ├── person/          # Componentes de perfil de pessoas
│   │   └── PersonCard.tsx
│   ├── profile/         # Componentes específicos do perfil
│   │   ├── PersonInfoTab.tsx
│   │   ├── PersonHistoryTab.tsx
│   │   ├── PersonFeedbackTab.tsx
│   │   ├── PersonChatTab.tsx
│   │   ├── MentionSuggestions.tsx
│   │   └── CreatePersonDialog.tsx
│   └── layout/          # Layout e navegação
│       ├── AppHeader.tsx
│       └── ThemeToggle.tsx
├── lib/                 # Utilitários e helpers
│   ├── stores/          # Zustand stores (companyStore, peopleStore)
│   ├── types/           # Definições TypeScript
│   ├── utils/           # Funções utilitárias
│   │   ├── dates.ts     # Formatação de datas
│   │   └── names.ts     # Utilitários de nomes
│   └── data/            # Dados mockados (temporário)
│       └── mockData.ts
└── hooks/               # Custom React hooks
    ├── useMentions.ts   # Hook para sistema de @menções
    └── useCreatePerson.ts # Hook para criação de pessoas
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

### Separação de Responsabilidades
- **Componentes de UI** - Focados apenas na apresentação
- **Hooks personalizados** - Lógica de negócio reutilizável
- **Utils** - Funções puras sem estado
- **Stores** - Gerenciamento de estado global

### Padrões Adotados
- **Composição sobre herança** - Componentes pequenos e compostos
- **Single Responsibility** - Cada componente tem uma função específica
- **Custom Hooks** - Abstração de lógica complexa (useMentions, useCreatePerson)
- **Utils centralizados** - Evita duplicação de código
- **TypeScript estrito** - Tipagem completa em todos os componentes

### Performance
- **Componentes otimizados** - Componentes pequenos carregam mais rápido
- **Memoização** - useMemo para cálculos pesados
- **Lazy loading** - Componentes carregados sob demanda
- **Estados locais** - Evita re-renders desnecessários

### Exemplo de Refatoração Realizada
**Antes:** ProfilePage tinha 676 linhas com tudo misturado
**Depois:** ProfilePage tem ~150 linhas + 6 componentes especializados

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
- Maria Santos (Analista de Sistemas)
- João Silva (Coordenador de Projetos)

**StartupXYZ:**
- Pedro Costa (Gerente de Projetos)
- Ana Lima (Analista Junior)

## Status de Implementação

1. ✅ Setup inicial do projeto
2. ✅ Configuração do shadcn/ui
3. ✅ Estrutura básica implementada
4. ✅ Stores Zustand (Company + People)
5. ✅ Componentes de layout (AppHeader, ThemeToggle)
6. ✅ Dashboard principal com métricas e cards
7. ✅ Sistema de multi-empresas com persistência
8. ✅ Perfis de pessoas com abas completas
9. ✅ Sistema de registro de 1:1s
10. ✅ Sistema de @mentions com autocomplete
11. ✅ Refatoração arquitetural (hooks + componentes modulares)
12. ✅ Sistema de feedbacks cruzados

## Próximos Passos
1. ⏳ Integração com backend real
2. ⏳ Sistema de autenticação
3. ⏳ IA contextual real (OpenAI/Claude API)
4. ⏳ Sistema de notificações
5. ⏳ Mobile responsiveness
6. ⏳ Testes unitários e E2E

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