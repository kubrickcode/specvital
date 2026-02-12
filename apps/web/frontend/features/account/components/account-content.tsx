"use client";

import { RequireAuth } from "@/features/auth";

import { useSubscription, useUsage } from "../hooks";
import { PlanSection } from "./plan-section";
import { UsageSection } from "./usage-section";

export const AccountContent = () => {
  return (
    <RequireAuth>
      <AccountContentInner />
    </RequireAuth>
  );
};

const AccountContentInner = () => {
  const {
    data: subscription,
    error: subscriptionError,
    isLoading: subscriptionLoading,
  } = useSubscription(true);
  const { data: usage, error: usageError, isLoading: usageLoading } = useUsage(true);

  const isLoading = subscriptionLoading || usageLoading;

  return (
    <div className="grid gap-6 md:grid-cols-2">
      <PlanSection error={subscriptionError} isLoading={isLoading} plan={subscription?.plan} />
      <UsageSection error={usageError} isLoading={isLoading} usage={usage} />
    </div>
  );
};
