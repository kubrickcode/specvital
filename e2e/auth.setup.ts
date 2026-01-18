import { test as setup } from "@playwright/test";

import { AUTH_FILE, BACKEND_URL } from "./constants";

setup("authenticate", async ({ request, context }) => {
  // Use dev-login endpoint to authenticate E2E test user
  const response = await request.post(`${BACKEND_URL}/api/auth/dev-login`);

  if (!response.ok()) {
    throw new Error(`Dev login failed: ${response.status()}`);
  }

  // Save the authenticated state (cookies) to file
  await context.storageState({ path: AUTH_FILE });
});
