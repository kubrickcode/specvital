"use client";

import { type ReactNode } from "react";

import { LoadingFallback } from "@/components/feedback";

import { useAuth } from "../hooks/use-auth";

type AuthLoadingWrapperProps = {
  children: ReactNode;
  loadingMessage?: string;
};

/**
 * Displays loading fallback until auth check completes.
 * Prevents flicker when navigating between auth-dependent routes.
 */
export const AuthLoadingWrapper = ({ children, loadingMessage }: AuthLoadingWrapperProps) => {
  const { isLoading } = useAuth();

  if (isLoading) {
    return <LoadingFallback message={loadingMessage} />;
  }

  return <>{children}</>;
};
