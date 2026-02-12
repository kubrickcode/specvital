"use client";

import { Bookmark } from "lucide-react";
import { useTranslations } from "next-intl";

import { Toggle } from "@/components/ui/toggle";

import { useBookmarkFilter } from "../hooks/use-bookmark-filter";

export const BookmarkToggle = () => {
  const t = useTranslations("dashboard.filter");
  const { bookmarkOnly, setBookmarkOnly } = useBookmarkFilter();

  const handleToggle = (pressed: boolean) => {
    setBookmarkOnly(pressed ? true : null);
  };

  return (
    <Toggle
      aria-label={t("bookmarkedLabel")}
      className="h-11 gap-2 px-3 sm:h-9"
      onPressedChange={handleToggle}
      pressed={bookmarkOnly}
      variant="outline"
    >
      <Bookmark
        aria-hidden="true"
        className={bookmarkOnly ? "fill-yellow-400 text-yellow-400" : ""}
      />
      <span className="hidden sm:inline">{t("bookmarked")}</span>
    </Toggle>
  );
};
