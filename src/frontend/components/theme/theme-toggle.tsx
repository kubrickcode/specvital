"use client";

import { Moon, Sun } from "lucide-react";
import { useTheme } from "next-themes";
import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";

export const ThemeToggle = () => {
  const { resolvedTheme, setTheme } = useTheme();
  const [mounted, setMounted] = useState(false);
  const sunRef = useRef<SVGSVGElement>(null);
  const moonRef = useRef<SVGSVGElement>(null);

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
      <Button disabled size="icon" variant="ghost">
        <div className="relative size-5">
          <Sun className="size-5" />
        </div>
        <span className="sr-only">Toggle theme</span>
      </Button>
    );
  }

  return (
    <Button aria-label="Toggle theme" onClick={toggleTheme} size="icon" variant="ghost">
      <div className="relative size-5">
        <Sun
          className={`absolute size-5 transition-[scale,rotate] duration-300 ${isDark ? "rotate-0 scale-100" : "rotate-90 scale-0"}`}
          ref={sunRef}
        />
        <Moon
          className={`absolute size-5 transition-[scale,rotate] duration-300 ${isDark ? "-rotate-90 scale-0" : "rotate-0 scale-100"}`}
          ref={moonRef}
        />
      </div>
      <span className="sr-only">Toggle theme</span>
    </Button>
  );
};
