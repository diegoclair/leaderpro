# API Client Documentation

This directory contains the API client implementation for LeaderPro frontend.

## Files Overview

- `client.ts` - Base API client with authentication handling
- `typed-client.ts` - Strongly typed wrapper for common API calls 
- `../types/api.ts` - API response type definitions

## Usage

### Basic API Client (client.ts)

The base API client provides generic methods for making HTTP requests:

```typescript
import { apiClient } from '@/lib/api/client'

// Public requests (no authentication)
const loginData = await apiClient.post<LoginResponse>('/auth/login', { email, password })

// Authenticated requests (automatic token handling + refresh)
const profile = await apiClient.authGet<UserProfileResponse>('/users/profile')
const companies = await apiClient.authGet<Company[]>('/companies')
```

### Typed API Client (typed-client.ts)

For better type safety and convenience, use the typed API client:

```typescript
import { typedApi } from '@/lib/api/typed-client'

// Authentication
const { user, tokens } = await typedApi.login(email, password)
await typedApi.logout()

// Companies
const companies = await typedApi.getCompanies()
const { company } = await typedApi.createCompany(companyData)

// People
const people = await typedApi.getPeople(companyUuid)
await typedApi.updatePerson(companyUuid, personUuid, personData)

// Dashboard
const { people, stats } = await typedApi.getDashboard(companyUuid)

// Timeline
const { data, pagination } = await typedApi.getTimeline(companyId, personUuid, params)
```

## Type Safety Improvements

### Before (using 'any')
```typescript
// ❌ No type safety
const response = await apiClient.authGet('/companies') // response: any
const companies = response // companies: any
```

### After (using proper types)
```typescript
// ✅ Full type safety
const companies = await typedApi.getCompanies() // companies: Company[]
const { people, stats } = await typedApi.getDashboard(companyUuid) // Destructured with proper types
```

## API Response Types

All API responses are properly typed in `/lib/types/api.ts`:

- `LoginResponse` - Login/register responses with user and tokens
- `CompaniesResponse` - Company listing
- `DashboardResponse` - Dashboard data with people and stats
- `TimelineResponse` - Timeline with pagination
- `VoidResponse` - For operations that don't return data

## Authentication Handling

The API client automatically handles:

- Adding `user-token` header to authenticated requests
- Token refresh on 401 responses
- Error notifications (except for login/register)
- Session expiration with user-friendly messages

## Migration Guide

To migrate existing code:

1. **Replace generic apiClient calls with typed alternatives:**
   ```typescript
   // Before
   const response = await apiClient.authGet('/companies')
   
   // After  
   const companies = await typedApi.getCompanies()
   ```

2. **Use proper types for responses:**
   ```typescript
   // Before
   const response: any = await apiClient.authGet('/dashboard')
   
   // After
   const { people, stats }: DashboardResponse = await typedApi.getDashboard(companyUuid)
   ```

3. **Import types from centralized location:**
   ```typescript
   import type { Company, Person, TimelineResponse } from '@/lib/types'
   ```

## Benefits

- **Type Safety**: Catch errors at compile time
- **Better IDE Support**: Autocomplete and inline documentation
- **Consistent Error Handling**: Centralized error management
- **Automatic Token Management**: No manual token handling needed
- **Maintainability**: Changes to API responses update types everywhere