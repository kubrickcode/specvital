"use client";

import { create } from "zustand";

import type { UsageStatusResponse } from "@/features/account/api/usage";

import { isValidSpecLanguage } from "../constants/spec-languages";
import type { SpecGenerationMode, SpecLanguage } from "../types";

type OpenOptions = {
  /**
   * Commit SHA of the current analysis being viewed.
   * Compared with documentCommitSha to determine if cache would be ineffective.
   */
  analysisCommitSha?: string;
  analysisId: string;
  /**
   * Commit SHA of the currently viewed spec document.
   * Used to detect if regeneration would produce identical results.
   */
  documentCommitSha?: string;
  estimatedCost?: number;
  initialLanguage?: SpecLanguage;
  isRegenerate?: boolean;
  locale?: string;
  onConfirm: (language: SpecLanguage, mode: SpecGenerationMode) => void;
  usage: UsageStatusResponse | null;
};

type QuotaConfirmDialogStore = {
  analysisId: string | null;
  close: () => void;
  confirm: () => void;
  estimatedCost: number | null;
  forceRegenerate: boolean;
  isOpen: boolean;
  isRegenerate: boolean;
  /**
   * True when regenerating the same commit - cache would be ineffective.
   * When true, "Use Cache" option is disabled and only "Fresh" is allowed.
   */
  isSameCommit: boolean;
  onConfirm: ((language: SpecLanguage, mode: SpecGenerationMode) => void) | null;
  onOpenChange: (open: boolean) => void;
  open: (options: OpenOptions) => void;
  regeneratingLanguage: SpecLanguage | null;
  selectedLanguage: SpecLanguage;
  setForceRegenerate: (value: boolean) => void;
  setSelectedLanguage: (language: SpecLanguage) => void;
  usage: UsageStatusResponse | null;
};

/**
 * Get the stored language preference for a specific analysis
 */
const getStoredLanguagePreference = (analysisId: string): SpecLanguage | null => {
  if (typeof window === "undefined") return null;
  const stored = localStorage.getItem(`spec-language-${analysisId}`);
  if (stored && isValidSpecLanguage(stored)) {
    return stored;
  }
  return null;
};

/**
 * Save language preference for a specific analysis and update global default
 */
const saveLanguagePreference = (analysisId: string, language: SpecLanguage) => {
  if (typeof window === "undefined") return;
  localStorage.setItem(`spec-language-${analysisId}`, language);
  localStorage.setItem("spec-language-default", language);
};

/**
 * Get the global default language (last used language across all analyses)
 */
const getGlobalLanguageDefault = (): SpecLanguage | null => {
  if (typeof window === "undefined") return null;
  const stored = localStorage.getItem("spec-language-default");
  if (stored && isValidSpecLanguage(stored)) {
    return stored;
  }
  return null;
};

/**
 * Map locale to SpecLanguage
 */
const localeToSpecLanguage = (locale: string): SpecLanguage | null => {
  const localeMap: Record<string, SpecLanguage> = {
    ar: "Arabic",
    cs: "Czech",
    da: "Danish",
    de: "German",
    el: "Greek",
    en: "English",
    es: "Spanish",
    fi: "Finnish",
    fr: "French",
    hi: "Hindi",
    id: "Indonesian",
    it: "Italian",
    ja: "Japanese",
    ko: "Korean",
    nl: "Dutch",
    pl: "Polish",
    pt: "Portuguese",
    ru: "Russian",
    sv: "Swedish",
    th: "Thai",
    tr: "Turkish",
    uk: "Ukrainian",
    vi: "Vietnamese",
    zh: "Chinese",
  };
  return localeMap[locale] ?? null;
};

const resolveDefaultLanguage = (
  analysisId: string,
  initialLanguage?: SpecLanguage,
  locale?: string
): SpecLanguage => {
  if (initialLanguage) return initialLanguage;

  const storedPreference = getStoredLanguagePreference(analysisId);
  if (storedPreference) return storedPreference;

  const globalDefault = getGlobalLanguageDefault();
  if (globalDefault) return globalDefault;

  if (locale) {
    const localeLanguage = localeToSpecLanguage(locale);
    if (localeLanguage) return localeLanguage;
  }

  return "English";
};

const INITIAL_STATE = {
  analysisId: null as string | null,
  estimatedCost: null as number | null,
  forceRegenerate: false,
  isOpen: false,
  isRegenerate: false,
  isSameCommit: false,
  onConfirm: null as ((language: SpecLanguage, mode: SpecGenerationMode) => void) | null,
  regeneratingLanguage: null as SpecLanguage | null,
  selectedLanguage: "English" as SpecLanguage,
  usage: null as UsageStatusResponse | null,
};

const useQuotaConfirmDialogStore = create<QuotaConfirmDialogStore>((set, get) => ({
  ...INITIAL_STATE,
  close: () => set(INITIAL_STATE),
  confirm: () => {
    const { forceRegenerate, isRegenerate, onConfirm, selectedLanguage } = get();
    set(INITIAL_STATE);

    const mode: SpecGenerationMode = isRegenerate
      ? forceRegenerate
        ? "regenerate_fresh"
        : "regenerate_cached"
      : "initial";

    onConfirm?.(selectedLanguage, mode);
  },
  onOpenChange: (open) => {
    if (!open) set(INITIAL_STATE);
  },
  open: ({
    analysisCommitSha,
    analysisId,
    documentCommitSha,
    estimatedCost,
    initialLanguage,
    isRegenerate = false,
    locale,
    onConfirm,
    usage,
  }) => {
    if (get().isOpen) return;

    // Detect if regenerating from the same commit - cache would be ineffective
    const isSameCommit =
      isRegenerate &&
      Boolean(analysisCommitSha) &&
      Boolean(documentCommitSha) &&
      analysisCommitSha === documentCommitSha;

    set({
      analysisId,
      estimatedCost: estimatedCost ?? null,
      // Force fresh mode when same commit (cache is pointless)
      forceRegenerate: isRegenerate || isSameCommit,
      isOpen: true,
      isRegenerate,
      isSameCommit,
      onConfirm,
      regeneratingLanguage: isRegenerate && initialLanguage ? initialLanguage : null,
      selectedLanguage: resolveDefaultLanguage(analysisId, initialLanguage, locale),
      usage,
    });
  },
  setForceRegenerate: (value) => set({ forceRegenerate: value }),
  setSelectedLanguage: (language) => {
    set({ selectedLanguage: language });
    const { analysisId } = get();
    if (analysisId) {
      saveLanguagePreference(analysisId, language);
    }
  },
}));

export const useQuotaConfirmDialog = () => useQuotaConfirmDialogStore();
