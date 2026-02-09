"use client";

import { useEffect, type ReactNode } from "react";

import { LoadingFallback } from "@/components/feedback";
import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

import { useAuth } from "../hooks/use-auth";

type RequireAuthProps = {
  children: ReactNode;
  loadingMessage?: string;
};

/**
 * Protects routes by redirecting unauthenticated users to home.
 * Displays loading fallback until auth check completes.
 */
export const RequireAuth = ({ children, loadingMessage }: RequireAuthProps) => {
  const { isAuthenticated, isLoading } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!isLoading && !isAuthenticated) {
      router.replace(ROUTES.HOME);
    }
  }, [isAuthenticated, isLoading, router]);

  if (isLoading) {
    return <LoadingFallback message={loadingMessage} />;
  }

  if (!isAuthenticated) {
    return <LoadingFallback />;
  }

  return <>{children}</>;
};
