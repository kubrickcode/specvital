"use client";

import { useState } from "react";

import { SpecViewContent } from "./spec-view-content";
import { SpecViewError } from "./spec-view-error";
import { SpecViewHeader } from "./spec-view-header";
import { SpecViewSkeleton } from "./spec-view-skeleton";
import { useSpecConversion } from "../../hooks";
import type { ConversionLanguage } from "../../types";

type SpecViewProps = {
  commitSha: string;
  initialLanguage: ConversionLanguage;
  owner: string;
  repo: string;
};

export const SpecView = ({ commitSha, initialLanguage, owner, repo }: SpecViewProps) => {
  const [language, setLanguage] = useState<ConversionLanguage>(initialLanguage);

  const { data, error, isLoading, isRegenerating, regenerate } = useSpecConversion({
    commitSha,
    language,
    owner,
    repo,
  });

  const handleLanguageChange = (newLanguage: ConversionLanguage) => {
    setLanguage(newLanguage);
  };

  if (isLoading || isRegenerating) {
    return <SpecViewSkeleton />;
  }

  if (error) {
    return <SpecViewError error={error} onRetry={regenerate} />;
  }

  if (!data) {
    return <SpecViewSkeleton />;
  }

  return (
    <div>
      <SpecViewHeader
        language={language}
        onLanguageChange={handleLanguageChange}
        summary={data.summary}
      />
      <SpecViewContent files={data.data} />
    </div>
  );
};
