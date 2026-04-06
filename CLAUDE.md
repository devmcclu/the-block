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
- **File-based routing** via `vue-router/auto-routes` — pages live in `frontend/src/pages/` (e.g., `index.vue` → `/`, `about.vue` → `/about`)
- **Tailwind CSS v4** with `@tailwindcss/vite` plugin
- **shadcn-vue** (new-york style, reka base, geist font) — UI components go in `src/components/ui/`, utils in `src/lib/utils.ts`. MCP server configured in `.mcp.json` for component installation
- **Icons:** lucide-vue-next (hugeicons configured as shadcn icon library)
- **Testing:** Vitest + jsdom + @vue/test-utils. Tests go in `src/components/__tests__/`
- **Linting:** oxlint (primary) → eslint (secondary, with vue + typescript + vitest plugins) → prettier for formatting
- **Path alias:** `@` → `frontend/src/`

### Backend (`backend/`)
- **Go** with **Fuego** web framework (auto-generates OpenAPI spec at `backend/doc/openapi.json`)
- Domain-driven structure: `backend/domains/<domain>/` with controller + service pattern
- Currently has a `cars` domain with CRUD scaffold (stub implementations)
- Fuego uses typed handlers: `ContextNoBody` for reads, `ContextWithBody[T]` for writes

### Formatting Conventions
- Frontend: no semicolons, single quotes, 100 char print width (prettier)
- Backend: gofumpt
