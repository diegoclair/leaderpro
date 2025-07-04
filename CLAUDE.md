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
- `/frontend/` - Empty directory (Next.js 14 to be implemented)
- `/plan/` - Complete project planning documentation
- `/CLAUDE.md` - This file with development guidance
- `/README.md` - Project overview and business context

### Planned Structure (When Implemented)
- `/backend/` - Go API com PostgreSQL e Redis
- `/frontend/` - Next.js 14 com TailwindCSS e shadcn/ui  
- `/plan/` - Documentação de produto e arquitetura

## Technology Stack

### Backend (Go/Golang)
- **Framework**: Gin ou Echo
- **ORM**: GORM
- **Database**: PostgreSQL (dados relacionais)
- **Cache**: Redis (sessões e cache de IA)
- **IA**: OpenAI GPT-4 ou Claude API
- **Vector DB**: Pinecone/Weaviate (memória contextual)

### Frontend (Next.js)
- **Framework**: Next.js 14 com App Router
- **Styling**: TailwindCSS + shadcn/ui
- **Estado**: Zustand
- **Auth**: Clerk ou Auth0
- **Deploy**: Vercel

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
**⚠️ PROJECT IS IN PLANNING PHASE - NO CODE IMPLEMENTED YET**

The backend/ and frontend/ directories are currently empty. Use the commands below when setting up the project:

### Project Initialization
```bash
# Create initial project structure
mkdir -p backend/{cmd/api,internal/{domain,service,repository,handler,ai},pkg,db/migrations,config}
mkdir -p frontend/{app,components,lib,hooks,stores}
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

### Frontend (Next.js) - To be implemented
```bash
# Initialize project (run once)
cd frontend && npx create-next-app@latest . --typescript --tailwind --app

# Install dependencies
npm install

# Run development server
npm run dev

# Run tests
npm test

# Build for production
npm run build

# Type checking
npm run type-check

# Linting
npm run lint
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
├── app/              # Next.js 14 App Router
│   ├── (auth)/       # Auth-required routes
│   ├── api/          # API routes
│   └── components/   # Page-specific components
├── components/       # Shared components
├── lib/              # Utilities and helpers
├── hooks/            # Custom React hooks
└── stores/           # Zustand state management
```

### AI Context Pipeline
1. **Data Collection**: User inputs, team member profiles, calendar data
2. **Vector Embedding**: Convert to embeddings using OpenAI/Claude
3. **Context Building**: Combine temporal + geographic + personal data
4. **Prompt Engineering**: Generate contextual suggestions
5. **Response Caching**: Store frequent queries in Redis

## Next Steps for Development

### Phase 1: Project Setup
1. **Backend Setup**:
   ```bash
   cd backend && go mod init github.com/diegoclair/leaderpro
   mkdir -p cmd/api internal/{domain,service,repository,handler,ai} pkg db/migrations config
   ```

2. **Frontend Setup**:
   ```bash
   cd frontend && npx create-next-app@latest . --typescript --tailwind --app
   npm install @radix-ui/react-* lucide-react zustand
   ```

3. **Database Setup**:
   - Configure PostgreSQL connection
   - Set up Redis for caching
   - Create initial migrations

### Phase 2: Core Features (MVP)
1. **User Authentication** - Clerk ou Auth0
2. **1:1 Management System** - CRUD operations
3. **AI Integration** - OpenAI/Claude API
4. **Basic Profile Management** - Team member profiles
5. **Note-taking System** - Structured feedback capture
6. **@Mention System** - Cross-reference parsing and profile linking

### Phase 3: Advanced Features
1. **Contextual AI Suggestions** - Temporal + personal data analysis
2. **Performance Review Generator** - Automated compilation
3. **Calendar Integration** - Google/Outlook sync
4. **Mobile App** - React Native ou PWA

## Important Notes

- Parte do ecossistema Smart (SmartBuy, SmartBank, LeaderPro)
- Modelo B2C: dados pertencem ao usuário, não à empresa
- Foco: líderes de qualquer área, adaptável a diferentes contextos organizacionais
- Validação contínua com early adopters essencial
- **Current Status**: Planning phase - no code implementation yet
- Directories `/backend/` and `/frontend/` are currently empty
- Project planning documentation available in `/plan/` directory

## Key Files and Documentation

- `/plan/000001-projeto-techlead.md` - Complete business plan and market analysis
- `/README.md` - Project overview and product vision
- This file contains technical architecture and development guidelines
- **`/frontend/README.md`** - Frontend-specific documentation (architecture, TypeScript types, Zustand stores, components)
- `/backend/README.md` - Backend API documentation (to be created)

## Important for AI Development

When working on this codebase:
1. **Always read module-specific documentation first**
2. **Frontend work**: Start by reading `/frontend/README.md` for current architecture
3. **Backend work**: Consult `/backend/README.md` for API structure
4. **Updates**: Keep documentation updated when adding new features