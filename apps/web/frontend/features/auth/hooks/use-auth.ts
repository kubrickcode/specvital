"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { usePathname as useNextPathname } from "next/navigation";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { usePathname, useRouter } from "@/i18n/navigation";
import { authKeys } from "@/lib/api/error-handler";
import type { UserInfo } from "@/lib/api/types";
import { RETURN_TO_KEY, ROUTES } from "@/lib/routes";

import { fetchCurrentUser, fetchLogin, fetchLogout } from "../api";

type UseAuthReturn = {
  isAuthenticated: boolean;
  isLoading: boolean;
  login: () => void;
  loginPending: boolean;
  logout: () => void;
  logoutPending: boolean;
  user: UserInfo | null;
};

const checkSessionCookie = (): boolean => {
  if (typeof document === "undefined") return false;
  return document.cookie.includes("has_session=1");
};

export const useAuth = (): UseAuthReturn => {
  const t = useTranslations("auth.toast");
  const queryClient = useQueryClient();
  const nextPathname = useNextPathname();
  const pathname = usePathname();
  const router = useRouter();

  // Prevent hydration mismatch: check cookie only on client
  // Re-check on pathname change to detect cookie changes during SPA navigation
  const [hasSession, setHasSession] = useState(false);
  const [isSessionChecked, setIsSessionChecked] = useState(false);
  // Track if initial session verification completed (prevents flicker on expired sessions)
  const [isSessionVerified, setIsSessionVerified] = useState(false);

  useEffect(() => {
    const newHasSession = checkSessionCookie();

    // Clear cache when cookie is missing to prevent stale auth state
    if (!newHasSession) {
      queryClient.setQueryData(authKeys.user(), null);
    }

    setHasSession(newHasSession);
    setIsSessionChecked(true);
  }, [nextPathname, queryClient]);

  // Reset verification state when session state changes (e.g., logout → login)
  useEffect(() => {
    setIsSessionVerified(false);
  }, [hasSession]);

  const userQuery = useQuery({
    enabled: hasSession,
    queryFn: fetchCurrentUser,
    queryKey: authKeys.user(),
    retry: false,
  });

  // Mark session as verified when initial fetch completes (success or error)
  // This runs only once per session, preventing infinite loops on subsequent refetches
  useEffect(() => {
    if (!hasSession) return;

    if (!userQuery.isPending && !userQuery.isFetching) {
      setIsSessionVerified(true);
    }
  }, [hasSession, userQuery.isPending, userQuery.isFetching]);

  // Re-check cookie when query errors (handles 401 cookie deletion by error-handler)
  useEffect(() => {
    if (userQuery.error) {
      const newHasSession = checkSessionCookie();
      if (!newHasSession && hasSession) {
        setHasSession(false);
      }
    }
  }, [userQuery.error, hasSession]);

  const loginMutation = useMutation({
    mutationFn: fetchLogin,
    onError: () => toast.error(t("loginFailed")),
    onSuccess: (data) => {
      if (typeof window !== "undefined") {
        // Store returnTo in cookie (readable by server)
        document.cookie = `${RETURN_TO_KEY}=${encodeURIComponent(pathname)}; path=/; max-age=300; SameSite=Lax`;
      }
      window.location.href = data.authUrl;
    },
  });

  const logoutMutation = useMutation({
    mutationFn: fetchLogout,
    onSuccess: () => {
      document.cookie = "has_session=; path=/; max-age=0";
      queryClient.setQueryData(authKeys.user(), null);
      router.push(ROUTES.HOME);
    },
  });

  // Wait for initial session verification before showing authenticated content
  // - No cookie: immediately not loading (unauthenticated)
  // - Has cookie: wait for first fetch to complete (prevents flicker on expired sessions)
  const isLoading = !isSessionChecked || (hasSession && !isSessionVerified);

  return {
    isAuthenticated: !!userQuery.data,
    isLoading,
    login: () => loginMutation.mutate(),
    loginPending: loginMutation.isPending,
    logout: () => logoutMutation.mutate(),
    logoutPending: logoutMutation.isPending,
    user: userQuery.data ?? null,
  };
};
