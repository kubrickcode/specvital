"use client";

import Image from "next/image";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { AuthErrorBoundary } from "@/components/feedback";
import { NavigationTabs } from "@/components/layout/navigation-tabs";
import { LanguageSelector, ThemeToggle } from "@/components/theme";
import { AuthStatus, LoginModal, LoginRequiredForSpecDialog, SignInTrigger } from "@/features/auth";
import { AnalyzeDialog } from "@/features/home";
import { Link, usePathname } from "@/i18n/navigation";

export const Header = () => {
  const pathname = usePathname();
  const tCommon = useTranslations("common");
  const tHeader = useTranslations("header");
  const isHomePage = pathname === "/";
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      setIsScrolled(window.scrollY > 0);
    };

    window.addEventListener("scroll", handleScroll, { passive: true });
    return () => window.removeEventListener("scroll", handleScroll);
  }, []);

  return (
    <>
      <a
        className="sr-only focus:not-sr-only focus:absolute focus:left-4 focus:top-4 focus:z-[60] focus:rounded-md focus:bg-background focus:px-4 focus:py-2 focus:text-foreground focus:outline-none focus:ring-2 focus:ring-ring"
        href="#main-content"
      >
        {tHeader("skipToContent")}
      </a>
      <header
        className={`sticky top-0 z-50 w-full border-b border-border/30 bg-background/80 backdrop-blur-xl backdrop-saturate-150 transition-all duration-300 supports-[backdrop-filter]:bg-background/50 ${isScrolled ? "shadow-[var(--shadow-card)]" : ""}`}
      >
        <div className="container mx-auto flex h-14 max-w-screen-2xl items-center justify-between px-4">
          <div className="flex items-center gap-6">
            <Link
              aria-label={tHeader("goToHomepage")}
              className="flex items-center gap-2.5 text-lg font-semibold transition-all duration-200 hover:opacity-90"
              href="/"
            >
              <Image alt={tCommon("appName")} height={28} src="/logo.png" width={28} />
              <span>{tCommon("appName")}</span>
            </Link>
            <NavigationTabs />
          </div>

          <div className="hidden items-center gap-2 md:flex">
            {!isHomePage && <AnalyzeDialog variant="header" />}
            <LanguageSelector />
            <ThemeToggle />
            <div
              aria-hidden="true"
              className="mx-1.5 h-6 w-px bg-gradient-to-b from-transparent via-border to-transparent"
            />
            <AuthErrorBoundary fallback={<SignInTrigger />}>
              <AuthStatus />
            </AuthErrorBoundary>
          </div>
        </div>
      </header>
      <LoginModal />
      <LoginRequiredForSpecDialog />
    </>
  );
};
