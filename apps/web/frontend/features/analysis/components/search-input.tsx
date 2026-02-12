"use client";

import { debounce } from "es-toolkit";
import { Search, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useEffect, useRef, useState } from "react";

import { Badge } from "@/components/ui/badge";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const DEBOUNCE_DELAY = 300;

type SearchInputProps = {
  matchCount?: number;
  onChange: (value: string) => void;
  placeholder?: string;
  value: string;
};

export const SearchInput = ({ matchCount, onChange, placeholder, value }: SearchInputProps) => {
  const t = useTranslations("analyze.filter");
  const inputRef = useRef<HTMLInputElement>(null);
  const [localValue, setLocalValue] = useState(value);

  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  const debouncedOnChange = debounce((newValue: string) => onChange(newValue), DEBOUNCE_DELAY);

  useEffect(() => {
    return () => {
      debouncedOnChange.cancel();
    };
  }, [debouncedOnChange]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setLocalValue(newValue);
    debouncedOnChange(newValue);
  };

  const handleClear = () => {
    setLocalValue("");
    onChange("");
    inputRef.current?.focus();
  };

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Escape") {
      handleClear();
    }
  };

  const hasQuery = localValue.trim().length > 0;
  const showMatchCount = hasQuery && matchCount !== undefined;

  return (
    <div className="flex items-center gap-2 flex-1">
      <div className="relative flex-1" role="search">
        <Search
          aria-hidden="true"
          className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
        />
        <Input
          className={cn("pl-9", localValue && "pr-9")}
          onChange={handleChange}
          onKeyDown={handleKeyDown}
          placeholder={placeholder ?? t("searchPlaceholder")}
          ref={inputRef}
          type="text"
          value={localValue}
        />
        {localValue && (
          <button
            aria-label={t("clearSearch")}
            className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground transition-colors"
            onClick={handleClear}
            type="button"
          >
            <X className="h-4 w-4" />
          </button>
        )}
      </div>
      {showMatchCount && (
        <Badge
          aria-live="polite"
          className="shrink-0"
          variant={matchCount === 0 ? "destructive" : "secondary"}
        >
          {t("matchCount", { count: matchCount })}
        </Badge>
      )}
    </div>
  );
};
