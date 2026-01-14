"use client";

import {
  AlertCircle,
  Building2,
  ChevronDown,
  ExternalLink,
  Link2,
  Lock,
  Loader2,
  RefreshCw,
  Search,
} from "lucide-react";
import { useFormatter, useNow, useTranslations } from "next-intl";
import { useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Input } from "@/components/ui/input";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { Skeleton } from "@/components/ui/skeleton";
import { OrgConnectionBanner } from "@/features/dashboard/components";
import {
  useAllOrgRepos,
  useGitHubAppInstallUrl,
  useOrganizations,
  usePaginatedRepositories,
} from "@/features/dashboard/hooks";
import { Link } from "@/i18n/navigation";
import type { GitHubOrganization, GitHubRepository } from "@/lib/api/types";
import { cn } from "@/lib/utils";

export const OrgReposTab = () => {
  const t = useTranslations("explore");
  const tApp = useTranslations("dashboard.githubApp");
  const format = useFormatter();
  const now = useNow({ updateInterval: 1000 * 60 });

  const [selectedOrg, setSelectedOrg] = useState<GitHubOrganization | null>(null);
  const [searchQuery, setSearchQuery] = useState("");

  const { data: analyzedRepositories, isLoading: isLoadingAnalyzed } = usePaginatedRepositories({});

  const {
    data: organizations,
    isLoading: isLoadingOrgs,
    isRefreshing: isRefreshingOrgs,
  } = useOrganizations();

  const {
    isLoading: isLoadingOrgRepos,
    orgData,
    refreshOrg,
  } = useAllOrgRepos({
    analyzedRepositories: isLoadingAnalyzed ? [] : analyzedRepositories,
    organizations,
  });

  const { install, isLoading: isInstallingApp } = useGitHubAppInstallUrl();

  const currentOrg = selectedOrg ?? (organizations.length === 1 ? organizations[0] : null);
  const currentOrgData = currentOrg ? orgData[currentOrg.login] : null;

  const analyzedFullNames = new Set(analyzedRepositories.map((r) => r.fullName.toLowerCase()));

  const allRepositories = currentOrgData ? currentOrgData.repositories : [];

  const filteredRepositories = (() => {
    if (!searchQuery.trim()) return allRepositories;
    const query = searchQuery.toLowerCase();
    return allRepositories.filter(
      (repo) =>
        repo.name.toLowerCase().includes(query) ||
        repo.fullName.toLowerCase().includes(query) ||
        repo.description?.toLowerCase().includes(query)
    );
  })();

  const isLoading = isLoadingAnalyzed || isLoadingOrgs || isLoadingOrgRepos;
  const isOrgRestricted = currentOrg && currentOrg.accessStatus !== "accessible";

  const handleSelectOrg = (org: GitHubOrganization) => {
    setSelectedOrg(org);
    setSearchQuery("");
  };

  const handleRefreshCurrentOrg = () => {
    if (currentOrg) {
      refreshOrg(currentOrg.login);
    }
  };

  const isRepoAnalyzed = (repo: GitHubRepository): boolean => {
    return analyzedFullNames.has(repo.fullName.toLowerCase());
  };

  const getAccessStatusBadge = (org: GitHubOrganization) => {
    if (org.accessStatus === "accessible") {
      return null;
    }

    if (org.accessStatus === "pending") {
      return (
        <Badge className="bg-yellow-100 text-yellow-800 hover:bg-yellow-100" variant="secondary">
          <AlertCircle className="mr-1 size-3" />
          {tApp("statusPending")}
        </Badge>
      );
    }

    return (
      <Badge className="bg-orange-100 text-orange-800 hover:bg-orange-100" variant="secondary">
        <Link2 className="mr-1 size-3" />
        {tApp("statusRestricted")}
      </Badge>
    );
  };

  if (isLoading) {
    return <OrgReposTabSkeleton />;
  }

  if (organizations.length === 0) {
    return <EmptyOrgsState isLoading={isInstallingApp} onInstall={install} />;
  }

  return (
    <div className="space-y-6">
      <Card>
        <CardHeader className="pb-4">
          <div className="flex items-center justify-between">
            <div>
              <CardTitle className="text-lg">{t("organizations.title")}</CardTitle>
              <CardDescription>{t("organizations.description")}</CardDescription>
            </div>
            <div className="flex items-center gap-2">
              {organizations.length > 1 && (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button className="gap-2" variant="outline">
                      {currentOrg ? (
                        <>
                          {currentOrg.avatarUrl ? (
                            <img
                              alt={currentOrg.login}
                              className="size-5 rounded"
                              src={currentOrg.avatarUrl}
                            />
                          ) : (
                            <Building2 className="size-4" />
                          )}
                          {currentOrg.login}
                        </>
                      ) : (
                        <>
                          <Building2 className="size-4" />
                          {t("organizations.selectOrg")}
                        </>
                      )}
                      <ChevronDown className="size-4 text-muted-foreground" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end" className="w-56">
                    {organizations.map((org) => (
                      <DropdownMenuItem
                        className="flex items-center gap-2"
                        key={org.id}
                        onClick={() => handleSelectOrg(org)}
                      >
                        {org.avatarUrl ? (
                          <img alt={org.login} className="size-5 rounded" src={org.avatarUrl} />
                        ) : (
                          <Building2 className="size-4" />
                        )}
                        <span className="flex-1">{org.login}</span>
                        {getAccessStatusBadge(org)}
                      </DropdownMenuItem>
                    ))}
                  </DropdownMenuContent>
                </DropdownMenu>
              )}
              {currentOrg && !isOrgRestricted && (
                <ResponsiveTooltip content={t("organizations.refresh")}>
                  <Button
                    aria-label={t("organizations.refresh")}
                    disabled={isRefreshingOrgs}
                    onClick={handleRefreshCurrentOrg}
                    size="icon"
                    variant="outline"
                  >
                    <RefreshCw className={cn("size-4", isRefreshingOrgs && "animate-spin")} />
                  </Button>
                </ResponsiveTooltip>
              )}
            </div>
          </div>
        </CardHeader>

        <CardContent className="space-y-4">
          {currentOrg && isOrgRestricted && (
            <OrgConnectionBanner
              accessStatus={currentOrg.accessStatus}
              orgLogin={currentOrg.login}
            />
          )}

          {currentOrg && !isOrgRestricted && (
            <>
              <div className="relative">
                <Search className="absolute left-3 top-1/2 -translate-y-1/2 size-4 text-muted-foreground" />
                <Input
                  className="pl-9"
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder={t("searchPlaceholder")}
                  type="search"
                  value={searchQuery}
                />
              </div>

              <ul className="space-y-2">
                {filteredRepositories.length === 0 ? (
                  <li className="text-center py-8 text-sm text-muted-foreground">
                    {searchQuery
                      ? t("myRepos.noSearchResults")
                      : t("organizations.noOrgsDescription")}
                  </li>
                ) : (
                  filteredRepositories.map((repo) => {
                    const isAnalyzed = isRepoAnalyzed(repo);
                    return (
                      <li key={repo.id}>
                        <RepoListItem
                          format={format}
                          isAnalyzed={isAnalyzed}
                          now={now}
                          repo={repo}
                          t={t}
                        />
                      </li>
                    );
                  })
                )}
              </ul>
            </>
          )}

          {!currentOrg && organizations.length > 1 && (
            <div className="text-center py-8 text-sm text-muted-foreground">
              {t("organizations.selectOrg")}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
};

type RepoListItemProps = {
  format: ReturnType<typeof useFormatter>;
  isAnalyzed: boolean;
  now: Date;
  repo: GitHubRepository;
  t: ReturnType<typeof useTranslations<"explore">>;
};

const RepoListItem = ({ format, isAnalyzed, now, repo, t }: RepoListItemProps) => {
  return (
    <div className="flex items-start gap-3 p-3 rounded-lg border bg-card shadow-sm hover:shadow-md hover:bg-accent/50 transition-all">
      <div className="flex-1 min-w-0">
        <div className="flex items-start gap-2">
          <ResponsiveTooltip content={repo.fullName}>
            <h4 className="font-medium text-sm line-clamp-2 break-all h-10">{repo.fullName}</h4>
          </ResponsiveTooltip>
          {repo.isPrivate && <Lock className="size-3 text-muted-foreground shrink-0 mt-0.5" />}
          {isAnalyzed && (
            <Badge className="text-xs" variant="secondary">
              {t("myRepos.analyzed")}
            </Badge>
          )}
        </div>
        {repo.description && (
          <p className="text-xs text-muted-foreground line-clamp-2 mt-1">{repo.description}</p>
        )}
        {repo.pushedAt && (
          <p className="text-xs text-muted-foreground mt-1">
            {t("myRepos.lastPush")}: {format.relativeTime(new Date(repo.pushedAt), now)}
          </p>
        )}
      </div>
      <div className="flex flex-col gap-1 shrink-0">
        {isAnalyzed ? (
          <Button asChild size="sm" variant="outline">
            <Link href={`/analyze/${repo.owner}/${repo.name}`}>{t("myRepos.view")}</Link>
          </Button>
        ) : (
          <Button asChild size="sm" variant="cta">
            <Link href={`/analyze/${repo.owner}/${repo.name}`}>{t("myRepos.analyze")}</Link>
          </Button>
        )}
        <Button asChild className="h-7" size="sm" variant="ghost">
          <a href={`https://github.com/${repo.fullName}`} rel="noopener noreferrer" target="_blank">
            <ExternalLink className="size-3 mr-1" />
            GitHub
          </a>
        </Button>
      </div>
    </div>
  );
};

type EmptyOrgsStateProps = {
  isLoading: boolean;
  onInstall: () => void;
};

const EmptyOrgsState = ({ isLoading, onInstall }: EmptyOrgsStateProps) => {
  const t = useTranslations("explore.organizations");
  const tApp = useTranslations("dashboard.githubApp");

  return (
    <Card>
      <CardContent className="flex flex-col items-center justify-center py-12 text-center">
        <Building2 className="size-12 text-muted-foreground mb-4" />
        <h3 className="font-semibold text-lg mb-2">{t("noOrgs")}</h3>
        <p className="text-sm text-muted-foreground mb-6 max-w-md">{t("noOrgsDescription")}</p>
        <Button disabled={isLoading} onClick={onInstall}>
          {isLoading ? (
            <Loader2 className="mr-2 size-4 animate-spin" />
          ) : (
            <ExternalLink className="mr-2 size-4" />
          )}
          {tApp("installButton")}
        </Button>
      </CardContent>
    </Card>
  );
};

const OrgReposTabSkeleton = () => {
  return (
    <Card>
      <CardHeader className="pb-4">
        <div className="flex items-center justify-between">
          <div className="space-y-2">
            <Skeleton className="h-5 w-40" />
            <Skeleton className="h-4 w-60" />
          </div>
          <Skeleton className="h-9 w-32" />
        </div>
      </CardHeader>
      <CardContent className="space-y-4">
        <Skeleton className="h-10 w-full" />
        <div className="space-y-2">
          {Array.from({ length: 3 }).map((_, i) => (
            <Skeleton className="h-24 w-full" key={i} />
          ))}
        </div>
      </CardContent>
    </Card>
  );
};
