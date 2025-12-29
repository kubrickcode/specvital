"use client";

import Image from "next/image";
import { useTranslations } from "next-intl";

import { AuthErrorBoundary } from "@/components/feedback";
import { LanguageSelector, ThemeToggle } from "@/components/theme";
import { AuthStatus, LoginButton } from "@/features/auth";
import { AnalyzeDialog } from "@/features/home";
import { Link, usePathname } from "@/i18n/navigation";

export const Header = () => {
  const pathname = usePathname();
  const tCommon = useTranslations("common");
  const tHeader = useTranslations("header");
  const isHomePage = pathname === "/";

  return (
    <>
      <a
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-[60] focus:rounded-md focus:bg-background focus:px-4 focus:py-2 focus:text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        href="#main-content"
      >
        {tHeader("skipToContent")}
      </a>
      <header className="sticky top-0 z-50 w-full bg-background/80 backdrop-blur-xl backdrop-saturate-150 supports-[backdrop-filter]:bg-background/50">
        <div className="absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-foreground/20 to-transparent" />
        <div className="absolute inset-x-0 bottom-0 h-px bg-gradient-to-r from-transparent via-foreground/20 to-transparent" />
        <div className="container mx-auto flex h-14 max-w-screen-2xl items-center justify-between px-4">
          <Link
            aria-label={tHeader("goToHomepage")}
            className="flex items-center gap-2.5 text-lg font-semibold transition-all duration-200 hover:opacity-90"
            href="/"
          >
            <Image alt={tCommon("appName")} height={28} src="/logo.png" width={28} />
            <span>{tCommon("appName")}</span>
          </Link>

          <div className="flex items-center gap-1">
            {!isHomePage && <AnalyzeDialog variant="header" />}
            <LanguageSelector />
            <ThemeToggle />
            <div aria-hidden="true" className="mx-1 h-5 w-px bg-border" />
            <AuthErrorBoundary fallback={<LoginButton />}>
              <AuthStatus />
            </AuthErrorBoundary>
          </div>
        </div>
      </header>
    </>
  );
};
