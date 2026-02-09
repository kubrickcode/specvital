"use client";

import type { LucideIcon } from "lucide-react";
import { AlertCircle, Bookmark, GitFork, Search } from "lucide-react";
import { useTranslations } from "next-intl";

import { Card } from "@/components/ui/card";

type EmptyStateVariant = "no-repos" | "no-search-results" | "no-bookmarks" | "error";

type EmptyStateVariantProps = {
  action?: React.ReactNode;
  searchQuery?: string;
  variant: EmptyStateVariant;
};

const VARIANT_CONFIG: Record<
  EmptyStateVariant,
  {
    Icon: LucideIcon;
    translationKey: string;
  }
> = {
  error: {
    Icon: AlertCircle,
    translationKey: "error",
  },
  "no-bookmarks": {
    Icon: Bookmark,
    translationKey: "noBookmarks",
  },
  "no-repos": {
    Icon: GitFork,
    translationKey: "noRepos",
  },
  "no-search-results": {
    Icon: Search,
    translationKey: "noSearchResults",
  },
};

export const EmptyStateVariant = ({ action, searchQuery, variant }: EmptyStateVariantProps) => {
  const t = useTranslations("dashboard.emptyStateVariant");
  const config = VARIANT_CONFIG[variant];
  const Icon = config.Icon;

  return (
    <Card className="flex flex-col items-center justify-center py-16 px-6 text-center">
      <div className="rounded-full bg-muted p-4 mb-6">
        <Icon aria-hidden="true" className="size-8 text-muted-foreground" />
      </div>

      <h2 className="text-xl font-semibold mb-2">{t(`${config.translationKey}.title`)}</h2>
      <p className="text-muted-foreground mb-6 max-w-md">
        {variant === "no-search-results" && searchQuery
          ? t(`${config.translationKey}.descriptionWithQuery`, {
              query: searchQuery,
            })
          : t(`${config.translationKey}.description`)}
      </p>

      {action && <div>{action}</div>}
    </Card>
  );
};
