# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LeaderPro** - Uma plataforma de IA que amplifica a inteligência de liderança, mantendo memória perfeita de cada interação com o time e sugerindo ações contextuais.

**Tagline:** "Torne-se um líder mais inteligente"

## Product Context

- **Problema**: 70% dos líderes são promovidos sem treinamento adequado, resultando em 1:1s ineficazes e performance reviews superficiais
- **Solução**: IA contextual multidimensional que combina dados pessoais + temporais + geográficos + histórico
- **Diferencial**: Coach virtual que lembra tudo e sugere ações baseadas em contexto completo
- **Modelo**: B2C (R$ 49,90/mês) - líder paga individualmente

## Project Structure

### Current State
- `/backend/` - Empty directory (Go API to be implemented)
- `/frontend/` - **✅ IMPLEMENTED** - Complete Next.js 15.3.5 frontend with all core features
- `/plan/` - Complete project planning documentation
- `/CLAUDE.md` - This file with development guidance
- `/README.md` - Project overview and business context

### Key Documentation
- **`/frontend/README.md`** - Frontend architecture, TypeScript types, components, and stores
- **`/plan/000001-projeto-techlead.md`** - Business plan, market analysis, and strategy

## Technology Stack

### Backend (Go/Golang)
- **Framework**: Gin ou Echo
- **ORM**: GORM
- **Database**: PostgreSQL (dados relacionais)
- **Cache**: Redis (sessões e cache de IA)
- **IA**: OpenAI GPT-4 ou Claude API
- **Vector DB**: Pinecone/Weaviate (memória contextual)

### Frontend (Next.js) - ✅ IMPLEMENTED
- **Framework**: Next.js 15.3.5 with App Router
- **Language**: TypeScript with strict mode
- **Styling**: TailwindCSS v4 + shadcn/ui components
- **State**: Zustand for global state management
- **Icons**: Lucide React
- **Date Utils**: date-fns
- **Auth**: Clerk ou Auth0 (to be integrated)
- **Deploy**: GitHub Pages via GitHub Actions

## Brand Identity

### Design System
- **Cores primárias**: Azul #2563eb, Verde #16a34a
- **Cores secundárias**: Cinza #475569, Branco #ffffff
- **Acentos**: Laranja #ea580c, Roxo #7c3aed
- **Tipografia**: Inter
- **Ícones**: Heroicons (outline style)
- **Estilo**: Minimalista, inspirado em Linear/Notion

### Tom de Voz
- Profissional mas acessível
- Inteligente sem ser arrogante
- Prático e orientado a resultados
- Empático com as dores do líder

## Core Features

1. **1:1s com IA Contextual**
   - Sugestões de perguntas baseadas em: perfil + data + localização + histórico
   - Exemplo: "João tem filhos + julho + SP = pergunte sobre férias escolares"

2. **Memória Contínua**
   - Timeline visual de cada pessoa
   - IA conecta eventos: "Maria melhorou comunicação desde feedback de março"

3. **Anotações Inteligentes**
   - Feedbacks do dia a dia
   - Conquistas e observações
   - Performance reviews automatizadas

4. **Sistema de Menções e Feedback Cruzado**
   - Durante 1:1s, use `@nome` para referenciar outras pessoas
   - Exemplo: "Maria disse que @João é muito atencioso e ajuda bastante"
   - Automaticamente cria feedback no perfil do João: "Feedback de Maria via 1:1"
   - Se a pessoa não existir, sugere criar perfil básico
   - Constrói rede de percepções e relacionamentos do time

## Development Guidelines

- **Backend**: Siga padrões Go idiomáticos, use interfaces para abstração
- **Frontend**: Componentes reutilizáveis com shadcn/ui, typecheck rigoroso
- **IA**: Otimize prompts para reduzir custos, implemente cache inteligente
- **UX**: Foque na simplicidade - onboarding em 5 minutos
- **Dados**: LGPD compliance obrigatório, criptografia end-to-end
- **Performance**: Otimize para mobile, lazy loading, edge functions

## Key Metrics

- **MVP**: 10 usuários ativos semanalmente
- **6 meses**: 100 usuários pagantes
- **North Star**: Número de 1:1s realizados por mês

## Development Commands

### Current Status
**✅ FRONTEND COMPLETED** - Full Next.js implementation with all core features  
**⚠️ BACKEND PENDING** - API implementation needed

### Frontend (Next.js) - ✅ IMPLEMENTED
```bash
# Development server (with Turbopack)
cd frontend && npm run dev

# Production build (static export)
npm run build

# Production server
npm run start

# Linting
npm run lint

# Type checking (manual - not in package.json)
npx tsc --noEmit

# Install dependencies
npm install

# Add new dependencies
npm install [package-name]
```

### Backend (Go) - To be implemented
```bash
# Initialize project (run once)
cd backend && go mod init github.com/diegoclair/leaderpro

# Run development server
go run cmd/api/main.go

# Run tests
go test ./...

# Build for production
go build -o bin/api cmd/api/main.go

# Database migrations (when implemented)
migrate -path db/migrations -database "postgresql://..." up
```

### Deployment
```bash
# Frontend deploys automatically via GitHub Actions on push to main branch
# Live at: https://diegoclair.github.io/leaderpro/
```

## Architecture Overview

### Backend Architecture
```
backend/
├── cmd/api/          # Application entrypoints
├── internal/         # Private application code
│   ├── domain/       # Business entities and interfaces
│   ├── service/      # Business logic (includes @mention parsing)
│   ├── repository/   # Data access layer
│   ├── handler/      # HTTP handlers
│   └── ai/           # AI integration (prompts, vector storage)
├── pkg/              # Public packages
├── db/migrations/    # Database migrations
└── config/           # Configuration files
```

### Frontend Architecture
```
frontend/
├── app/                    # Next.js 15 App Router
│   ├── components/         # Page-specific components
│   ├── fonts/              # Local font files
│   ├── globals.css         # Global styles and Tailwind
│   ├── layout.tsx          # Root layout with providers
│   └── page.tsx            # Main dashboard page
├── components/             
│   ├── ui/                 # shadcn/ui components
│   ├── layout/             # Header, navigation
│   ├── people/             # Person profiles, forms
│   ├── one-on-ones/        # 1:1 meeting components
│   ├── feedback/           # Feedback management
│   └── shared/             # Reusable components
├── hooks/                  # Custom React hooks
│   ├── use-mentions.ts     # @mention functionality
│   └── use-create-person.ts# Auto-create from mentions
├── lib/                    
│   ├── utils/              # Date, name, formatting utils
│   └── utils.ts            # cn() utility for classes
└── stores/                 # Zustand state stores
    ├── company-store.ts    # Company/project management
    ├── people-store.ts     # Team member profiles
    └── one-on-one-store.ts # 1:1 meetings and notes
```

### AI Context Pipeline
1. **Data Collection**: User inputs, team member profiles, calendar data
2. **Vector Embedding**: Convert to embeddings using OpenAI/Claude
3. **Context Building**: Combine temporal + geographic + personal data
4. **Prompt Engineering**: Generate contextual suggestions
5. **Response Caching**: Store frequent queries in Redis

## Implementation Status

### ✅ Completed Features (Frontend)
- Multi-company/project support with data persistence
- Complete person profile management (personal + professional data)
- 1:1 meeting system with notes and action items
- @mention system with auto-suggestions and cross-referencing
- Feedback tracking (direct + mentioned)
- Dark/light theme toggle
- Responsive design with shadcn/ui components
- Local data persistence (localStorage + Zustand)

### ⏳ Pending Implementation
1. **Backend API** (Go + PostgreSQL + Redis)
2. **AI Integration** (OpenAI/Claude for contextual suggestions)
3. **Authentication** (Clerk/Auth0)
4. **Vector Database** (Pinecone/Weaviate for AI memory)
5. **Real-time Sync** (WebSockets for collaborative features)
6. **Testing Framework** (Jest/Vitest + Playwright)
7. **Performance Reviews** (Auto-generation from 1:1s)
8. **Calendar Integration** (Google/Outlook)

## Code Quality and Best Practices

### TypeScript Configuration
- Strict mode enabled with all checks
- Path aliases: `@/*` maps to `./src/*`
- Target: ES2022
- Module resolution: bundler

### Common Patterns
```typescript
// Component variants with CVA
const buttonVariants = cva("base-classes", {
  variants: {
    size: { sm: "...", md: "...", lg: "..." }
  }
})

// Zustand store pattern
interface StoreState {
  items: Item[]
  addItem: (item: Item) => void
  updateItem: (id: string, updates: Partial<Item>) => void
}

// Date formatting
formatRelativeDate(date) // "2 days ago"
formatDate(date, "PPP") // "January 7, 2025"
```

### Performance Considerations
- Static export for fast loading
- Turbopack for faster dev builds
- Component lazy loading where appropriate
- Zustand persist middleware for offline support

## Important Notes

- **Current Status**: Frontend completed, backend in planning phase
- **Live Demo**: https://diegoclair.github.io/leaderpro/
- **Data Storage**: Currently using localStorage (temporary)
- **AI Features**: Mocked in frontend, pending backend implementation
- **Authentication**: Not implemented, all features publicly accessible

## Quick Start

```bash
# Clone and install
git clone https://github.com/diegoclair/leaderpro.git
cd leaderpro/frontend
npm install

# Run development server
npm run dev

# Open http://localhost:3000
```

## Common Tasks

### Adding a New Feature
1. Check `/frontend/README.md` for architecture patterns
2. Create components in appropriate directories
3. Update relevant Zustand stores
4. Follow existing TypeScript patterns
5. Test with multiple companies/profiles

### Debugging
- Check browser console for Zustand state updates
- Use React Developer Tools for component inspection
- TypeScript errors: run `npx tsc --noEmit`
- Build issues: check `next.config.js` for static export settings