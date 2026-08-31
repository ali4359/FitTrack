# Iron & Spice

Gym + regional nutrition app for the Pakistani / South Asian market. Closes the
loop between training and eating: finish a workout → get a MET-based calorie
estimate → get a region- and budget-aware meal suggestion for your goal.

## Monorepo layout

```
iron-and-spice/
├── apps/
│   ├── mobile/     React Native (Expo) app  ← the frontend
│   └── backend/    Go API (Gin + GORM)      ← standalone Go module
├── packages/
│   └── shared/     TypeScript types shared by the mobile app (mirror the Go structs)
├── pnpm-workspace.yaml
└── package.json
```

Go and Node don't share a package manager, so this is just folders — no
Turborepo/Nx. `pnpm` workspaces cover `apps/mobile` + `packages/*` only;
`apps/backend` is a normal Go module.

## Prerequisites

- Node ≥ 20 and `pnpm` (`npm i -g pnpm`)
- Go ≥ 1.24
- Docker (optional — only if you want Postgres instead of the SQLite fallback)
- Xcode / Android Studio, or the Expo Go app on a phone

## Run it locally

### 1. Backend

```bash
cd apps/backend
go run ./cmd/api            # zero-setup: local SQLite file + seeded demo data
```

Or against Postgres:

```bash
cd apps/backend
docker compose up -d
DATABASE_URL='postgres://postgres:postgres@localhost:5432/ironandspice?sslmode=disable' \
  go run ./cmd/api
```

API listens on `http://localhost:8080`. Demo login: `demo@ironandspice.app` / `password123`.
See [apps/backend/README.md](apps/backend/README.md) for the full endpoint list.

> The backend here is a **stub** scaffolded so the mobile app has something to
> build against. Swap in the real Go service — keep routes and JSON shapes
> identical to `packages/shared`.

### 2. Mobile

```bash
pnpm install                          # from the repo root, once
cp apps/mobile/.env.example apps/mobile/.env
# edit apps/mobile/.env → API_BASE_URL (use your LAN IP for a physical device)
pnpm --filter mobile start            # Expo dev server
```

Then press `i` / `a` in the Expo CLI, or scan the QR code with Expo Go.

## Handy scripts (repo root)

| Command | What |
| --- | --- |
| `pnpm mobile` | start the Expo dev server |
| `pnpm backend` | `go run ./cmd/api` |
| `pnpm shared:check` | typecheck `packages/shared` |
| `pnpm --filter mobile typecheck` | typecheck the app |

## Tech

- **Mobile:** Expo (SDK 57), React Navigation (native stack + bottom tabs),
  TanStack Query for server state, Zustand for local/UI state, axios API client
  with JWT attach + 401 → sign-out, `expo-secure-store` for the token.
- **Design system:** dark UI only. Tokens live in
  [apps/mobile/src/theme](apps/mobile/src/theme). Oswald (display), Inter (body),
  JetBrains Mono (all numeric data).
- **Backend:** Go, Gin, GORM, JWT. Postgres in prod, SQLite fallback for local.

## Build status

Screens 1–4 built (Auth, Onboarding, Home/Today, Workout Session). Session
Complete / Meal Suggestion / Progress / Profile are wired-up placeholders
pending design sign-off.
