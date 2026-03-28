# PokéStore — Oack Browser Checks

Playwright test suite for [poke-store.oack.io](https://poke-store.oack.io).
Can be run locally, in Docker, or published to [Oack](https://oack.io) for continuous monitoring.

## Test coverage

| Check | File | What it tests |
|-------|------|---------------|
| Login | `checks/login.js` | Navigate to login, submit credentials, verify redirect |
| Search | `checks/search.js` | Name search, type filtering |
| Cart | `checks/cart.js` | Add item, view cart, remove item |
| Checkout | `checks/checkout.js` | Add item, fill form, place order, verify confirmation |
| Full flow | `checks/full-flow.js` | Login → search → add to cart → checkout → logout |

The `e2e/poke-store.spec.ts` file contains the same scenarios as a standard Playwright Test suite.

## Run locally

```bash
npm install
npx playwright install chromium
npx playwright test
```

Override the target URL:

```bash
BASE_URL=http://localhost:6001 npx playwright test
```

## Run in Docker

```bash
docker build -t poke-store-checks .
docker run --rm poke-store-checks
```

Mount volumes to get results out:

```bash
docker run --rm \
  -v ./test-results:/checks/test-results \
  -v ./playwright-report:/checks/playwright-report \
  poke-store-checks
```

Override the target URL:

```bash
docker run --rm -e BASE_URL=http://host.docker.internal:6001 poke-store-checks
```

## Publish to Oack via oackctl

### Prerequisites

1. Install `oackctl` (see [oack docs](https://docs.oack.io/cli)).
2. Authenticate: `oackctl auth login`.
3. Have a team and a browser-type monitor ready.

### Set environment variables

Store credentials as team-level secrets so scripts can access them at runtime:

```bash
oackctl env set LOGIN_EMAIL "ash@pokemon.com" --team <TEAM_ID>
oackctl env set LOGIN_PASSWORD "pikachu123" --team <TEAM_ID> --secret
oackctl env set BASE_URL "https://poke-store.oack.io" --team <TEAM_ID>
```

### Test a script

Run a one-off test against the remote browser-checker infrastructure:

```bash
oackctl test \
  --team <TEAM_ID> \
  --monitor <MONITOR_ID> \
  --script checks/full-flow.js
```

With local env overrides (useful for staging):

```bash
oackctl test \
  --team <TEAM_ID> \
  --monitor <MONITOR_ID> \
  --script checks/login.js \
  --env BASE_URL=https://staging.poke-store.oack.io \
  --env LOGIN_EMAIL=ash@pokemon.com \
  --env LOGIN_PASSWORD=pikachu123
```

Add `--json` for machine-readable output:

```bash
oackctl test \
  --team <TEAM_ID> \
  --monitor <MONITOR_ID> \
  --script checks/full-flow.js \
  --json
```

### Deploy scripts to monitors

Create a browser monitor for each check you want to run continuously:

```bash
# Create monitors
oackctl monitors create --team <TEAM_ID> --name "PokéStore Login"    --url https://poke-store.oack.io --type browser
oackctl monitors create --team <TEAM_ID> --name "PokéStore Search"   --url https://poke-store.oack.io --type browser
oackctl monitors create --team <TEAM_ID> --name "PokéStore Cart"     --url https://poke-store.oack.io --type browser
oackctl monitors create --team <TEAM_ID> --name "PokéStore Checkout" --url https://poke-store.oack.io --type browser
oackctl monitors create --team <TEAM_ID> --name "PokéStore Full"     --url https://poke-store.oack.io --type browser
```

Then deploy scripts to each monitor (the monitor runs the script on its configured interval):

```bash
oackctl deploy --team <TEAM_ID> --monitor <LOGIN_MONITOR_ID>    --script checks/login.js
oackctl deploy --team <TEAM_ID> --monitor <SEARCH_MONITOR_ID>   --script checks/search.js
oackctl deploy --team <TEAM_ID> --monitor <CART_MONITOR_ID>     --script checks/cart.js
oackctl deploy --team <TEAM_ID> --monitor <CHECKOUT_MONITOR_ID> --script checks/checkout.js
oackctl deploy --team <TEAM_ID> --monitor <FULL_MONITOR_ID>     --script checks/full-flow.js
```

### Script format requirements

Oack scripts must:

- Export a function via `module.exports` (or `export default`)
- Accept `(page, context)` — Playwright page + oack context
- Use `context.step(name, fn)` for named steps with timing
- Use `context.screenshot(name)` to capture screenshots
- Access env vars via `context.KEY` (e.g., `context.LOGIN_EMAIL`)
- Stay under 256 KB after bundling
- Not use `require()`, `process.`, `eval()`, or other forbidden patterns (the runtime provides Playwright)

Local imports are bundled automatically by `oackctl` via esbuild — `playwright` is marked as external.
