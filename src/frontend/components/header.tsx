"use client";

import { Home } from "lucide-react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { LanguageSelector } from "@/components/language-selector";
import { ThemeToggle } from "@/components/theme-toggle";
import { Link, usePathname } from "@/i18n/navigation";

export const Header = () => {
  const pathname = usePathname();
  const t = useTranslations("header");
  const tCommon = useTranslations("common");
  const isHomePage = pathname === "/";

  return (
    <header className="sticky top-0 z-50 w-full border-b border-border/40 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="container mx-auto flex h-14 max-w-screen-2xl items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link
            href="/"
            className="flex items-center space-x-2 text-lg font-semibold hover:opacity-80 transition-opacity"
          >
            <span>{tCommon("appName")}</span>
          </Link>
        </div>

        <div className="flex items-center gap-2">
          {!isHomePage && (
            <Button variant="ghost" size="sm" asChild>
              <Link href="/" className="flex items-center gap-2">
                <Home className="h-4 w-4" />
                <span className="hidden sm:inline">{t("analyzeAnother")}</span>
              </Link>
            </Button>
          )}
          <LanguageSelector />
          <ThemeToggle />
        </div>
      </div>
    </header>
  );
};
