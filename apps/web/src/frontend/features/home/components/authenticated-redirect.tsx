"use client";

import type { ReactNode } from "react";
import { useEffect } from "react";

import { LoadingFallback } from "@/components/feedback";
import { useAuth } from "@/features/auth";
import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

type AuthenticatedRedirectProps = {
  children?: ReactNode;
  showLoading?: boolean;
};

/**
 * Redirects authenticated users to dashboard.
 *
 * Usage patterns:
 * - Without children: Acts as redirect-only component
 * - With children: Wraps content, showing loading during auth check and redirect
 *
 * @param showLoading - When true, shows loading fallback during auth check (for non-children mode)
 * @param children - Content to render for unauthenticated users
 */
export const AuthenticatedRedirect = ({
  children,
  showLoading = false,
}: AuthenticatedRedirectProps) => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.replace(ROUTES.DASHBOARD);
    }
  }, [isAuthenticated, isLoading, router]);

  // When children are provided, show loading during auth check and redirect
  if (children !== undefined) {
    if (isLoading || isAuthenticated) {
      return <LoadingFallback />;
    }
    return <>{children}</>;
  }

  // Without children, optionally show loading during auth check only
  if (showLoading && isLoading) {
    return <LoadingFallback />;
  }

  return null;
};
