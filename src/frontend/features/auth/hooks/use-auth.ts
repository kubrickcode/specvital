"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { usePathname as useNextPathname } from "next/navigation";
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
  const queryClient = useQueryClient();
  const nextPathname = useNextPathname();
  const pathname = usePathname();
  const router = useRouter();

  // Prevent hydration mismatch: check cookie only on client
  // Re-check on pathname change to detect cookie changes during SPA navigation
  const [hasSession, setHasSession] = useState(false);
  const [isSessionChecked, setIsSessionChecked] = useState(false);
  useEffect(() => {
    setHasSession(checkSessionCookie());
    setIsSessionChecked(true);
  }, [nextPathname]);

  const userQuery = useQuery({
    enabled: hasSession,
    queryFn: fetchCurrentUser,
    queryKey: authKeys.user(),
    retry: false,
    staleTime: 5 * 60 * 1000,
  });

  const loginMutation = useMutation({
    mutationFn: fetchLogin,
    onError: () => toast.error("Failed to start login process"),
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

  const isLoading = !isSessionChecked || userQuery.isFetching;

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
