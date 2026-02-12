"use client";

import { AlertCircle, ArrowRight, CheckCircle, Github, HelpCircle, Loader2 } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState, useTransition } from "react";

import { Alert, AlertDescription } from "@/components/ui/alert";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { ResponsiveTooltip } from "@/components/ui/responsive-tooltip";
import { useRouter } from "@/i18n/navigation";
import { useMediaQuery } from "@/lib/hooks";
import { cn } from "@/lib/utils";

import { useDebouncedValidation } from "../hooks/use-debounced-validation";
import { getInputFeedback, isValidGitHubUrl, parseGitHubUrl } from "../lib";

type UrlInputFormProps = {
  onSuccess?: () => void;
};

export const UrlInputForm = ({ onSuccess }: UrlInputFormProps) => {
  const router = useRouter();
  const t = useTranslations("home");
  const isDesktop = useMediaQuery("(min-width: 640px)");
  const [error, setError] = useState<string | null>(null);
  const [url, setUrl] = useState("");
  const [isPending, startTransition] = useTransition();

  const placeholder = isDesktop ? t("inputPlaceholder") : t("inputPlaceholderShort");

  const validationState = useDebouncedValidation(url, isValidGitHubUrl, {
    delay: 500,
    minLength: 3,
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

    onSuccess?.();
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

  const renderInputHint = () => {
    const feedback = getInputFeedback(url);

    if (feedback.type === "shorthand" || feedback.type === "deeplink") {
      return (
        <p aria-live="polite" className="text-xs text-muted-foreground" id="url-hint">
          → {feedback.normalized}
        </p>
      );
    }

    return null;
  };

  return (
    <form className="w-full space-y-3" onSubmit={handleSubmit}>
      <label className="sr-only" htmlFor="github-url">
        {t("inputLabel")}
      </label>
      <div className="flex flex-col sm:flex-row gap-2">
        <div className="relative flex-1">
          <Github
            aria-hidden="true"
            className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground"
          />
          <Input
            aria-describedby={
              error ? "url-error" : validationState !== "idle" ? "url-validation-status" : undefined
            }
            aria-invalid={!!error || validationState === "invalid"}
            aria-label={t("inputLabel")}
            className={cn(
              "pl-10 pr-16 h-11 sm:h-10",
              validationState === "valid" && "border-green-500/50 focus-visible:ring-green-500/20",
              validationState === "invalid" && "border-amber-500/50 focus-visible:ring-amber-500/20"
            )}
            disabled={isPending}
            id="github-url"
            onChange={handleChange}
            placeholder={placeholder}
            type="text"
            value={url}
          />
          <div className="absolute right-3 top-1/2 -translate-y-1/2 flex items-center gap-1.5">
            <ResponsiveTooltip
              content={
                <div className="space-y-2">
                  <p className="font-medium text-sm">{t("supportedFormats.title")}</p>
                  <ul className="text-xs space-y-1 list-disc list-inside">
                    <li>https://github.com/owner/repo</li>
                    <li>github.com/owner/repo</li>
                    <li>owner/repo</li>
                  </ul>
                  <p className="text-xs text-background/70">{t("supportedFormats.note")}</p>
                </div>
              }
              contentClassName="max-w-xs"
              side="top"
              sideOffset={8}
            >
              <button
                aria-label={t("supportedFormats.label")}
                className="text-muted-foreground hover:text-foreground transition-colors"
                type="button"
              >
                <HelpCircle className="h-4 w-4" />
              </button>
            </ResponsiveTooltip>
            {validationState === "valid" && (
              <CheckCircle aria-hidden="true" className="h-4 w-4 text-green-500" />
            )}
            {validationState === "invalid" && (
              <AlertCircle aria-hidden="true" className="h-4 w-4 text-amber-500" />
            )}
          </div>
          {validationState !== "idle" && (
            <span aria-live="polite" className="sr-only" id="url-validation-status">
              {validationState === "valid" ? "Valid GitHub URL" : "Invalid GitHub URL format"}
            </span>
          )}
        </div>
        <Button
          aria-label={t("analyzeButton")}
          className="w-full sm:w-auto sm:min-w-[120px] h-11 sm:h-10"
          disabled={isPending}
          size="lg"
          type="submit"
          variant="cta"
        >
          {isPending ? (
            <>
              <Loader2 aria-hidden="true" className="h-4 w-4 animate-spin" />
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
      {renderInputHint()}
      {error && (
        <Alert id="url-error" variant="destructive">
          <AlertCircle className="h-4 w-4" />
          <AlertDescription>{error}</AlertDescription>
        </Alert>
      )}
    </form>
  );
};
