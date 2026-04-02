import { defineConfig, devices } from '@playwright/test';

const baseURL = process.env.BASE_URL || 'https://poke-store.oack.io';
const isLocal = baseURL.includes('localhost') || baseURL.includes('127.0.0.1');

// Playwright browsers ignore HTTP_PROXY env vars — pass proxy explicitly.
const proxyServer = process.env.HTTPS_PROXY || process.env.HTTP_PROXY;

export default defineConfig({
  testDir: './tests/e2e',
  fullyParallel: true,
  forbidOnly: !!process.env.CI,
  retries: process.env.CI ? 2 : 0,
  workers: process.env.CI ? 1 : undefined,
  reporter: 'html',
  use: {
    baseURL,
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
    actionTimeout: isLocal ? 5_000 : 15_000,
    navigationTimeout: isLocal ? 10_000 : 30_000,
    ...(proxyServer && !isLocal ? { proxy: { server: proxyServer } } : {}),
  },
  expect: {
    timeout: isLocal ? 5_000 : 15_000,
  },
  projects: [
    {
      name: 'chromium',
      use: { ...devices['Desktop Chrome'] },
    },
  ],
  ...(isLocal
    ? {
        webServer: {
          command: 'cd .. && make build && ./bin/poke-store',
          url: 'http://localhost:6001/api/version',
          reuseExistingServer: !process.env.CI,
          timeout: 30_000,
        },
      }
    : {}),
});
