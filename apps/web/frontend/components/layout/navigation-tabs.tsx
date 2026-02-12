"use client";

import { useTranslations } from "next-intl";

import { useAuth } from "@/features/auth";
import { Link, usePathname } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

type NavItem = {
  href: string;
  labelKey: "dashboard" | "explore" | "pricing" | "docs";
  requiresAuth?: boolean;
};

const NAV_ITEMS: NavItem[] = [
  { href: "/dashboard", labelKey: "dashboard", requiresAuth: true },
  { href: "/explore", labelKey: "explore" },
  { href: "/docs", labelKey: "docs" },
  { href: "/pricing", labelKey: "pricing" },
];

type NavigationTabsProps = {
  className?: string;
};

export const NavigationTabs = ({ className }: NavigationTabsProps) => {
  const pathname = usePathname();
  const t = useTranslations("navigation");
  const { isAuthenticated } = useAuth();

  const visibleItems = NAV_ITEMS.filter((item) => !item.requiresAuth || isAuthenticated);

  const isActive = (href: string) => {
    if (href === "/dashboard") {
      return pathname === "/dashboard" || pathname.startsWith("/dashboard/");
    }
    if (href === "/explore") {
      return pathname === "/explore" || pathname.startsWith("/explore/");
    }
    if (href === "/docs") {
      return pathname === "/docs" || pathname.startsWith("/docs/");
    }
    return pathname === href;
  };

  return (
    <nav aria-label={t("ariaLabel")} className={cn("hidden items-center gap-1 md:flex", className)}>
      {visibleItems.map((item) => {
        const active = isActive(item.href);
        return (
          <Link
            aria-current={active ? "page" : undefined}
            className={cn(
              "relative px-3 py-1.5 text-sm font-medium transition-colors duration-200",
              "rounded-md outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-background",
              active
                ? "text-foreground"
                : "text-muted-foreground hover:text-foreground hover:bg-accent/50"
            )}
            href={item.href}
            key={item.href}
          >
            {t(item.labelKey)}
            {active && (
              <span
                aria-hidden="true"
                className="absolute inset-x-1 -bottom-[13px] h-0.5 rounded-full bg-foreground"
              />
            )}
          </Link>
        );
      })}
    </nav>
  );
};
