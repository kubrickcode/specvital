"use client";

import { useEffect } from "react";

import { useAuth } from "@/features/auth";
import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

/**
 * Redirects authenticated users from home page to dashboard.
 * Only affects home page - other pages are unaffected.
 */
export const AuthenticatedRedirect = () => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && isAuthenticated) {
      router.replace(ROUTES.DASHBOARD);
    }
  }, [isAuthenticated, isLoading, router]);

  return null;
};
