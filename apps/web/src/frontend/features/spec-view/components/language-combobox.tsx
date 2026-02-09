"use client";

import { Check, ChevronsUpDown } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Command,
  CommandEmpty,
  CommandGroup,
  CommandInput,
  CommandItem,
  CommandList,
} from "@/components/ui/command";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

import { getLanguageDisplayLabel, SPEC_LANGUAGE_INFO } from "../constants/spec-languages";
import type { SpecLanguage } from "../types";

type LanguageComboboxProps = {
  onValueChange: (value: SpecLanguage) => void;
  value: SpecLanguage;
};

export const LanguageCombobox = ({ onValueChange, value }: LanguageComboboxProps) => {
  const t = useTranslations("specView.quotaConfirm");
  const [open, setOpen] = useState(false);

  const selectedLabel = getLanguageDisplayLabel(value);

  return (
    <Popover modal onOpenChange={setOpen} open={open}>
      <PopoverTrigger asChild>
        <Button
          aria-expanded={open}
          aria-label={t("selectLanguage")}
          className="w-full justify-between font-normal"
          id="language-select"
          role="combobox"
          variant="outline"
        >
          <span className="truncate">{selectedLabel}</span>
          <ChevronsUpDown className="ml-2 h-4 w-4 shrink-0 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-[var(--radix-popover-trigger-width)] p-0">
        <Command
          filter={(value, search) => {
            const info = SPEC_LANGUAGE_INFO.find((lang) => lang.code === value);
            if (!info) return 0;
            const searchLower = search.toLowerCase();
            const codeLower = info.code.toLowerCase();
            const nativeLower = info.nativeName.toLowerCase();
            if (codeLower.includes(searchLower) || nativeLower.includes(searchLower)) {
              return 1;
            }
            return 0;
          }}
        >
          <CommandInput placeholder={t("searchLanguage")} />
          <CommandList className="max-h-[200px]">
            <CommandEmpty>{t("noLanguageFound")}</CommandEmpty>
            <CommandGroup>
              {SPEC_LANGUAGE_INFO.map((lang) => (
                <CommandItem
                  key={lang.code}
                  keywords={[lang.nativeName]}
                  onSelect={() => {
                    onValueChange(lang.code);
                    setOpen(false);
                  }}
                  value={lang.code}
                >
                  <Check
                    className={cn(
                      "mr-2 h-4 w-4",
                      value === lang.code ? "opacity-100" : "opacity-0"
                    )}
                  />
                  <span>{getLanguageDisplayLabel(lang.code)}</span>
                </CommandItem>
              ))}
            </CommandGroup>
          </CommandList>
        </Command>
      </PopoverContent>
    </Popover>
  );
};
