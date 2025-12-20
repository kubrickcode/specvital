"use client";

import { useState, useTransition } from "react";
import { AlertCircle, ArrowRight, CheckCircle, Github, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";
import { useDebouncedValidation } from "../hooks/use-debounced-validation";
import { isValidGitHubUrl, parseGitHubUrl } from "../lib";
import { useRouter } from "@/i18n/navigation";

export const UrlInputForm = () => {
  const router = useRouter();
  const t = useTranslations("home");
  const [error, setError] = useState<string | null>(null);
  const [url, setUrl] = useState("");
  const [isPending, startTransition] = useTransition();

  const validationState = useDebouncedValidation(url, isValidGitHubUrl, {
    delay: 500,
    minLength: 10,
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();

    const result = parseGitHubUrl(url);

    if (!result.success) {
      setError(result.error);
      return;
    }

    setError(null);
    const { owner, repo } = result.data;

    startTransition(() => {
      router.push(`/analyze/${owner}/${repo}`);
    });
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setUrl(e.target.value);
    if (error) {
      setError(null);
    }
  };

  const renderValidationIcon = () => {
    if (validationState === "valid") {
      return (
        <CheckCircle
          className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-green-500"
          aria-hidden="true"
        />
      );
    }
    if (validationState === "invalid") {
      return (
        <AlertCircle
          className="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4 text-amber-500"
          aria-hidden="true"
        />
      );
    }
    return null;
  };

  return (
    <form onSubmit={handleSubmit} className="w-full space-y-3">
      <label htmlFor="github-url" className="sr-only">
        {t("inputLabel")}
      </label>
      <div className="flex gap-2">
        <div className="relative flex-1">
          <Github
            className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
            aria-hidden="true"
          />
          <Input
            id="github-url"
            type="url"
            placeholder={t("inputPlaceholder")}
            value={url}
            onChange={handleChange}
            disabled={isPending}
            aria-invalid={!!error || validationState === "invalid"}
            aria-describedby={
              error ? "url-error" : validationState !== "idle" ? "url-validation-status" : undefined
            }
            aria-label={t("inputLabel")}
            className={cn(
              "pl-10 pr-10 h-11 sm:h-10",
              validationState === "valid" && "border-green-500/50 focus-visible:ring-green-500/20",
              validationState === "invalid" && "border-amber-500/50 focus-visible:ring-amber-500/20"
            )}
          />
          {renderValidationIcon()}
          {validationState !== "idle" && (
            <span id="url-validation-status" className="sr-only" aria-live="polite">
              {validationState === "valid" ? "Valid GitHub URL" : "Invalid GitHub URL format"}
            </span>
          )}
        </div>
        <Button
          type="submit"
          size="lg"
          disabled={isPending}
          aria-label={t("analyzeButton")}
          className="min-w-[120px]"
        >
          {isPending ? (
            <>
              <Loader2 className="h-4 w-4 animate-spin" aria-hidden="true" />
              <span className="sr-only">Loading...</span>
            </>
          ) : (
            <>
              {t("analyzeButton")}
              <ArrowRight aria-hidden="true" />
            </>
          )}
        </Button>
      </div>
      {error && (
        <Alert variant="destructive" id="url-error">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </form>
  );
};
