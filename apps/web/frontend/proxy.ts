import type { NextRequest } from "next/server";
import { NextResponse } from "next/server";
import createIntlMiddleware from "next-intl/middleware";

import { defaultLocale, type Locale, locales } from "@/i18n/config";
import { routing } from "@/i18n/routing";

// Maintenance Mode
const MAINTENANCE_MODE = process.env.NEXT_PUBLIC_MAINTENANCE_MODE === "true";
const MAINTENANCE_ALLOWED_PATHS = ["/maintenance", "/_next", "/favicon", "/api/health"];

const AUTH_COOKIE_NAME = "auth_token";
const REFRESH_COOKIE_NAME = "refresh_token";
const SESSION_INDICATOR_COOKIE = "has_session";
const JWT_MIN_LENGTH = 20;
const BACKEND_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";
const ACCESS_TOKEN_MAX_AGE = 15 * 60;
const REFRESH_TOKEN_MAX_AGE = 7 * 24 * 60 * 60;

type AuthResult = {
  isValid: boolean;
  newTokens?: { accessToken: string; refreshToken: string };
};

const intlMiddleware = createIntlMiddleware(routing);

/**
 * Routes that require authentication.
 * Unauthenticated users will be redirected to home.
 */
const PROTECTED_ROUTES = ["dashboard"] as const;

const isHomePage = (pathname: string): boolean => {
  if (pathname === "/") return true;
  return locales.some((locale) => pathname === `/${locale}` || pathname === `/${locale}/`);
};

const isProtectedRoute = (pathname: string): boolean => {
  return locales.some((locale) =>
    PROTECTED_ROUTES.some(
      (route) => pathname.startsWith(`/${locale}/${route}`) || pathname === `/${locale}/${route}`
    )
  );
};

const getLocaleFromRequest = (request: NextRequest, pathname: string): string => {
  for (const locale of locales) {
    if (pathname.startsWith(`/${locale}/`) || pathname === `/${locale}`) {
      return locale;
    }
  }

  const localeCookie = request.cookies.get("NEXT_LOCALE");
  if (localeCookie?.value && locales.includes(localeCookie.value as Locale)) {
    return localeCookie.value;
  }

  return defaultLocale;
};

const hasValidAuthCookie = (request: NextRequest): boolean => {
  const cookie = request.cookies.get(AUTH_COOKIE_NAME);

  if (!cookie?.value) {
    return false;
  }

  if (cookie.value.length < JWT_MIN_LENGTH) {
    return false;
  }

  const parts = cookie.value.split(".");
  if (parts.length !== 3) {
    return false;
  }

  return true;
};

const refreshTokens = async (
  refreshToken: string
): Promise<{ accessToken: string; refreshToken: string } | null> => {
  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000);

    const response = await fetch(`${BACKEND_URL}/api/auth/refresh`, {
      headers: {
        Accept: "application/json",
        Cookie: `${REFRESH_COOKIE_NAME}=${refreshToken}`,
      },
      method: "POST",
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (!response.ok) {
      return null;
    }

    const setCookieHeaders = response.headers.getSetCookie();
    let accessToken: string | null = null;
    let newRefreshToken: string | null = null;

    for (const cookieHeader of setCookieHeaders) {
      const authMatch = cookieHeader.match(/^auth_token=([^;]+)/);
      if (authMatch) {
        accessToken = authMatch[1];
      }
      const refreshMatch = cookieHeader.match(/^refresh_token=([^;]+)/);
      if (refreshMatch) {
        newRefreshToken = refreshMatch[1];
      }
    }

    if (accessToken && newRefreshToken) {
      return { accessToken, refreshToken: newRefreshToken };
    }

    return null;
  } catch {
    return null;
  }
};

const verifyAuthToken = async (request: NextRequest): Promise<AuthResult> => {
  const authCookie = request.cookies.get(AUTH_COOKIE_NAME);
  const refreshCookie = request.cookies.get(REFRESH_COOKIE_NAME);

  if (!authCookie?.value) {
    return { isValid: false };
  }

  const cookieHeader = [
    `${AUTH_COOKIE_NAME}=${authCookie.value}`,
    refreshCookie?.value ? `${REFRESH_COOKIE_NAME}=${refreshCookie.value}` : "",
  ]
    .filter(Boolean)
    .join("; ");

  try {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000);

    const response = await fetch(`${BACKEND_URL}/api/auth/me`, {
      headers: {
        Accept: "application/json",
        Cookie: cookieHeader,
      },
      signal: controller.signal,
    });

    clearTimeout(timeoutId);

    if (response.ok) {
      return { isValid: true };
    }

    if (response.status === 401 && refreshCookie?.value) {
      const newTokens = await refreshTokens(refreshCookie.value);
      if (newTokens) {
        return { isValid: true, newTokens };
      }
    }

    return { isValid: false };
  } catch {
    // Allow optimistically if cookie is valid on network error (re-validated client-side)
    if (hasValidAuthCookie(request)) {
      return { isValid: true };
    }
    return { isValid: false };
  }
};

const createRedirectWithCookieClear = (url: URL): NextResponse => {
  const response = NextResponse.redirect(url, 307);
  response.cookies.delete(AUTH_COOKIE_NAME);
  response.cookies.delete(REFRESH_COOKIE_NAME);
  response.cookies.delete(SESSION_INDICATOR_COOKIE);
  return response;
};

const setTokenCookies = (
  response: NextResponse,
  tokens: { accessToken: string; refreshToken: string }
): void => {
  const isProduction = process.env.NODE_ENV === "production";

  response.cookies.set(AUTH_COOKIE_NAME, tokens.accessToken, {
    httpOnly: true,
    maxAge: ACCESS_TOKEN_MAX_AGE,
    path: "/",
    sameSite: "lax",
    secure: isProduction,
  });

  response.cookies.set(REFRESH_COOKIE_NAME, tokens.refreshToken, {
    httpOnly: true,
    maxAge: REFRESH_TOKEN_MAX_AGE,
    path: "/",
    sameSite: "strict",
    secure: isProduction,
  });

  setSessionIndicator(response);
};

const setSessionIndicator = (response: NextResponse): void => {
  const isProduction = process.env.NODE_ENV === "production";

  response.cookies.set(SESSION_INDICATOR_COOKIE, "1", {
    httpOnly: false,
    maxAge: REFRESH_TOKEN_MAX_AGE,
    path: "/",
    sameSite: "lax",
    secure: isProduction,
  });
};

const proxy = async (request: NextRequest) => {
  const { pathname } = request.nextUrl;
  const locale = getLocaleFromRequest(request, pathname);

  // Maintenance mode: redirect all routes except allowed paths
  if (MAINTENANCE_MODE) {
    const isAllowedPath = MAINTENANCE_ALLOWED_PATHS.some((path) => pathname.startsWith(path));
    const isMaintenancePath = pathname.includes("/maintenance");
    const isStaticFile = /\.(ico|png|jpg|jpeg|svg|gif|webp|css|js|woff2?)$/i.test(pathname);

    if (!isAllowedPath && !isMaintenancePath && !isStaticFile) {
      const maintenanceUrl = new URL(`/${locale}/maintenance`, request.url);
      return NextResponse.redirect(maintenanceUrl, 307);
    }
  }

  const hasCookie = hasValidAuthCookie(request);

  // Protected routes: require authentication
  if (isProtectedRoute(pathname)) {
    if (!hasCookie) {
      const homeUrl = new URL(`/${locale}`, request.url);
      return NextResponse.redirect(homeUrl, 307);
    }

    const authResult = await verifyAuthToken(request);
    if (!authResult.isValid) {
      const homeUrl = new URL(`/${locale}`, request.url);
      return createRedirectWithCookieClear(homeUrl);
    }

    const response = authResult.newTokens ? NextResponse.next() : intlMiddleware(request);
    if (authResult.newTokens) {
      setTokenCookies(response, authResult.newTokens);
    } else {
      setSessionIndicator(response);
    }
    return response;
  }

  // Optimistic redirect: authenticated users on home → redirect to dashboard
  // No backend API call here - token validation happens in dashboard's protected route logic
  if (isHomePage(pathname) && hasCookie) {
    const dashboardUrl = new URL(`/${locale}/dashboard`, request.url);
    return NextResponse.redirect(dashboardUrl, 307);
  }

  return intlMiddleware(request);
};

export default proxy;

export const config = {
  // Exclude: api, _next, _vercel, root-level static files (e.g., favicon.ico)
  // Allow: deep paths with dots (e.g., /analyze/vercel/next.js)
  matcher: ["/((?!api|_next|_vercel|[^/]+\\.[\\w]+$).*)"],
};
