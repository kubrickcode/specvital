"use client";

import { Loader2 } from "lucide-react";

import { useAuth } from "../hooks";
import { LoginButton } from "./login-button";
import { UserMenu } from "./user-menu";

export const AuthStatus = () => {
  const { isAuthenticated, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex h-8 w-8 items-center justify-center">
        <Loader2 className="h-4 w-4 animate-spin text-muted-foreground" />
      </div>
    );
  }

  return isAuthenticated ? <UserMenu /> : <LoginButton />;
};
