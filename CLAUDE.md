# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Vehicle auction platform (buyer side) built for the OPENLANE coding challenge. Buyers browse inventory, inspect vehicle details, and place bids on 200 vehicles from `data/vehicles.json`.

## Commands

### Run everything (backend + frontend in parallel)
```bash
npm run dev          # or: npm start
```

### Frontend only (from repo root)
```bash
cd frontend && npm run dev
```

### Backend only (from repo root)
```bash
cd backend && air -c .air.toml
```

### Linting & formatting (frontend)
```bash
cd frontend
npm run lint              # runs oxlint then eslint (both with --fix)
npm run lint:oxlint       # oxlint only
npm run lint:eslint       # eslint only
npm run format            # prettier
```

### Testing (frontend)
```bash
cd frontend
npm run test:unit         # vitest (watch mode)
npx vitest run            # single run
npx vitest run src/components/__tests__/HelloWorld.spec.ts  # single file
```

### Build & type-check (frontend)
```bash
cd frontend
npm run build             # type-check + vite build
npm run type-check        # vue-tsc only
```

### Backend linting
```bash
cd backend && golangci-lint run
```

## Tooling

- **mise** manages tool versions: Node 24, npm 11, Go 1.26, air, golangci-lint, gofumpt (see `mise.toml`)
- **air** provides Go hot-reload for the backend (config: `backend/.air.toml`)

## Architecture

### Frontend (`frontend/`)
- **Vue 3** + **TypeScript** + **Vite 8** + **Pinia** for state
- **File-based routing** via `vue-router/auto-routes` — pages live in `frontend/src/pages/` (e.g., `index.vue` → `/`, `vehicles/[id].vue` → `/vehicles/:id`). No manual route config needed.
- **Tailwind CSS v4** with `@tailwindcss/vite` plugin — all styling via utility classes, no component-scoped `<style>` blocks
- **shadcn-vue** (new-york style, reka base, geist font) — UI primitives in `src/components/ui/`, utils in `src/lib/utils.ts`. MCP server configured in `.mcp.json` for component installation
- **Icons:** iconify-vue with hugeicons library (`@iconify-json/hugeicons`)
- **API client:** `openapi-fetch` in `src/lib/api/client.ts`, typed from generated `src/lib/api/v1.d.ts`. Base URL `/api` is proxied by Vite to `localhost:9999`
- **Testing:** Vitest + jsdom + @vue/test-utils. Tests go in `src/components/__tests__/`
- **Linting:** oxlint (primary) → eslint (secondary, with vue + typescript + vitest plugins) → prettier for formatting
- **Path alias:** `@` → `frontend/src/`

### Backend (`backend/`)
- **Go** with **Fuego** web framework (auto-generates OpenAPI spec at `backend/doc/openapi.json`)
- **SQLite + GORM** — database auto-created at `backend/vehicles.db`, auto-migrated on startup, seeded from `data/vehicles.json` if empty
- Domain-driven structure: `backend/domains/<domain>/` with controller + service pattern
- `vehicles` domain: controller.go (HTTP handlers + route registration), service.go (interface + GORM implementation)
- Fuego uses typed handlers: `ContextNoBody` for reads, `ContextWithBody[T]` for writes
- Vehicles use `external_id` (UUID) as public API identifier; internal GORM `id` is never exposed
- Optional fields (`ReservePrice`, `BuyNowPrice`) use Go pointer types (`*int`) to distinguish null from zero

### API Routes
```
GET    /vehicles/         # List with filters + sorting (query params)
GET    /vehicles/config   # Returns { max_auction_duration_hours: 720 }
GET    /vehicles/filters  # Returns distinct values for all filterable fields
GET    /vehicles/{id}     # Single vehicle by external_id
POST   /vehicles/         # Create
PUT    /vehicles/{id}     # Update (CurrentBid, BidCount)
DELETE /vehicles/{id}     # Delete
```

### Key Conventions
- **Frontend state:** All shared state lives in Pinia stores (`src/stores/`), not component state. Filters use `useDebounceFn` (300ms) to avoid excessive API calls
- **API types:** Use `components["schemas"]["Vehicle"]` from `v1.d.ts` — don't manually define API types
- **Parameter naming:** Backend uses snake_case in JSON; frontend store uses camelCase, converts when building query params in `fetchVehicles()`
- **Auction timing:** `auction_start` (ISO 8601) + server-configured `max_auction_duration_hours` (720h = 30 days). `useAuctionTime` composable handles countdown logic
- **Dark mode:** Handled by `@vueuse/core`'s `useColorMode` + Tailwind `dark:` prefix

### Formatting Conventions
- Frontend: semicolons, double quotes, 100 char print width (prettier)
- Backend: gofumpt
