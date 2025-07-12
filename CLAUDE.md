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

### Frontend (Next.js) - ✅ FULLY IMPLEMENTED  
- **Next.js 15.3.5** with App Router and React 19.0.0
- **TypeScript** with strict mode, path alias `@/*` → `./src/*`
- **Styling**: TailwindCSS v4 + shadcn/ui components
- **State**: Zustand stores with localStorage persistence
- **Icons**: Lucide React, **Date Utils**: date-fns
- **Deploy**: GitHub Pages via GitHub Actions (auto-deploy on main push)

## Development Commands

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
go test -v -cover ./internal/domain/...  # Run specific test
```

**Note**: The user typically handles `make start` manually, so don't run it automatically.

### Key URLs After Starting
- **Frontend**: http://localhost:3000
- **Backend API**: http://localhost:5000  
- **Swagger Docs**: http://localhost:5000/swagger/
- **Live Demo**: https://diegoclair.github.io/leaderpro/

## Architecture Overview

### Backend - Clean Architecture Flow
**Dependency Rule**: `transport` → `application` → `domain` ← `infra`

```
backend/
├── cmd/main.go               # App entry point
├── internal/
│   ├── domain/entity/        # Business entities (User, Company, Person)
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
│   └── profile/[id]/       # Person profile pages
├── components/
│   ├── ui/                 # shadcn/ui base components
│   ├── auth/               # AuthGuard, auth forms
│   ├── company/            # CompanySelector
│   └── profile/            # Profile tabs + @mention system
├── lib/
│   ├── stores/             # Zustand stores (auth, company, people)
│   ├── api/client.ts       # ⭐ Centralized API client 
│   └── types/              # TypeScript definitions
└── hooks/                  # Custom React hooks
```

## Implementation Status

### ✅ Completed Features
**Frontend (Next.js)**: Complete implementation with all core features
- Multi-company management with persistence
- Person profile management with @mention system  
- 1:1 meeting system with notes and AI suggestions (mocked)
- Feedback tracking (direct + mentioned across people)
- Dark/light theme toggle + responsive design

**Backend (Go)**: Authentication + company + person management implemented
- PASETO-based JWT authentication (15min/24h tokens)
- User registration/login with automatic token refresh
- Company CRUD operations with user ownership model
- **Person Management API** - Complete CRUD endpoints for people
- MySQL integration with Clean Architecture
- Onboarding flow (frontend wizard → backend company creation)

### ⏳ Next Priorities
1. **1:1 Meetings API** - Backend endpoints for meeting notes
2. **AI Integration** - OpenAI/Claude API for contextual suggestions
3. **Member Get Member System** - Referral tracking with discounts
4. **Advanced Person Features** - Manager relationships, search, filtering

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

## Critical Development Patterns

### ⭐ Frontend API Client - ALWAYS USE
**NEVER use fetch() directly. Always use the centralized apiClient:**

```typescript
import { apiClient } from '@/lib/api/client'

// ✅ Public requests (no auth)
await apiClient.post('/auth/login', { email, password })
await apiClient.post('/users', userData)

// ✅ Authenticated requests (auto token + refresh)
await apiClient.authGet('/users/profile')
await apiClient.authPost('/companies', companyData)
await apiClient.authPut('/users/profile', updates)

// ❌ NEVER do this
// fetch('/api/endpoint') // DON'T DO THIS
```

**Why apiClient is mandatory:**
- Automatic `user-token` header injection
- Automatic token refresh on 401 errors  
- Centralized error handling + notifications
- TypeScript type safety

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

## Key Entry Points & Files

### Backend Critical Files
- `backend/cmd/main.go` - Application startup
- `backend/internal/domain/entity/` - Business entities (User, Company, Person, OneOnOne)
- `backend/internal/application/service/` - Business logic services
- `backend/internal/transport/rest/routes/` - HTTP route handlers
- `backend/migrator/mysql/sql/` - Database migrations

### Frontend Critical Files  
- `frontend/src/app/` - Next.js App Router pages (auth, dashboard, profile)
- `frontend/src/lib/api/client.ts` - ⭐ **Centralized API client** (use always!)
- `frontend/src/lib/stores/` - Zustand state stores (auth, company, people)
- `frontend/src/components/` - UI components organized by feature
- `frontend/src/lib/types/index.ts` - TypeScript type definitions

## Environment & Configuration

### Prerequisites
- **Docker** (backend services)
- **Go 1.24.5+** (backend development)
- **Node.js 18+** (frontend development)

### Configuration
- **Backend Config**: `backend/deployment/config-local.toml`  
- **Frontend ENV**: `NEXT_PUBLIC_API_URL` (defaults to http://localhost:5000)
- **Default Ports**: Frontend (3000), Backend (5000), MySQL (3306), Redis (6379)