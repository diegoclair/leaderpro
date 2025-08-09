# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LeaderPro** - AI platform that amplifies leadership intelligence, maintaining perfect memory of every team interaction and suggesting contextual actions.

**Tagline:** "Become a smarter leader"

**Live Demo**: https://diegoclair.github.io/leaderpro/

## Quick Start

```bash
# Full stack development with hot reload
tilt up
# Visit: http://localhost:10350 (Tilt UI)
# Frontend: http://localhost:3000
# Backend: http://localhost:5000
# Swagger: http://localhost:5000/swagger/

# Stop everything
tilt down
```

## Architecture & Tech Stack

### Backend (Go) - Clean Architecture + DDD
- **Go 1.24.5+** with Echo v4.13.3 framework
- **Database**: MySQL 8.0.32 (no GORM - raw SQL for performance)
- **Cache**: Redis 7.4.2 for sessions and caching
- **Auth**: PASETO tokens (15min access, 24h refresh) via `user-token` header
- **Testing**: Testcontainers + go-sqlmock for unit/integration tests
- **Observability**: Prometheus, Grafana, Jaeger tracing
- **Architecture**: Domain-Driven Design with Clean Architecture principles
- **AI Providers**: OpenAI integrated, Anthropic support ready

### Frontend (Next.js)
- **Next.js 15.3.5** with App Router and React 19.0.0
- **TypeScript** with strict mode, path alias `@/*` → `./src/*`
- **Styling**: TailwindCSS v4 + shadcn/ui components
- **State**: Zustand stores with localStorage persistence
- **Icons**: Lucide React, **Date Utils**: date-fns
- **Deploy**: GitHub Pages via GitHub Actions (auto-deploy on main push)
- **Animation**: Performance-optimized with tw-animate-css

## Development Commands

### Frontend (Next.js)
```bash
cd frontend
npm run dev        # Development with Turbopack (http://localhost:3000)
npm run build      # Production build with static export
npm run lint       # ESLint linting
npx tsc --noEmit   # TypeScript type checking
```

### Backend (Go) 
```bash
cd backend
make start         # Start all services via Docker
make tests         # Run all tests with coverage
make mocks         # Generate mocks for testing
make docs          # Generate Swagger API docs
make clean-volumes # Clean Docker volumes (full reset)

# Testing commands
go test -v -cover ./...                           # All tests
go test -v -cover ./internal/application/service/ # Package tests
go test -v -run TestCreateCompany ./internal/application/service/  # Single test
go test -tags=integration -v ./...                # Integration tests
govulncheck ./...                                 # Security scan
```

## Architecture Overview

### Backend - Clean Architecture Flow
**Dependency Rule**: `transport` → `application` → `domain` ← `infra`

```
backend/
├── cmd/main.go               # App entry point
├── internal/
│   ├── domain/entity/        # Business entities (User, Company, Person, OneOnOne)
│   ├── application/service/  # Use cases/business logic  
│   └── transport/rest/       # HTTP handlers + routes
├── infra/
│   ├── data/mysql/          # Database repositories
│   ├── cache/               # Redis implementation
│   └── auth/                # PASETO token management
├── migrator/mysql/sql/      # Database migrations
└── mocks/                   # Generated test mocks
```

### Frontend - Feature-Based Architecture  
```
frontend/src/
├── app/                     # Next.js App Router pages
│   ├── auth/               # Login/register pages
│   ├── dashboard/          # Main dashboard  
│   ├── profile/[id]/       # Person profile pages
│   └── settings/           # User preferences/settings
├── components/
│   ├── ui/                 # shadcn/ui base components
│   ├── auth/               # AuthGuard, auth forms
│   ├── company/            # CompanySelector
│   ├── profile/            # Profile tabs + @mention system
│   └── timeline/           # UnifiedTimeline, FilterBar, Pagination
├── lib/
│   ├── stores/             # Zustand stores (auth, company, people)
│   ├── api/client.ts       # ⭐ Centralized API client 
│   ├── utils/              # storageManager, dates, gender formatting
│   └── types/              # TypeScript definitions
└── hooks/                  # Custom React hooks (useThemePreferences)
```

## Implementation Status

### ✅ Completed Features

**Backend (Go)**
- PASETO authentication (15min access/24h refresh tokens)
- Company/Person/Timeline APIs with full CRUD
- AI System: Chat API, usage tracking, feedback, attribute extraction
- Multi-provider support (OpenAI integrated, Anthropic ready)
- Unified Timeline API with server-side filtering and pagination
- Generic parameter utilities for type-safe request handling
- User preferences with theme persistence

**Frontend (Next.js)**
- Multi-company management with localStorage persistence
- Person profiles with @mention system
- Unified Timeline with advanced filtering
- Dark/light theme with backend sync
- Responsive design with tw-animate-css animations
- Settings page with user preferences

### ⏳ Next Priorities
1. Frontend AI Integration (chat interface)
2. 1:1 Meetings API
3. Member referral system

## AI System Implementation (✅ Complete)

### API Endpoints
- **Chat**: `POST /companies/{uuid}/people/{uuid}/ai/chat` - Leadership coaching
- **Usage Report**: `GET /companies/{uuid}/ai/usage?period={period}` - Token/cost analytics  
- **Feedback**: `POST /companies/{uuid}/ai/usage/{id}/feedback` - Response rating
- **Extraction**: Automatic attribute extraction from notes (background process)

### Architecture
```
Backend AI Flow:
├── Transport Layer: /routes/airoute/ (handlers + validation)
├── Application Layer: /application/service/ai.go (business logic)
├── Domain Layer: /domain/entity/ai.go (entities + interfaces)
├── Infrastructure: /infra/ai/ (openai + anthropic providers)
└── Data Layer: /infra/data/mysql/ai.go (persistence)
```

### Key Features
- **Contextual AI**: Complete person context (attributes, timeline, preferences)
- **Multi-Provider**: OpenAI and Anthropic support with easy extensibility
- **Usage Tracking**: Token counting, cost calculation, usage reports
- **Attribute System**: Flexible key-value storage for person insights
- **Feedback Loop**: Response quality tracking for continuous improvement
- **Background Jobs**: Automatic cleanup of old conversations (30 days)

### Configuration (backend/deployment/config-local.toml)
```toml
[ai.openai]
api_key = "${OPENAI_API_KEY}"  # Environment variable
model = "gpt-4o"
max_tokens = 4096

[ai.anthropic]  
api_key = "${ANTHROPIC_API_KEY}"  # Environment variable
model = "claude-3-5-sonnet-latest"
max_tokens = 4096
```

**Implementation Status**: ✅ Backend complete (Plan 003), ⏳ Frontend next (Plan 004)




## 🚨 CRITICAL DEVELOPMENT RULES - READ BEFORE CODING

### 🔍 COMPONENT CREATION RULES ⚠️ EXTREMAMENTE IMPORTANTE
**SEMPRE procure por componentes compartilhados antes de criar novos componentes para evitar duplicidade!**

**Checklist obrigatório ANTES de criar qualquer componente:**
1. **Verificar `/components/ui/`** - Existe componente similar?
2. **Analisar reutilização** - Este componente será usado em outros lugares?
3. **Local correto** - Se reutilizável, criar IMEDIATAMENTE em `/components/ui/`
4. **Constantes** - Verificar `/lib/constants/` antes de criar valores fixos

**Componentes já disponíveis para reutilização:**
- `LoadingSpinner` - Loading indicators (10+ usos eliminados)
- `AppLogo` - Logo da aplicação (4+ usos eliminados)
- `ErrorMessage` - Mensagens de erro compartilhadas
- `PasswordInput` - Input de senha com toggle visibility
- `PhoneInput` - Input com máscara brasileira
- `SubmitButton` - Botões de envio padronizados
- `MentionsInputComponent` - Sistema de @mentions
- `Pagination` - Componente de paginação customizável
- `EditNoteModal` - Modal para edição de notas
- `FilterBar` - Barra de filtros com quick views

### 🗄️ Storage Manager - USAR SEMPRE ⚠️ CRÍTICO PARA SEGURANÇA
**NUNCA use localStorage diretamente! SEMPRE use o Storage Manager:**

```typescript
// ✅ CORRETO - Uso obrigatório
import { storageManager } from '@/lib/utils/storageManager'
storageManager.set('leaderpro-active-company', companyId)
const companyId = storageManager.get<string>('leaderpro-active-company')

// ❌ INCORRETO - Causa vazamento de dados entre usuários
localStorage.setItem('leaderpro-active-company', companyId)  // NÃO FAÇA ISSO
```

**Por que Storage Manager é obrigatório:**
- **Segurança**: Previne vazamento de dados entre usuários
- **Logout seguro**: `storageManager.clearAll()` limpa TUDO de uma vez
- **Debugável**: `storageManager.debug()` mostra todos os dados
- **Centralizado**: Todas as chaves gerenciadas em um lugar

### ⭐ Frontend API Client - USAR SEMPRE
**NUNCA use fetch() diretamente. SEMPRE use o apiClient centralizado:**

```typescript
// ✅ CORRETO - Uso obrigatório  
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
```

**Por que apiClient é obrigatório:**
- Headers de autenticação (`user-token`) incluídos automaticamente
- Renovação automática de token em caso de 401
- Gerenciamento centralizado de erros
- Configuração única - mudanças de header em 1 lugar só
- Type safety com TypeScript


## Key APIs

### Unified Timeline API
`GET /companies/:company_uuid/people/:person_uuid/timeline`

**Query Parameters:**
- `page`, `quantity` - Pagination (default: page=1, quantity=25)
- `search_query` - Text search across all fields
- `types[]` - Filter: `["feedback", "one_on_one", "observation", "mention"]`
- `feedback_types[]` - Filter: `["positive", "constructive", "neutral"]`
- `direction` - Filter: `"all"`, `"about-person"`, `"from-person"`, `"bilateral"`
- `period` - Filter: `"7d"`, `"30d"`, `"3m"`, `"6m"`, `"1y"`, `"all"`

### AI Chat API
`POST /companies/:company_uuid/people/:person_uuid/ai/chat`
```json
{
  "message": "How can I improve my 1:1s with this person?"
}
```

### AI Usage Report
`GET /companies/:company_uuid/ai/usage?period={period}`

### Feedback API
`POST /companies/:company_uuid/ai/usage/:id/feedback`
```json
{
  "rating": 5,
  "comment": "Very helpful suggestion"
}
```


## Database Migrations

**Location**: `backend/migrator/mysql/sql/`

**Current Migrations:**
- 000001: Initial setup (users, sessions, companies, persons)
- 000002: Notes and mentions system
- 000003: Gender field for contextual text
- 000004: Address management
- 000005: User preferences with theme
- 000006: AI infrastructure (attributes, prompts, usage, conversations)
- 000007: AI conversation cleanup job

Migrations run automatically on backend startup.



