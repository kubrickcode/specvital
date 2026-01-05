"use client";

import { Globe } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";

import { type ConversionLanguage, CONVERSION_LANGUAGES, LANGUAGE_LABELS } from "../../types";

type LanguageSelectorProps = {
  onChange: (language: ConversionLanguage) => void;
  value: ConversionLanguage;
};

export const LanguageSelector = ({ onChange, value }: LanguageSelectorProps) => {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button size="sm" variant="outline">
          <Globe className="mr-2 h-4 w-4" />
          {LANGUAGE_LABELS[value]}
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="max-h-80 overflow-y-auto">
        <DropdownMenuRadioGroup
          onValueChange={(v) => onChange(v as ConversionLanguage)}
          value={value}
        >
          {CONVERSION_LANGUAGES.map((lang) => (
            <DropdownMenuRadioItem key={lang} value={lang}>
              {LANGUAGE_LABELS[lang]}
            </DropdownMenuRadioItem>
          ))}
        </DropdownMenuRadioGroup>
      </DropdownMenuContent>
    </DropdownMenu>
  );
};
