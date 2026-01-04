"use client";

import { debounce } from "es-toolkit";
import { Search, X } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback, useEffect, useMemo, useRef, useState } from "react";

import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

const DEBOUNCE_DELAY = 300;

type SearchInputProps = {
  onChange: (value: string) => void;
  value: string;
};

export const SearchInput = ({ onChange, value }: SearchInputProps) => {
  const t = useTranslations("analyze.filter");
  const inputRef = useRef<HTMLInputElement>(null);
  const [localValue, setLocalValue] = useState(value);

  useEffect(() => {
    setLocalValue(value);
  }, [value]);

  const debouncedOnChange = useMemo(
    () => debounce((newValue: string) => onChange(newValue), DEBOUNCE_DELAY),
    [onChange]
  );

  useEffect(() => {
    return () => {
      debouncedOnChange.cancel();
    };
  }, [debouncedOnChange]);

  const handleChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      const newValue = e.target.value;
      setLocalValue(newValue);
      debouncedOnChange(newValue);
    },
    [debouncedOnChange]
  );

  const handleClear = useCallback(() => {
    setLocalValue("");
    onChange("");
    inputRef.current?.focus();
  }, [onChange]);

  const handleKeyDown = useCallback(
    (e: React.KeyboardEvent<HTMLInputElement>) => {
      if (e.key === "Escape") {
        handleClear();
      }
    },
    [handleClear]
  );

  return (
    <div className="relative flex-1" role="search">
      <Search
        aria-hidden="true"
        className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground"
      />
      <Input
        className={cn("pl-9", localValue && "pr-9")}
        onChange={handleChange}
        onKeyDown={handleKeyDown}
        placeholder={t("searchPlaceholder")}
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
  );
};
