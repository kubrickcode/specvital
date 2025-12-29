"use client";

import { Globe } from "lucide-react";
import { useParams } from "next/navigation";
import { useTranslations } from "next-intl";
import { useId, useState, useTransition } from "react";

import { Button } from "@/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { isValidLocale, LANGUAGE_NAMES } from "@/i18n/config";
import { usePathname, useRouter } from "@/i18n/navigation";
import { routing } from "@/i18n/routing";

export const LanguageSelector = () => {
  const id = useId();
  const router = useRouter();
  const pathname = usePathname();
  const params = useParams();
  const [isPending, startTransition] = useTransition();
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
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
    <DropdownMenu onOpenChange={setIsDropdownOpen} open={isDropdownOpen}>
      <Tooltip open={isDropdownOpen ? false : undefined}>
        <TooltipTrigger asChild>
          <DropdownMenuTrigger asChild id={id}>
            <Button
              aria-label={t("selectLanguage")}
              disabled={isPending}
              size="header-icon"
              variant="header-action"
            >
              <Globe className="size-4" />
              <span className="sr-only">{t("selectLanguage")}</span>
            </Button>
          </DropdownMenuTrigger>
        </TooltipTrigger>
        <TooltipContent side="bottom" sideOffset={8}>
          {t("selectLanguage")}
        </TooltipContent>
      </Tooltip>
      <DropdownMenuContent align="end">
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
