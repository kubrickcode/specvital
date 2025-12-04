import type { NextConfig } from "next";

const isDev = process.env.NODE_ENV === "development";

// Development: Allow unsafe-inline/eval for HMR and dev tools
// Production: Strict CSP without unsafe directives
const cspValue = isDev
  ? "default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self' data:;"
  : "default-src 'self'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self' data:;";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  headers: async () => [
    {
      source: "/:path*",
      headers: [
        { key: "X-Frame-Options", value: "DENY" },
        { key: "X-Content-Type-Options", value: "nosniff" },
        { key: "Referrer-Policy", value: "strict-origin-when-cross-origin" },
        { key: "Content-Security-Policy", value: cspValue },
      ],
    },
  ],
};

export default nextConfig;
