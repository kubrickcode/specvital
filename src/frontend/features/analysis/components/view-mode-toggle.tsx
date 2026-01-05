"use client";

import { BookOpen, GitBranch, List } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { ToggleGroup, ToggleGroupItem } from "@/components/ui/toggle-group";

import type { ConversionLanguage, ViewMode } from "../types";
import { VIEW_MODES } from "../types";
import { ConsentDialog } from "./spec-view";

type ViewModeToggleProps = {
  onChange: (value: ViewMode, language?: ConversionLanguage) => void;
  value: ViewMode;
};

export const ViewModeToggle = ({ onChange, value }: ViewModeToggleProps) => {
  const t = useTranslations("analyze.viewMode");
  const [showConsentDialog, setShowConsentDialog] = useState(false);

  const handleValueChange = (newValue: string) => {
    if (!VIEW_MODES.includes(newValue as ViewMode)) return;

    if (newValue === "spec") {
      setShowConsentDialog(true);
      return;
    }

    onChange(newValue as ViewMode);
  };

  const handleConsent = (language: ConversionLanguage) => {
    setShowConsentDialog(false);
    onChange("spec", language);
  };

  const handleCancelConsent = () => {
    setShowConsentDialog(false);
  };

  return (
    <>
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
        <ToggleGroupItem aria-label={t("spec")} value="spec">
          <BookOpen className="h-4 w-4" />
        </ToggleGroupItem>
      </ToggleGroup>

      <ConsentDialog
        onCancel={handleCancelConsent}
        onConsent={handleConsent}
        open={showConsentDialog}
      />
    </>
  );
};
