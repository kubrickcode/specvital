/**
 * Application route constants
 * Centralized route definitions for type-safe navigation
 */
export const ROUTES = {
  ACCOUNT: "/account",
  analyze: (owner: string, repo: string) => `/analyze/${owner}/${repo}`,
  DASHBOARD: "/dashboard",
  EXPLORE: "/explore",
  HOME: "/",
  PRICING: "/pricing",
} as const;

/**
 * Session storage key for storing return URL after OAuth login
 */
export const RETURN_TO_KEY = "auth_return_to";
