import createMiddleware from "next-intl/middleware";

import { routing } from "@/i18n/routing";

export default createMiddleware(routing);

export const config = {
  matcher: [
    // Match all paths except:
    // - API routes, Next.js internal routes
    // - Static files (any path containing a dot: .ico, .png, .txt, etc.)
    "/((?!api|_next|_vercel|.*\\..*).*)",
  ],
};
