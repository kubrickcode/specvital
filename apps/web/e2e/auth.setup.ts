import { test as setup } from "@playwright/test";

import { AUTH_FILE, BACKEND_URL, FRONTEND_URL } from "./constants";

setup("authenticate", async ({ page }) => {
  // Navigate to frontend first to establish browser context
  await page.goto(FRONTEND_URL);

  // Call dev-login endpoint and capture cookies via page context
  const response = await page.request.post(`${BACKEND_URL}/api/auth/dev-login`, {
    data: {},
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!response.ok()) {
    throw new Error(`Dev login failed: ${response.status()}`);
  }

  // Reload page to apply cookies
  await page.reload();

  // Save the authenticated state (cookies) to file
  await page.context().storageState({ path: AUTH_FILE });
});
