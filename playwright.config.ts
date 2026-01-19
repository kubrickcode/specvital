import { defineConfig, devices } from "@playwright/test";

import { AUTH_FILE, BACKEND_URL, FRONTEND_URL } from "./e2e/constants";

const WEBSERVER_TIMEOUT_MS = 3 * 60 * 1000;
const isCI = !!process.env.CI;

const backendCommand = isCI
  ? "cd src/backend && go run ./cmd/server"
  : 'cd src/backend && DATABASE_URL="$LOCAL_DATABASE_URL" air';

export default defineConfig({
  testDir: "./e2e",
  outputDir: "./test-results",
  fullyParallel: true,
  forbidOnly: isCI,
  retries: isCI ? 2 : 0,
  workers: isCI ? 2 : 4,
  reporter: isCI
    ? [["blob"], ["html", { open: "never", outputFolder: "playwright-report" }]]
    : [["html", { open: "never", outputFolder: "playwright-report" }]],
  use: {
    baseURL: FRONTEND_URL,
    trace: isCI ? "on-first-retry" : "on",
    screenshot: isCI ? "only-on-failure" : "on",
    video: isCI ? "retain-on-failure" : "on",
  },
  projects: [
    {
      name: "setup",
      testMatch: /.*\.setup\.ts/,
    },
    {
      name: "chromium",
      dependencies: ["setup"],
      testMatch: /.*\.auth\.spec\.ts/,
      use: {
        ...devices["Desktop Chrome"],
        storageState: AUTH_FILE,
      },
    },
    {
      name: "chromium-mocked",
      dependencies: ["setup"],
      testMatch: /.*\.mocked\.spec\.ts/,
      use: {
        ...devices["Desktop Chrome"],
        storageState: AUTH_FILE,
      },
    },
    {
      name: "chromium-no-auth",
      testIgnore: /.*\.(auth|mocked)\.spec\.ts/,
      use: { ...devices["Desktop Chrome"] },
    },
  ],
  webServer: [
    {
      command: backendCommand,
      url: `${BACKEND_URL}/health`,
      reuseExistingServer: !isCI,
      timeout: WEBSERVER_TIMEOUT_MS,
    },
    {
      command: "cd src/frontend && PLAYWRIGHT=1 pnpm dev",
      url: FRONTEND_URL,
      reuseExistingServer: !isCI,
      timeout: WEBSERVER_TIMEOUT_MS,
    },
  ],
});
