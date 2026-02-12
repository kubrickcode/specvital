"use client";

import { Moon, Sun } from "lucide-react";
import { useTranslations } from "next-intl";
import { useTheme } from "next-themes";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";

export const ThemeToggle = () => {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const sunRef = useRef<SVGSVGElement>(null);
  const moonRef = useRef<SVGSVGElement>(null);
  const t = useTranslations("header");

  useEffect(() => {
    setMounted(true);
  }, []);

  const toggleTheme = () => {
    const newTheme = resolvedTheme === "dark" ? "light" : "dark";

    if (sunRef.current && moonRef.current) {
      if (newTheme === "dark") {
        sunRef.current.classList.remove("rotate-90", "scale-0");
        sunRef.current.classList.add("rotate-0", "scale-100");
        moonRef.current.classList.remove("rotate-0", "scale-100");
        moonRef.current.classList.add("-rotate-90", "scale-0");
      } else {
        sunRef.current.classList.remove("rotate-0", "scale-100");
        sunRef.current.classList.add("rotate-90", "scale-0");
        moonRef.current.classList.remove("-rotate-90", "scale-0");
        moonRef.current.classList.add("rotate-0", "scale-100");
      }
    }

    // Delay setTheme until animation completes (300ms duration)
    setTimeout(() => setTheme(newTheme), 300);
  };

  const isDark = resolvedTheme === "dark";

  if (!mounted) {
    return (
      <Button disabled size="header-icon" variant="header-action">
        <div className="relative size-4">
          <Sun className="size-4" />
        </div>
        <span className="sr-only">{t("toggleTheme")}</span>
      </Button>
    );
  }

  return (
    <ResponsiveTooltip content={t("toggleTheme")} side="bottom" sideOffset={8}>
      <Button
        aria-label={t("toggleTheme")}
        onClick={toggleTheme}
        size="header-icon"
        variant="header-action"
      >
        <div className="relative size-4">
          <Sun
            className={`absolute size-4 transition-[scale,rotate,opacity] duration-300 ease-out ${isDark ? "rotate-0 scale-100 opacity-100" : "rotate-90 scale-0 opacity-0"}`}
            ref={sunRef}
          />
          <Moon
            className={`absolute size-4 transition-[scale,rotate,opacity] duration-300 ease-out ${isDark ? "-rotate-90 scale-0 opacity-0" : "rotate-0 scale-100 opacity-100"}`}
            ref={moonRef}
          />
        </div>
        <span className="sr-only">{t("toggleTheme")}</span>
      </Button>
    </ResponsiveTooltip>
  );
};
