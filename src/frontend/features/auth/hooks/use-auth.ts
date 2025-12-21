"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { authKeys } from "@/lib/api/error-handler";
import type { UserInfo } from "@/lib/api/types";

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

export const useAuth = (): UseAuthReturn => {
  const queryClient = useQueryClient();

  const userQuery = useQuery({
    queryFn: fetchCurrentUser,
    queryKey: authKeys.user(),
    retry: false,
    staleTime: 5 * 60 * 1000,
  });

  const loginMutation = useMutation({
    mutationFn: fetchLogin,
    onError: () => toast.error("Failed to start login process"),
    onSuccess: (data) => {
      window.location.href = data.authUrl;
    },
  });

  const logoutMutation = useMutation({
    mutationFn: fetchLogout,
    onSuccess: () => queryClient.setQueryData(authKeys.user(), null),
  });

  return {
    isAuthenticated: !!userQuery.data,
    isLoading: userQuery.isPending,
    login: () => loginMutation.mutate(),
    loginPending: loginMutation.isPending,
    logout: () => logoutMutation.mutate(),
    logoutPending: logoutMutation.isPending,
    user: userQuery.data ?? null,
  };
};
