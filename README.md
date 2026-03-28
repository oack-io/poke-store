# PokéStore

A demo Pokémon store web application built for Playwright testing demonstrations. Features login, search, cart, and checkout flows with original Generation I Pokémon.

**Live demo**: [poke-store.oack.io](https://poke-store.oack.io/)

## Demo Accounts

| Name              | Email              | Password    |
|-------------------|--------------------|-------------|
| Ash Ketchum       | ash@pokemon.com    | pikachu123  |
| Misty Waterflower | misty@pokemon.com  | starmie123  |
| Brock Harrison    | brock@pokemon.com  | onix123     |

## Playwright Tests

Run tests against the live site or a local instance:

```bash
cd web
npx playwright install   # first time only
npx playwright test
npx playwright show-report
```

Against localhost:

```bash
BASE_URL=http://localhost:6001 npx playwright test
# or via Makefile:
make test-oack
```

## Oack Browser Monitoring

The `web/tests/e2e/` directory contains standard Playwright tests that double as [Oack](https://oack.io) continuous monitoring checks. No custom format — the same tests you run locally with `npx playwright test` run on the Oack platform unchanged.

### Run on Oack (one-off)

Upload your Playwright project directory for a one-off test run on Oack's browser infrastructure:

```bash
# Install and authenticate
brew install oack-io/tap/oackctl
oackctl login

# Run tests (opens HTML report in browser)
oackctl test --team <TEAM> --monitor <MONITOR> --dir web
```

The output shows pass/fail counts and a link to the full Playwright HTML report.

### Skip repetitive flags with .oackctl.env

Create a `.oackctl.env` file in the project root to avoid passing `--team` and `--monitor` on every command:

```env
OACKCTL_TEAM=<your-team-id>
OACKCTL_MONITOR=<your-monitor-id>
```

Now commands are shorter:

```bash
oackctl test --dir web
oackctl deploy --dir web
```

Every CLI flag maps to an `OACKCTL_` env var (`--team` → `OACKCTL_TEAM`, `--pw-grep` → `OACKCTL_PW_GREP`, etc.). Add `.oackctl.env` to `.gitignore` if it contains team-specific IDs.

### Deploy for continuous monitoring

Deploy the test suite to run on a schedule. Every check produces an HTML report and alerts you on failure:

```bash
oackctl deploy --team <TEAM> --monitor <MONITOR> --dir web
```

The deploy captures your git commit, branch, and who deployed — visible in the dashboard.

### Multi-monitor config

Define multiple check suites in a config file:

```bash
oackctl config-deploy --config oack.config.json
```

See [`oack.config.json`](oack.config.json) for an example.

### What you get

- **Playwright HTML report** — full test breakdown with screenshots and error details
- **Pass/fail alerts** — any test failure = monitor DOWN, alerts via email/Slack/PagerDuty
- **Git metadata** — commit SHA, branch, and deployer tracked per deploy
- **Filters** — run subsets with `--pw-grep`, `--pw-project`, or `--pw-tag`

---

## Development

### Tech Stack

- **Backend**: Go 1.26 (net/http, no frameworks)
- **Frontend**: Astro (static site generation)
- **Data**: In-memory (no database)
- **Auth**: Cookie-based sessions with HMAC tokens

### Quick Start

```bash
# Install frontend dependencies
make install-web

# Build everything (Astro + Go)
make build

# Run the server
make run
# or for development (without building):
make dev
```

The app runs on `http://localhost:6001` by default. Override with `ADDR=:8080`.

### Pages

- `/` — Landing page
- `/login` — Login with demo account quick-fill buttons
- `/store` — Browse, search, and filter 25 Gen I Pokémon; add to cart
- `/cart` — View cart, remove items, proceed to checkout
- `/checkout` — Fill delivery & payment info, place order, see confirmation

### API Endpoints

| Method | Path              | Auth | Description            |
|--------|-------------------|------|------------------------|
| POST   | /api/login        | No   | Authenticate user      |
| POST   | /api/logout       | Yes  | End session            |
| GET    | /api/me           | Yes  | Current user + cart count |
| GET    | /api/pokemon      | No   | List/search Pokémon (`?q=`, `?type=`) |
| GET    | /api/pokemon/{id} | No   | Get single Pokémon     |
| GET    | /api/cart         | Yes  | Get cart details       |
| POST   | /api/cart/add     | Yes  | Add to cart            |
| POST   | /api/cart/remove  | Yes  | Remove from cart       |
| POST   | /api/cart/clear   | Yes  | Clear cart             |
| POST   | /api/checkout     | Yes  | Place order            |
| GET    | /api/version      | No   | Build info             |

### Makefile Targets

| Target       | Description                          |
|--------------|--------------------------------------|
| `build`      | Build Astro frontend + Go binary     |
| `build-web`  | Build Astro frontend only            |
| `build-server` | Build Go binary only               |
| `install-web`| Install npm dependencies             |
| `run`        | Build and run                        |
| `dev`        | Run Go server (expects pre-built static) |
| `test`       | Run Go tests                         |
| `test-e2e`   | Run Playwright E2E tests             |
| `test-oack`  | Run oack Playwright checks           |
| `lint`       | Run golangci-lint                    |
| `clean`      | Remove build artifacts               |

### Project Structure

```
poke-store/
├── cmd/server/          # Go entrypoint + embedded static files
├── internal/
│   ├── data/            # Pokemon catalog + demo users
│   ├── handler/         # HTTP handlers (auth, store, cart, checkout)
│   ├── middleware/       # Auth, logging, recovery middleware
│   ├── model/           # Data models
│   └── store/           # Session + cart in-memory stores
├── web/                 # Astro frontend
│   ├── src/pages/       # Astro pages (login, store, cart, checkout)
│   ├── src/layouts/     # Shared layout
│   ├── public/          # Static assets (Pokemon pixel sprites)
│   └── tests/e2e/       # Playwright test specs
├── scripts/             # Build scripts
├── CLAUDE.md            # AI assistant guidelines
└── Makefile             # Build targets
```
