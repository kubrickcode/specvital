"use client";

import { GitBranch, List } from "lucide-react";
import { useTranslations } from "next-intl";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

import type { ViewMode } from "../types";

type ViewModeToggleProps = {
  onChange: (value: ViewMode) => void;
  value: ViewMode;
};

export const ViewModeToggle = ({ onChange, value }: ViewModeToggleProps) => {
  const t = useTranslations("analyze.viewMode");

  const handleValueChange = (newValue: string) => {
    if (newValue === "list" || newValue === "tree") {
      onChange(newValue);
    }
  };

  return (
    <ToggleGroup
      onValueChange={handleValueChange}
      size="default"
      type="single"
      value={value}
      variant="outline"
    >
      <ToggleGroupItem aria-label={t("list")} value="list">
        <List className="h-4 w-4" />
      </ToggleGroupItem>
      <ToggleGroupItem aria-label={t("tree")} value="tree">
        <GitBranch className="h-4 w-4" />
      </ToggleGroupItem>
    </ToggleGroup>
  );
};
