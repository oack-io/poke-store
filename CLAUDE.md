# Production

- This is a demo Pokémon store for Playwright testing. No real production environment.
- Default environment is localhost:3000.

# Git

- Only use git commands scoped to this repository (`/Users/gkomissarov/Documents/vcs/public/poke-store`). Never operate on repositories outside this folder.

# Architecture

- **Backend**: Go 1.26, internet-facing (no reverse proxy). Serves both the API and the Astro-built static frontend.
- **Frontend**: Astro (static site generation). Built artifacts in `web/dist/` are embedded into the Go binary via `embed.FS`.
- **Data**: In-memory. No database. Pokemon catalog and user sessions live in Go structs.
- **Auth**: Cookie-based sessions with HMAC-signed tokens. Demo users hardcoded in `internal/data/users.go`.

# ProjectFlow

## Post-feature checklist

- After completing any feature, run `make lint` to ensure the code passes all linters.
- Run `make test` to verify nothing is broken.
- Run `make build` to verify the full build (Go + Astro) succeeds.
- If there are Playwright test scenarios affected, update `tests/e2e/`.

## Go style

- **Max line length: 120 characters.** Break longer lines.
- **Parameter structs over long signatures.** When a function has more than 5 parameters, use a named params struct.
- **Concise error wrapping.** Use short noun phrases in `fmt.Errorf`: `fmt.Errorf("create session: %w", err)`. Avoid verbose prefixes like "failed to", "unable to", "could not".
- **Prefer strconv over fmt for primitives.** Use `strconv.Itoa(n)` instead of `fmt.Sprintf("%d", n)`.
- **Don't embed mutexes.** Always use a named field (`mu sync.Mutex`), never embed `sync.Mutex` directly.
- **Reduce nesting with early returns.** Handle error/special cases first and return early.
- **Safe type assertions.** Always use the two-value form `v, ok := x.(T)`.

## Logging conventions

- All logs use `slog` with JSON output format and UTC timestamps.
- Every `slog` call must include a `"scope"` attribute as the first key-value pair. Scope values:
  - `http` — HTTP handlers, middleware, request/response logging
  - `auth` — authentication, session management
  - `store` — product catalog, search
  - `cart` — cart operations
  - `lifecycle` — server startup/shutdown
