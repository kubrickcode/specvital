"use client";

import { AlertCircle, LogIn, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  CommandDialog,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
  CommandShortcut,
} from "@/components/ui/command";
import { useAuth } from "@/features/auth";
import { AnalyzeDialog } from "@/features/home";

import {
  useDebouncedSearch,
  useGlobalSearchStore,
  useRecentItems,
  useStaticActions,
} from "../hooks";
import { RecentSearchItem } from "./recent-search-item";
import { RepositorySearchItem } from "./repository-search-item";
import { SearchKeyboardHints } from "./search-keyboard-hints";
import { SearchSkeleton } from "./search-skeleton";

export const GlobalSearchDialog = () => {
  const { close, isOpen } = useGlobalSearchStore();
  const t = useTranslations("globalSearch");
  const { analyzeDialogOpen, commandItems, navigationItems, setAnalyzeDialogOpen } =
    useStaticActions();
  const { login } = useAuth();
  const { addItem, recentItems } = useRecentItems();

  const {
    groupedResults,
    hasError,
    hasResults,
    isAuthenticated,
    isLoading,
    navigateToRepository,
    setQuery,
  } = useDebouncedSearch();

  const [inputValue, setInputValue] = useState("");

  const handleNavigateToRepository = (
    owner: string,
    repo: string,
    repoId: string,
    fullName: string
  ) => {
    addItem({
      fullName,
      id: repoId,
      name: repo,
      owner,
    });
    navigateToRepository(owner, repo);
  };

  const handleValueChange = (value: string) => {
    setInputValue(value);
    setQuery(value);
  };

  const handleOpenChange = (open: boolean) => {
    if (!open) {
      close();
      setInputValue("");
      setQuery("");
    }
  };

  const hasQuery = inputValue.trim().length > 0;
  const showRepositoryResults = hasQuery && hasResults;
  const showStaticActions = !hasQuery;
  const showRecentItems = !hasQuery && recentItems.length > 0;
  const showLoginPrompt = hasQuery && !isAuthenticated;
  const showNoResults = hasQuery && !hasResults && !isLoading;

  // Calculate total results for screen reader announcement
  const totalResults = hasQuery
    ? groupedResults.repositories.length +
      groupedResults.bookmarks.length +
      groupedResults.community.length
    : 0;

  return (
    <>
      <CommandDialog
        description={t("description")}
        onOpenChange={handleOpenChange}
        open={isOpen}
        shouldFilter={!hasQuery}
        showCloseButton={false}
        title={t("title")}
      >
        <div className="relative">
          <CommandInput
            className="pr-12 md:pr-3"
            onValueChange={handleValueChange}
            placeholder={t("placeholder")}
            value={inputValue}
          />
          {/* Mobile close button */}
          <Button
            aria-label={t("hints.close")}
            className="absolute right-2 top-1/2 size-8 -translate-y-1/2 md:hidden"
            onClick={close}
            size="icon"
            variant="ghost"
          >
            <X className="size-4" />
          </Button>
        </div>
        {/* Screen reader announcement for result count */}
        <div aria-atomic="true" aria-live="polite" className="sr-only" role="status">
          {hasQuery && !isLoading && t("aria.resultsFound", { count: totalResults })}
        </div>
        <CommandList>
          {showNoResults && (
            <CommandEmpty>
              {inputValue.trim()
                ? t("noResults", { query: inputValue.trim() })
                : t("noResultsDefault")}
            </CommandEmpty>
          )}

          {isLoading && hasQuery && <SearchSkeleton />}

          {hasError && hasQuery && !isLoading && (
            <div className="py-6 text-center text-sm text-destructive">
              <AlertCircle className="size-4 mx-auto mb-2" />
              {t("error")}
            </div>
          )}

          {showLoginPrompt && (
            <CommandGroup>
              <CommandItem
                className="flex items-center gap-2 text-muted-foreground"
                onSelect={() => {
                  close();
                  login();
                }}
              >
                <LogIn className="size-4" />
                <span>{t("loginPrompt")}</span>
              </CommandItem>
            </CommandGroup>
          )}

          {showRepositoryResults && (
            <>
              {groupedResults.repositories.length > 0 && (
                <CommandGroup heading={t("categories.repositories")}>
                  {groupedResults.repositories.map((result) => (
                    <RepositorySearchItem
                      key={result.item.id}
                      onSelect={() =>
                        handleNavigateToRepository(
                          result.item.owner,
                          result.item.name,
                          result.item.id,
                          result.item.fullName
                        )
                      }
                      result={result}
                    />
                  ))}
                </CommandGroup>
              )}

              {groupedResults.bookmarks.length > 0 && (
                <CommandGroup heading={t("categories.bookmarks")}>
                  {groupedResults.bookmarks.map((result) => (
                    <RepositorySearchItem
                      key={result.item.id}
                      onSelect={() =>
                        handleNavigateToRepository(
                          result.item.owner,
                          result.item.name,
                          result.item.id,
                          result.item.fullName
                        )
                      }
                      result={result}
                    />
                  ))}
                </CommandGroup>
              )}

              {groupedResults.community.length > 0 && (
                <CommandGroup heading={t("categories.community")}>
                  {groupedResults.community.map((result) => (
                    <RepositorySearchItem
                      key={result.item.id}
                      onSelect={() =>
                        handleNavigateToRepository(
                          result.item.owner,
                          result.item.name,
                          result.item.id,
                          result.item.fullName
                        )
                      }
                      result={result}
                    />
                  ))}
                </CommandGroup>
              )}
            </>
          )}

          {showRecentItems && (
            <CommandGroup heading={t("categories.recent")}>
              {recentItems.map((item) => (
                <RecentSearchItem
                  item={item}
                  key={item.id}
                  onSelect={() =>
                    handleNavigateToRepository(item.owner, item.name, item.id, item.fullName)
                  }
                />
              ))}
            </CommandGroup>
          )}

          {showStaticActions && (
            <>
              {navigationItems.length > 0 && (
                <CommandGroup heading={t("categories.navigation")}>
                  {navigationItems.map((item) => (
                    <CommandItem key={item.id} onSelect={item.onSelect}>
                      {item.icon && <item.icon className="mr-2 size-4" />}
                      <span>{item.label}</span>
                      {item.shortcut && <CommandShortcut>{item.shortcut}</CommandShortcut>}
                    </CommandItem>
                  ))}
                </CommandGroup>
              )}

              {commandItems.length > 0 && (
                <CommandGroup heading={t("categories.commands")}>
                  {commandItems.map((item) => (
                    <CommandItem key={item.id} onSelect={item.onSelect}>
                      {item.icon && <item.icon className="mr-2 size-4" />}
                      <span>{item.label}</span>
                      {item.shortcut && <CommandShortcut>{item.shortcut}</CommandShortcut>}
                    </CommandItem>
                  ))}
                </CommandGroup>
              )}
            </>
          )}
        </CommandList>
        <SearchKeyboardHints />
      </CommandDialog>

      <AnalyzeDialog
        onOpenChange={setAnalyzeDialogOpen}
        open={analyzeDialogOpen}
        variant="header"
      />
    </>
  );
};
