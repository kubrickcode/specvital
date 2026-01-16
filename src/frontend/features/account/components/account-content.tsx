"use client";

import { useEffect } from "react";

import { useAuth } from "@/features/auth";
import { useRouter } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

import { useSubscription, useUsage } from "../hooks";
import { PlanSection } from "./plan-section";
import { UsageSection } from "./usage-section";

export const AccountContent = () => {
  const router = useRouter();
  const { isAuthenticated, isLoading: authLoading } = useAuth();

  const shouldFetch = !authLoading && isAuthenticated;
  const {
    data: subscription,
    error: subscriptionError,
    isLoading: subscriptionLoading,
  } = useSubscription(shouldFetch);
  const { data: usage, error: usageError, isLoading: usageLoading } = useUsage(shouldFetch);

  useEffect(() => {
    if (!authLoading && !isAuthenticated) {
      router.push(ROUTES.HOME);
    }
  }, [authLoading, isAuthenticated, router]);

  if (!authLoading && !isAuthenticated) {
    return null;
  }

  const isLoading = authLoading || subscriptionLoading || usageLoading;

  return (
    <div className="grid gap-6 md:grid-cols-2">
      <PlanSection error={subscriptionError} isLoading={isLoading} plan={subscription?.plan} />
      <UsageSection error={usageError} isLoading={isLoading} usage={usage} />
    </div>
  );
};
