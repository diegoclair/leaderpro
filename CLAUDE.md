# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LeaderPro** - AI platform that amplifies leadership intelligence, maintaining perfect memory of every team interaction and suggesting contextual actions.

**Tagline:** "Become a smarter leader"

## Architecture & Tech Stack

### Backend (Go) - Clean Architecture + DDD
- **Go 1.24.5** with Echo v4.13.3 framework
- **Database**: MySQL 8.0.32 with GORM ORM  
- **Cache**: Redis 7.4.2 for sessions and caching
- **Auth**: PASETO tokens (15min access, 24h refresh) via `user-token` header
- **Testing**: Testcontainers + go-sqlmock for unit/integration tests
- **Observability**: Prometheus, Grafana, Jaeger tracing
- **Architecture**: Domain-Driven Design with Clean Architecture principles
- **CI/CD**: GitHub Actions with coverage reports via Coveralls

### Frontend (Next.js) - ✅ FULLY IMPLEMENTED  
- **Next.js 15.3.5** with App Router and React 19.0.0
- **TypeScript** with strict mode, path alias `@/*` → `./src/*`
- **Styling**: TailwindCSS v4 + shadcn/ui components
- **State**: Zustand stores with localStorage persistence
- **Icons**: Lucide React, **Date Utils**: date-fns
- **Deploy**: GitHub Pages via GitHub Actions (auto-deploy on main push) - **TEMPORARY SETUP**
- **Animation**: Performance-optimized animation system with tw-animate-css

## Development Commands

### Using Tilt (Recommended for Full Stack Development)
```bash
# Start all services with hot reload (frontend + backend + infra)
tilt up

# Stop all services
tilt down

# View Tilt UI
# Open http://localhost:10350 in browser
```

### Frontend (Next.js)
```bash
cd frontend
npm run dev        # Development server with Turbopack (http://localhost:3000)
npm run build      # Production build with static export
npm run lint       # ESLint linting
npx tsc --noEmit   # TypeScript type checking
```

### Backend (Go)
```bash
cd backend
make start         # Start all services (MySQL, Redis, Go app) via Docker
make tests         # Run all tests with coverage
make mocks         # Generate mocks for testing
make docs          # Generate Swagger API docs
make clean-volumes # Clean Docker volumes if needed

# Individual commands
go test -v -cover ./...                      # Run all tests
go test -v -cover ./internal/domain/...      # Run specific package tests
go test -v -run TestFunctionName ./path/to/package  # Run single test
govulncheck ./...                            # Check for vulnerabilities
```

**Note**: The user typically handles `make start` or `tilt up` manually, so don't run it automatically.

### Key URLs After Starting
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:5000  
- **Swagger Docs**: http://localhost:5000/swagger/
- **Live Demo**: https://diegoclair.github.io/leaderpro/ (temporary)
- **Tilt UI**: http://localhost:10350

### Service Ports
- **MySQL**: 3306
- **Redis**: 6379  
- **Prometheus**: 9090 (http://localhost:9090)
- **Grafana**: 3001 (http://localhost:3001)
- **Jaeger**: 16686 (http://localhost:16686)

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
**Frontend (Next.js)**: Complete implementation with all core features
- Multi-company management with persistence
- Person profile management with @mention system  
- 1:1 meeting system with notes and AI suggestions (mocked)
- **Unified Timeline System** - Complete timeline with server-side filtering, search, and pagination
- Smart date formatting (days/months/years ago) with exact date tooltips
- Dark/light theme toggle + responsive design
- **Settings Page** - User preferences with backend-persisted theme
- **Animation System** - Performance-optimized with tw-animate-css

**Backend (Go)**: Authentication + company + person + timeline + **AI system** implemented
- PASETO-based JWT authentication (15min/24h tokens)
- User registration/login with automatic token refresh
- Company CRUD operations with user ownership model
- **Person Management API** - Complete CRUD endpoints for people
- **Unified Timeline API** - Server-side filtering, search, and pagination for all person activities
- **Generic Parameter Utilities** - Type-safe parameter parsing with Go generics
- **User Preferences API** - Theme persistence and user settings
- **🤖 AI System Complete** - Chat API, attribute extraction, usage tracking, feedback system
- **🏗️ Provider Architecture** - OpenAI integrated, ready for Claude/local models
- **📊 AI Analytics** - Token usage, cost tracking, performance metrics
- **Address Management** - Person addresses (migration 000004)
- MySQL integration with Clean Architecture
- Onboarding flow (frontend wizard → backend company creation)

### ⏳ Next Priorities
1. **Frontend AI Integration** - Chat interface and AI components (see `/plan/000004-frontend-ai-integration.md`)
2. **1:1 Meetings API** - Backend endpoints for meeting notes
3. **Member Get Member System** - Referral tracking with discounts

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
├── Infrastructure: /infra/ai/openai/ (provider implementation)
└── Data Layer: /infra/data/mysql/ai.go (persistence)
```

### Key Features
- **Contextual AI**: Complete person context (attributes, timeline, preferences)
- **Provider Pattern**: Easy to swap OpenAI/Claude/local models
- **Usage Tracking**: Token counting, cost calculation, usage reports
- **Attribute System**: Flexible key-value storage for person insights
- **Feedback Loop**: Response quality tracking for continuous improvement

**Implementation Status**: ✅ Backend complete (Plan 003), ⏳ Frontend next (Plan 004)

## Member Get Member Strategy

### Concept (Smart Cash Flow Strategy)
- **For Referred User**: 10% discount only on first month (R$ 44.91)
- **For Referrer**: 50% discount starting from SECOND month of payment
- **Trigger**: Referred user completes trial + pays first full month
- **Accumulative**: 10 referrals = 10 months at R$ 24.95 instead of R$ 49.90
- **No Limits**: Users can accumulate unlimited discounts

### Strategic Advantages
- **Healthy Cash Flow**: First month always guarantees real revenue (R$ 44.91 minimum)
- **User Qualification**: Filters "discount hunters", attracts genuinely interested users
- **Sustainability**: Program can scale without compromising margins
- **Fraud Prevention**: Deferred discount reduces fake referral attempts
- **Higher Retention**: Both referrer and referred have skin in the game

### Implementation Requirements
1. **Unique Referral Codes**: Each user gets a personalized referral link
2. **Tracking System**: Track clicks, registrations, and conversions
3. **Smart Billing**: 10% referred (1st month) + 50% referrer (2nd month+)
4. **Payment Validation**: Credit only activates after first payment clears
5. **Dashboard**: User interface to track referrals and credits
6. **Gamification**: Achievement badges and progress tracking
7. **Templates**: Pre-written messages for sharing (WhatsApp, LinkedIn, Email)
8. **Fraud Prevention**: Same IP, card, device detection

## AI Context Pipeline (To Implement)
1. **Data Collection**: User inputs, profiles, calendar
2. **Vector Embedding**: Convert to embeddings
3. **Context Building**: Combine temporal + geographic + personal
4. **Prompt Engineering**: Generate suggestions
5. **Response Caching**: Redis for frequent queries

**Note**: Detailed implementation plan available at `/plan/000003-ai-implementation-plan.md`

## Testing Strategy

### Backend Testing
- **Unit Tests**: `go test -v -cover ./...` for all packages
- **Specific Tests**: `go test -v -cover ./internal/domain/...` for targeted testing
- **Single Test**: `go test -v -run TestName ./path/to/package`
- **Integration Tests**: Testcontainers automatically spins up MySQL/Redis in Docker
- **Mock Generation**: `make mocks` generates mocks using uber/mock
- **Test Coverage**: Full coverage expected with `-cover` flag
- **Vulnerability Check**: `govulncheck ./...` for security scanning
- **CI/CD**: GitHub Actions runs tests on every push with Coveralls integration

### Frontend Testing  
- **Linting**: `npm run lint` (ESLint + Next.js rules)
- **Type Checking**: `npx tsc --noEmit` (strict TypeScript)
- **Development**: `npm run dev` with Turbopack hot reload

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

### ⚡ Constantes Centralizadas - USAR SEMPRE
**SEMPRE verificar `/lib/constants/` antes de criar valores fixos:**

```typescript
// ✅ CORRETO - Use constantes centralizadas
import { API_ENDPOINTS } from '@/lib/constants/api'
import { COMPANY_SIZES } from '@/lib/constants/company'  // Padrões brasileiros
import { MESSAGES } from '@/lib/constants/messages'
import { VALIDATION } from '@/lib/constants/validation'
import { 
  NOTE_SOURCE_TYPES, 
  getNoteSourceTypeLabel, 
  getFeedbackTypeColor 
} from '@/lib/constants/notes'  // Types de notas e feedbacks

// ❌ INCORRETO - Valores hardcoded espalhados
const endpoint = '/companies'  // NÃO FAÇA ISSO
if (note.type === 'feedback') // NÃO FAÇA ISSO - use NOTE_SOURCE_TYPES.FEEDBACK
```

### Backend Authentication
- **Header**: `user-token` (NOT `Authorization`)
- **Middleware**: `serverMiddleware/auth.go` validates protected routes  
- **Tokens**: PASETO with 15min access + 24h refresh
- **Session Storage**: Redis with user context

### Go Clean Architecture Pattern
```go
// Flow: transport → application → domain ← infra
// Repository interface (domain layer)
type CompanyRepository interface {
    CreateCompany(ctx context.Context, company entity.Company) (int64, error)
    GetCompaniesByUser(ctx context.Context, userID int64) ([]entity.Company, error)
}

// Service (application layer)  
func (s *CompanyService) CreateCompany(ctx context.Context, company entity.Company) (entity.Company, error) {
    userID, err := s.authApp.GetLoggedUserID(ctx) // Get from auth context
    company.UserOwnerID = userID
    return s.dm.Company().CreateCompany(ctx, company)
}
```

## Unified Timeline API

### Endpoint
**Route**: `GET /companies/:company_uuid/people/:person_uuid/timeline`
**Description**: Unified timeline combining direct notes and mentions with server-side filtering

### Query Parameters (All Optional)
- `page` - Page number for pagination (default: 1)
- `quantity` - Items per page (default: 25, max: 50)
- `search_query` - Text search across content, author, category, type
- `types[]` - Filter by activity types: `["feedback", "one_on_one", "observation", "mention"]`
- `feedback_types[]` - Filter by feedback type: `["positive", "constructive", "neutral"]`
- `direction` - Filter by direction: `"all"`, `"about-person"`, `"from-person"`, `"bilateral"`
- `period` - Filter by time period: `"7d"`, `"30d"`, `"3m"`, `"6m"`, `"1y"`, `"all"`

### Response Format
```json
{
  "data": [
    {
      "uuid": "note-uuid-123",
      "type": "feedback",
      "content": "Great performance on {{person:uuid|John}} project",
      "author_name": "Manager Name",
      "created_at": "2024-01-15T10:30:00Z",
      "feedback_type": "positive",
      "feedback_category": "performance",
      "person_name": "John Doe",
      "source_person_name": null,
      "entry_source": "direct"
    },
    {
      "uuid": "mention-uuid-456", 
      "type": "mention",
      "content": "{{person:target-uuid|Target}} was mentioned in this feedback",
      "author_name": "Another Manager",
      "created_at": "2024-01-14T15:20:00Z",
      "feedback_type": "constructive",
      "person_name": "Source Person",
      "source_person_name": "Target Person",
      "entry_source": "mention"
    }
  ],
  "pagination": {
    "page": 1,
    "quantity": 25,
    "total_records": 150,
    "total_pages": 6
  }
}
```

### Key Benefits
- **Single Request**: Eliminates need for separate `/timeline` and `/mentions` calls
- **Server-Side Filtering**: Proper filtering across all data, not just current page
- **Better Performance**: Reduced network requests and client-side processing
- **Unified Sorting**: Chronological order across all activity types
- **Traditional Pagination**: User-controlled page size (10, 25, 50) with page navigation
- **Scalable**: Handles large datasets with efficient database queries

### Pagination Features
- **Page Size Selection**: 10, 25, or 50 items per page
- **Smart Navigation**: Previous/Next + numbered pages with ellipsis
- **Search Debounce**: 800ms delay for search queries to avoid excessive requests
- **State Management**: Page resets to 1 when filters change
- **User Feedback**: "Showing X to Y of Z records" information

## Key Entry Points & Files

### Backend Critical Files
- `backend/cmd/main.go` - Application startup
- `backend/internal/domain/entity/` - Business entities (User, Company, Person, OneOnOne)
- `backend/internal/application/service/` - Business logic services
- `backend/internal/transport/rest/routes/` - HTTP route handlers
- `backend/internal/transport/rest/routeutils/request.go` - ⭐ **Generic parameter utilities** (use sempre!)
- `backend/migrator/mysql/sql/` - Database migrations
- `backend/Tiltfile` - Tilt orchestration configuration

### Frontend Critical Files  
- `frontend/src/app/` - Next.js App Router pages (auth, dashboard, profile, settings)
- `frontend/src/lib/stores/authStore.ts` - ⭐ **apiClient centralizado** (use sempre!)
- `frontend/src/lib/utils/storageManager.ts` - ⭐ **Storage Manager** (use sempre!)
- `frontend/src/lib/utils/dates.ts` - ⭐ **Smart date formatting** (use sempre!)
- `frontend/src/lib/utils/gender.ts` - **Gender-aware Portuguese formatting**
- `frontend/src/lib/stores/` - Zustand state stores (auth, company, people)
- `frontend/src/components/ui/` - ⭐ **Componentes compartilhados** (verifique primeiro!)
- `frontend/src/components/timeline/` - ⭐ **Timeline components** (UnifiedTimeline, FilterBar, etc.)
- `frontend/src/lib/constants/` - ⭐ **Constantes centralizadas** (use sempre!)
- `frontend/src/lib/types/index.ts` - TypeScript type definitions
- `frontend/src/hooks/useThemePreferences.ts` - Theme persistence hook
- `frontend/Tiltfile` - Tilt orchestration configuration

## Environment & Configuration

### Prerequisites
- **Docker** (backend services)
- **Go 1.24.5+** (backend development)
- **Node.js 18+** (frontend development)
- **Tilt** (optional, for orchestrated development)

### Configuration
- **Backend Config**: `backend/deployment/config-local.toml`  
- **Frontend ENV**: `NEXT_PUBLIC_API_URL` (defaults to http://localhost:5000)
- **Default Ports**: Frontend (3000), Backend (5000), MySQL (3306), Redis (6379)
- **Docker Services**: `backend/docker-compose.yml` orchestrates MySQL, Redis, Prometheus, Grafana, Jaeger
- **Tilt Config**: Root `Tiltfile` includes both backend and frontend configurations

### Deployment
- **Frontend**: Auto-deploys to GitHub Pages on push to main branch via GitHub Actions (**TEMPORARY**)
- **Build Output**: `frontend/out` directory (Next.js static export)
- **CI/CD**: Uses Node.js 18 in GitHub Actions pipeline
- **Backend CI**: GitHub Actions runs tests with coverage on every push

## Database Migrations

### Migration Files Location
- **Path**: `backend/migrator/mysql/sql/`
- **Naming**: Sequential numbering (000001_, 000002_, etc.)

### Current Migrations
- `000001_initial_setup.sql` - Users, sessions, companies, persons
- `000002_create_note_tables.sql` - Notes and mentions system  
- `000003_add_gender_to_person.sql` - Gender field for contextual text
- `000004_create_address_table.sql` - Address management
- `000005_create_user_preferences.sql` - User preferences with theme support

### Running Migrations
Migrations run automatically on backend startup. Check `backend/cmd/main.go` for migration logic.

## Recent Technical Improvements

### ✅ Gender System Complete
Full-stack implementation for contextual text in Portuguese:
- Backend: ENUM field with validation
- Frontend: Gender select with Portuguese labels
- Utils: `gender.ts` for contextual formatting ("mencionada" vs "mencionado")

### ✅ Unified Timeline System
- Single API endpoint combining timeline + mentions
- Server-side filtering and search
- Traditional pagination with page controls
- Smart date formatting with tooltips

### ✅ Generic Parameter Utilities (Go)
Type-safe parameter parsing with Go generics:
```go
// Example usage from backend/internal/transport/rest/routeutils/request.go
noteTypes, _ := routeutils.GetArrayParam(c.QueryParam("types"), ",", routeutils.StringConverter)
companyID, _ := routeutils.GetRequiredParam(c.Param("company_id"), routeutils.StringConverter, "company_id required")
includeArchived := routeutils.GetBoolQueryParam(c, "include_archived")
page := routeutils.GetIntQueryParam(c, "page", 1)
```

### ✅ User Preferences with Theme Persistence
- Backend API for user preferences
- Frontend settings page with theme toggle
- `useThemePreferences` hook for backend sync
- Automatic theme application on login

### ✅ Animation System
- Performance-optimized with tw-animate-css
- Documentation at `/frontend/ANIMATION_SYSTEM.md`
- Consistent animation patterns across components

## Documentation References

### Main Documentation Files
- `/README.md` - Project overview and business context
- `/plan/000001-projeto-leaderpro.md` - Complete business plan
- `/plan/000003-ai-implementation-plan.md` - AI implementation strategy
- `/frontend/README.md` - Frontend architecture details
- `/frontend/ANIMATION_SYSTEM.md` - Animation system guide
- `/frontend/STORAGE_MANAGER_DEMO.md` - Storage manager usage examples
- `/backend/README.md` - Backend API documentation
- `/github-pages.md` - Temporary GitHub Pages setup guide

## Common Development Workflows

### Full Stack Development with Tilt
```bash
# Start everything with hot reload
tilt up

# Check service status at http://localhost:10350
# Frontend changes auto-reload
# Backend changes trigger rebuilds
```

### Running a Single Backend Test
```bash
cd backend
go test -v -run TestCreateCompany ./internal/application/service/
```

### Checking Frontend Type Errors
```bash
cd frontend
npx tsc --noEmit  # Check for TypeScript errors without building
```

### Viewing API Documentation
After running `make start` in backend or `tilt up`:
- Swagger UI: http://localhost:5000/swagger/

### Clearing All Docker Data (Full Reset)
```bash
cd backend
make clean-volumes  # Removes all Docker volumes
make start         # Restart with fresh database
# OR
tilt up            # If using Tilt
```

## Business Context

### Product Vision
**LeaderPro** helps tech leads become better managers by maintaining perfect memory of team interactions and providing AI-powered contextual suggestions. Target market: 4.4M new tech leads promoted annually without management training.

### Pricing Model
- **Basic**: R$ 29.90/month (50 people)
- **Standard**: R$ 49.90/month (200 people) - Main offering
- **Unlimited**: R$ 79.90/month (unlimited people)
- **Trial**: 30 days free
- **Early Adopter**: 6 months for R$ 9.90

### AI Integration Strategy
The AI assistant will be the key differentiator, acting as a 24/7 leadership coach. Implementation uses a flexible attribute system (`person_attributes` table) to store dynamic person data that feeds into AI context generation. See `/plan/000003-ai-implementation-plan.md` for detailed implementation strategy.