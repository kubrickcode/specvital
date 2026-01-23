"use client";

import { AlertCircle, Loader2, LogIn } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

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

import { useDebouncedSearch, useGlobalSearchStore, useStaticActions } from "../hooks";
import { RepositorySearchItem } from "./repository-search-item";
import { SearchKeyboardHints } from "./search-keyboard-hints";

export const GlobalSearchDialog = () => {
  const { close, isOpen } = useGlobalSearchStore();
  const t = useTranslations("globalSearch");
  const { analyzeDialogOpen, commandItems, navigationItems, setAnalyzeDialogOpen } =
    useStaticActions();
  const { login } = useAuth();

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
  const showLoginPrompt = hasQuery && !isAuthenticated;
  const showNoResults = hasQuery && !hasResults && !isLoading;

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
        <CommandInput
          onValueChange={handleValueChange}
          placeholder={t("placeholder")}
          value={inputValue}
        />
        <CommandList>
          {showNoResults && (
            <CommandEmpty>
              {inputValue.trim()
                ? t("noResults", { query: inputValue.trim() })
                : t("noResultsDefault")}
            </CommandEmpty>
          )}

          {isLoading && hasQuery && (
            <div className="py-6 text-center text-sm text-muted-foreground">
              <Loader2 className="size-4 animate-spin mx-auto mb-2" />
              {t("loading")}
            </div>
          )}

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
                      onSelect={() => navigateToRepository(result.item.owner, result.item.name)}
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
                      onSelect={() => navigateToRepository(result.item.owner, result.item.name)}
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
                      onSelect={() => navigateToRepository(result.item.owner, result.item.name)}
                      result={result}
                    />
                  ))}
                </CommandGroup>
              )}
            </>
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
