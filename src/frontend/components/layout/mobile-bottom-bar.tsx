"use client";

import {
  Compass,
  Globe,
  LayoutDashboard,
  Loader2,
  LogIn,
  LogOut,
  Moon,
  Plus,
  Sun,
  User,
} from "lucide-react";
import { motion } from "motion/react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useTheme } from "next-themes";
import { useEffect, useState, useTransition } from "react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useAuth, useLoginModal } from "@/features/auth";
import { AnalyzeDialog } from "@/features/home";
import { isValidLocale, LANGUAGE_NAMES } from "@/i18n/config";
import { Link, usePathname, useRouter } from "@/i18n/navigation";
import { routing } from "@/i18n/routing";
import { slideInUp, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

const MobileLanguageSelector = () => {
  const router = useRouter();
  const pathname = usePathname();
  const params = useParams();
  const [isPending, startTransition] = useTransition();
  const t = useTranslations("header");

  const localeParam = params.locale;
  const currentLocale =
    typeof localeParam === "string" && isValidLocale(localeParam)
      ? localeParam
      : routing.defaultLocale;

  const handleLocaleChange = (locale: string) => {
    startTransition(() => {
      router.replace(pathname, { locale });
    });
  };

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button
          aria-label={t("selectLanguage")}
          className="flex-col gap-0.5"
          disabled={isPending}
          size="mobile-nav"
          variant="mobile-nav"
        >
          <Globe className="size-5" />
          <span className="text-[10px] font-normal">{t("language")}</span>
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="center" side="top" sideOffset={12}>
        {routing.locales.map((locale) => (
          <DropdownMenuItem
            disabled={locale === currentLocale}
            key={locale}
            onClick={() => handleLocaleChange(locale)}
          >
            {LANGUAGE_NAMES[locale]}
          </DropdownMenuItem>
        ))}
      </DropdownMenuContent>
    </DropdownMenu>
  );
};

const MobileThemeToggle = () => {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const t = useTranslations("header");

  useEffect(() => {
    setMounted(true);
  }, []);

  const toggleTheme = () => {
    setTheme(resolvedTheme === "dark" ? "light" : "dark");
  };

  const isDark = resolvedTheme === "dark";

  if (!mounted) {
    return (
      <Button className="flex-col gap-0.5" disabled size="mobile-nav" variant="mobile-nav">
        <Sun className="size-5" />
        <span className="text-[10px] font-normal">{t("theme")}</span>
      </Button>
    );
  }

  return (
    <Button
      aria-label={t("toggleTheme")}
      className="flex-col gap-0.5"
      onClick={toggleTheme}
      size="mobile-nav"
      variant="mobile-nav"
    >
      {isDark ? <Sun className="size-5" /> : <Moon className="size-5" />}
      <span className="text-[10px] font-normal">{t("theme")}</span>
    </Button>
  );
};

const MobileAuthAction = () => {
  const t = useTranslations("header");
  const tAuth = useTranslations("auth");
  const { isAuthenticated, isLoading, logout, logoutPending, user } = useAuth();
  const { open: openLoginModal } = useLoginModal();

  if (isLoading) {
    return (
      <Button className="flex-col gap-0.5" disabled size="mobile-nav" variant="mobile-nav">
        <Loader2 className="size-5 animate-spin" />
        <span className="text-[10px] font-normal">{t("account")}</span>
      </Button>
    );
  }

  if (isAuthenticated && user) {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button
            aria-label={t("account")}
            className="flex-col gap-0.5"
            size="mobile-nav"
            variant="mobile-nav"
          >
            {user.avatarUrl ? (
              <img
                alt={user.login}
                className="size-6 rounded-full"
                crossOrigin="anonymous"
                referrerPolicy="no-referrer"
                src={user.avatarUrl}
              />
            ) : (
              <User className="size-5" />
            )}
            <span className="text-[10px] font-normal">{t("account")}</span>
          </Button>
        </DropdownMenuTrigger>
        <DropdownMenuContent align="end" side="top" sideOffset={12}>
          <DropdownMenuItem asChild>
            <Link href="/dashboard">
              <LayoutDashboard className="mr-2 size-4" />
              {t("dashboard")}
            </Link>
          </DropdownMenuItem>
          <DropdownMenuSeparator />
          <DropdownMenuItem disabled={logoutPending} onClick={logout}>
            <LogOut className="mr-2 size-4" />
            {tAuth("logout")}
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    );
  }

  return (
    <Button
      aria-label={tAuth("login")}
      className="flex-col gap-0.5"
      onClick={openLoginModal}
      size="mobile-nav"
      variant="mobile-nav"
    >
      <LogIn className="size-5" />
      <span className="text-[10px] font-normal">{tAuth("login")}</span>
    </Button>
  );
};

type MobileNavItemProps = {
  href: string;
  icon: React.ReactNode;
  isActive: boolean;
  label: string;
};

const MobileNavItem = ({ href, icon, isActive, label }: MobileNavItemProps) => {
  return (
    <Link
      aria-current={isActive ? "page" : undefined}
      className={cn(
        "flex h-full flex-col items-center justify-center gap-0.5 px-3 text-muted-foreground transition-colors",
        isActive && "text-primary"
      )}
      href={href}
    >
      {icon}
      <span className="text-[10px] font-normal">{label}</span>
    </Link>
  );
};

export const MobileBottomBar = () => {
  const t = useTranslations("header");
  const tNav = useTranslations("navigation");
  const pathname = usePathname();
  const shouldReduceMotion = useReducedMotion();
  const { isAuthenticated } = useAuth();
  const isHomePage = pathname === "/";

  const isDashboardActive = pathname === "/dashboard" || pathname.startsWith("/dashboard/");
  const isExploreActive = pathname === "/explore" || pathname.startsWith("/explore/");

  return (
    <motion.nav
      animate="visible"
      aria-label={t("mobileNavigation")}
      className="pb-safe fixed inset-x-0 bottom-0 z-50 border-t border-border/40 bg-background/90 shadow-[0_-2px_8px_0_rgb(0_0_0/0.04)] backdrop-blur-xl backdrop-saturate-150 supports-[backdrop-filter]:bg-background/60 dark:shadow-[0_-2px_8px_0_rgb(0_0_0/0.2)] md:hidden"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={shouldReduceMotion ? undefined : slideInUp}
    >
      <div className="flex h-16 items-center justify-between px-2">
        <div className="flex h-full items-center">
          {isAuthenticated && (
            <MobileNavItem
              href="/dashboard"
              icon={<LayoutDashboard className="size-5" />}
              isActive={isDashboardActive}
              label={tNav("dashboard")}
            />
          )}
          <MobileNavItem
            href="/explore"
            icon={<Compass className="size-5" />}
            isActive={isExploreActive}
            label={tNav("explore")}
          />
        </div>
        <div className="flex h-full items-center">
          {!isHomePage && (
            <AnalyzeDialog variant="header">
              <Button className="flex-col gap-0.5" size="mobile-nav" variant="mobile-nav">
                <Plus className="size-5" />
                <span className="text-[10px] font-normal">{t("analyze")}</span>
              </Button>
            </AnalyzeDialog>
          )}
          <MobileThemeToggle />
          <MobileLanguageSelector />
          <MobileAuthAction />
        </div>
      </div>
    </motion.nav>
  );
};
