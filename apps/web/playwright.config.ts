import path from "node:path";

import { defineConfig, devices } from "@playwright/test";

import { AUTH_FILE, BACKEND_URL, FRONTEND_URL } from "./e2e/constants";

const WEBSERVER_TIMEOUT_MS = 3 * 60 * 1000;
const isCI = !!process.env.CI;

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
    trace: "on-first-retry",
    screenshot: "only-on-failure",
    video: "retain-on-failure",
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
      command: isCI ? "go run ./cmd/server" : "air",
      cwd: path.join(__dirname, "src/backend"),
      env: isCI ? {} : { DATABASE_URL: process.env.LOCAL_DATABASE_URL ?? "" },
      url: `${BACKEND_URL}/health`,
      reuseExistingServer: !isCI,
      timeout: WEBSERVER_TIMEOUT_MS,
      gracefulShutdown: { signal: "SIGTERM", timeout: 5000 },
    },
    {
      command: "pnpm dev",
      cwd: path.join(__dirname, "src/frontend"),
      env: { PLAYWRIGHT: "1" },
      url: FRONTEND_URL,
      reuseExistingServer: !isCI,
      timeout: WEBSERVER_TIMEOUT_MS,
      gracefulShutdown: { signal: "SIGTERM", timeout: 5000 },
    },
  ],
});
