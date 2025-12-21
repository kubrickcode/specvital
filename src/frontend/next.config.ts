import type { NextConfig } from "next";
import createNextIntlPlugin from "next-intl/plugin";

const withNextIntl = createNextIntlPlugin("./i18n/request.ts");

const isDev = process.env.NODE_ENV === "development";

// Development: Allow unsafe-eval for HMR and dev tools
// Production: Allow unsafe-inline for Next.js hydration scripts
const cspValue = isDev
  ? "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://avatars.githubusercontent.com; font-src 'self' data:;"
  : "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://avatars.githubusercontent.com; font-src 'self' data:;";

const API_URL = process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:8000";

const nextConfig: NextConfig = {
  headers: async () => [
    {
      headers: [
        { key: "X-Frame-Options", value: "DENY" },
        { key: "X-Content-Type-Options", value: "nosniff" },
        { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
        { key: "Content-Security-Policy", value: cspValue },
      ],
      source: "/:path*",
    },
  ],
  reactStrictMode: true,
  // Proxy API requests to backend (solves CORS and localhost access issues)
  rewrites: async () => [
    {
      destination: `${API_URL}/api/:path*`,
      source: "/api/:path*",
    },
  ],
};

export default withNextIntl(nextConfig);
