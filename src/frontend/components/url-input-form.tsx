"use client";

import { useState } from "react";
import { ArrowRight } from "lucide-react";
import { useTranslations } from "next-intl";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { parseGitHubUrl } from "@/lib/github-url";
import { useRouter } from "@/i18n/navigation";

export const UrlInputForm = () => {
  const router = useRouter();
  const t = useTranslations("home");
  const [error, setError] = useState<string | null>(null);
  const [url, setUrl] = useState("");

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const result = parseGitHubUrl(url);

    if (!result.success) {
      setError(result.error);
      return;
    }

    setError(null);
    const { owner, repo } = result.data;
    router.push(`/analyze/${owner}/${repo}`);
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
    if (error) {
      setError(null);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="w-full space-y-2">
      <label htmlFor="github-url" className="sr-only">
        {t("inputLabel")}
      </label>
      <div className="flex gap-2">
        <Input
          id="github-url"
          type="url"
          placeholder={t("inputPlaceholder")}
          value={url}
          onChange={handleChange}
          aria-invalid={!!error}
          aria-describedby={error ? "url-error" : undefined}
          aria-label={t("inputLabel")}
          className="flex-1"
        />
        <Button type="submit" size="lg" aria-label={t("analyzeButton")}>
          {t("analyzeButton")}
          <ArrowRight aria-hidden="true" />
        </Button>
      </div>
      {error && (
        <p id="url-error" role="alert" className="text-sm text-destructive">
          {error}
        </p>
      )}
    </form>
  );
};
