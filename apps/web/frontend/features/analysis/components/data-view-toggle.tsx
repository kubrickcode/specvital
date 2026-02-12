"use client";

import { GitBranch, List } from "lucide-react";
import { useTranslations } from "next-intl";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";

import type { DataViewMode } from "../types/data-view-mode";
import { DATA_VIEW_MODES } from "../types/data-view-mode";

type DataViewToggleProps = {
  onChange: (value: DataViewMode) => void;
  value: DataViewMode;
};

export const DataViewToggle = ({ onChange, value }: DataViewToggleProps) => {
  const t = useTranslations("analyze.viewMode");

  const handleValueChange = (newValue: string) => {
    if (!DATA_VIEW_MODES.includes(newValue as DataViewMode)) return;
    onChange(newValue as DataViewMode);
  };

  return (
    <ToggleGroup
      onValueChange={handleValueChange}
      size="default"
      type="single"
      value={value}
      variant="outline"
    >
      <Tooltip>
        <TooltipTrigger asChild>
          <ToggleGroupItem aria-label={t("list")} value="list">
            <List className="h-4 w-4" />
          </ToggleGroupItem>
        </TooltipTrigger>
        <TooltipContent>{t("list")}</TooltipContent>
      </Tooltip>
      <Tooltip>
        <TooltipTrigger asChild>
          <ToggleGroupItem aria-label={t("tree")} value="tree">
            <GitBranch className="h-4 w-4" />
          </ToggleGroupItem>
        </TooltipTrigger>
        <TooltipContent>{t("tree")}</TooltipContent>
      </Tooltip>
    </ToggleGroup>
  );
};
