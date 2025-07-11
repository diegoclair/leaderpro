# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**LeaderPro** - AI platform that amplifies leadership intelligence, maintaining perfect memory of every team interaction and suggesting contextual actions.

**Tagline:** "Become a smarter leader"

## Product Context

- **Problem**: 70% of leaders are promoted without adequate training, resulting in ineffective 1:1s and superficial performance reviews
- **Solution**: Multidimensional contextual AI that combines personal + temporal + geographic + historical data
- **Differentiator**: Virtual coach that remembers everything and suggests actions based on complete context
- **Model**: B2C (R$ 49.90/month) - individual leader pays

## Project Structure

### Current State
- `/backend/` - Go boilerplate with DDD/Clean Architecture (needs LeaderPro implementation)
- `/frontend/` - **✅ IMPLEMENTED** - Complete Next.js 15.3.5 frontend with all core features
- `/plan/` - Complete project planning documentation
- `/.github/workflows/` - GitHub Actions for frontend deployment to GitHub Pages

### Key Documentation
- **`/frontend/README.md`** - Frontend architecture, TypeScript types, components, and stores
- **`/backend/README.md`** - Go boilerplate documentation with Clean Architecture
- **`/plan/000001-projeto-leaderpro.md`** - Business plan, market analysis, and strategy

## Technology Stack

### Backend (Go/Golang)
- **Current**: Go boilerplate with Echo, GORM, MySQL, Redis
- **To implement**: LeaderPro-specific endpoints for:
  - Companies and team member management
  - 1:1 meetings and notes
  - AI integration (OpenAI/Claude API)
  - Vector database for contextual memory

### Frontend (Next.js) - ✅ IMPLEMENTED
- **Framework**: Next.js 15.3.5 with App Router
- **Language**: TypeScript with strict mode
- **Styling**: TailwindCSS v4 + shadcn/ui components
- **State**: Zustand for global state management
- **Icons**: Lucide React
- **Date Utils**: date-fns
- **Deploy**: GitHub Pages via GitHub Actions

## Development Commands

### Frontend (Next.js) - ✅ WORKING
```bash
cd frontend

# Development server (with Turbopack)
npm run dev

# Production build (static export)
npm run build

# Linting
npm run lint

# Type checking (manual - not in package.json)
npx tsc --noEmit
```

### Backend (Go) - BOILERPLATE READY
```bash
cd backend

# Start all services (MySQL, Redis, App)
make start

# Run tests with coverage
make tests

# Generate mocks for testing
make mocks

# Generate API documentation (Swagger)
make docs

# Access Swagger UI (after start)
# http://localhost:5000/swagger/

# Individual operations
docker compose up --build    # Same as make start
go test -v -cover ./...      # Same as make tests
```

### Deployment
```bash
# Frontend deploys automatically via GitHub Actions on push to main branch
# Live at: https://diegoclair.github.io/leaderpro/
```

## Architecture Overview

### Backend Architecture (Clean Architecture)
```
backend/
├── cmd/              # Application entry point
├── internal/         # Core application code
│   ├── domain/       # Business logic, entities, interfaces
│   ├── application/  # Use cases/services
│   └── transport/    # HTTP handlers (Echo framework)
├── infra/            # Infrastructure implementations
│   ├── data/         # Database repositories
│   ├── cache/        # Redis implementation
│   └── logger/       # Logging infrastructure
├── migrator/         # Database migrations
├── docs/             # Generated Swagger docs
└── mocks/            # Generated mocks for testing
```

### Frontend Architecture
```
frontend/
├── app/                    # Next.js 15 App Router
│   ├── page.tsx            # Main dashboard
│   └── profile/[id]/       # Person profile pages
├── components/             
│   ├── ui/                 # shadcn/ui components
│   ├── person/             # Person profiles, cards
│   ├── profile/            # Profile tabs (info, history, feedback, chat)
│   └── layout/             # Header, navigation
├── hooks/                  # Custom React hooks
│   ├── use-mentions.ts     # @mention functionality
│   └── use-create-person.ts# Auto-create from mentions
├── lib/                    
│   ├── stores/             # Zustand state stores
│   ├── types/              # TypeScript definitions
│   └── utils/              # Date, name utilities
└── stores/                 # Zustand stores
    ├── company-store.ts    # Company/project management
    ├── people-store.ts     # Team member profiles
    └── one-on-one-store.ts # 1:1 meetings and notes
```

## Core Features

### ✅ Implemented (Frontend)
- Multi-company/project support with persistence
- Complete person profile management
- 1:1 meeting system with notes
- @mention system with auto-suggestions
- Feedback tracking (direct + mentioned)
- Dark/light theme toggle
- Responsive design

### ⏳ Pending Implementation
1. **Backend API** - Implement LeaderPro endpoints in Go boilerplate
2. **AI Integration** - OpenAI/Claude for contextual suggestions
3. **Authentication** - User login/signup
4. **Vector Database** - AI memory with Pinecone/Weaviate
5. **Real-time Sync** - WebSockets for collaboration
6. **Calendar Integration** - Google/Outlook

## AI Context Pipeline (To Implement)
1. **Data Collection**: User inputs, profiles, calendar
2. **Vector Embedding**: Convert to embeddings
3. **Context Building**: Combine temporal + geographic + personal
4. **Prompt Engineering**: Generate suggestions
5. **Response Caching**: Redis for frequent queries

## Common Patterns

### TypeScript (Frontend)
```typescript
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

### Go (Backend)
```go
// Clean Architecture flow
// transport -> application -> domain -> infra

// Repository interface (domain layer)
type UserRepository interface {
    GetByID(ctx context.Context, id string) (*User, error)
}

// Use case (application layer)
type GetUserUseCase struct {
    userRepo UserRepository
}
```

## Important Notes

- **Current Status**: Frontend completed, backend needs LeaderPro implementation
- **Live Demo**: https://diegoclair.github.io/leaderpro/
- **Data Storage**: Currently localStorage (frontend only)
- **AI Features**: Mocked in frontend
- **Next Steps**: Implement backend API endpoints for LeaderPro features

## Key Files and Entry Points

### Backend Entry Points
- `backend/cmd/main.go` - Application startup
- `backend/internal/transport/rest/server.go` - HTTP server setup
- `backend/internal/domain/entity/` - Core business entities (Person, Company, OneOnOne)
- `backend/migrator/mysql/sql/` - Database migrations

### Frontend Entry Points  
- `frontend/src/app/page.tsx` - Main dashboard
- `frontend/src/app/profile/[id]/page.tsx` - Person profile pages
- `frontend/src/lib/stores/` - Zustand state management
- `frontend/src/components/` - Reusable UI components