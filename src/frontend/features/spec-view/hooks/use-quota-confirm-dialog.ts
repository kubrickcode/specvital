"use client";

import { useSyncExternalStore } from "react";

import type { UsageStatusResponse } from "@/features/account/api/usage";

import type { SpecLanguage } from "../types";

type QuotaConfirmDialogStore = {
  analysisId: string | null;
  estimatedCost: number | null;
  isOpen: boolean;
  isRegenerate: boolean;
  listeners: Set<() => void>;
  onConfirm: ((language: SpecLanguage, isForceRegenerate: boolean) => void) | null;
  selectedLanguage: SpecLanguage;
  usage: UsageStatusResponse | null;
};

const store: QuotaConfirmDialogStore = {
  analysisId: null,
  estimatedCost: null,
  isOpen: false,
  isRegenerate: false,
  listeners: new Set(),
  onConfirm: null,
  selectedLanguage: "English",
  usage: null,
};

const notifyListeners = () => {
  store.listeners.forEach((listener) => listener());
};

const subscribe = (listener: () => void) => {
  store.listeners.add(listener);
  return () => store.listeners.delete(listener);
};

const getSnapshot = () => store.isOpen;

const getServerSnapshot = () => false;

const getUsageSnapshot = () => store.usage;

const getOnConfirmSnapshot = () => store.onConfirm;

const getEstimatedCostSnapshot = () => store.estimatedCost;

const getSelectedLanguageSnapshot = () => store.selectedLanguage;

const getAnalysisIdSnapshot = () => store.analysisId;

const getIsRegenerateSnapshot = () => store.isRegenerate;

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
 * Save language preference for a specific analysis
 */
const saveLanguagePreference = (analysisId: string, language: SpecLanguage) => {
  if (typeof window === "undefined") return;
  localStorage.setItem(`spec-language-${analysisId}`, language);
};

const SPEC_LANGUAGES: SpecLanguage[] = [
  "Arabic",
  "Chinese",
  "Czech",
  "Danish",
  "Dutch",
  "English",
  "Finnish",
  "French",
  "German",
  "Greek",
  "Hindi",
  "Indonesian",
  "Italian",
  "Japanese",
  "Korean",
  "Polish",
  "Portuguese",
  "Russian",
  "Spanish",
  "Swedish",
  "Thai",
  "Turkish",
  "Ukrainian",
  "Vietnamese",
];

const isValidSpecLanguage = (value: string): value is SpecLanguage => {
  return SPEC_LANGUAGES.includes(value as SpecLanguage);
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

type OpenOptions = {
  analysisId: string;
  estimatedCost?: number;
  isRegenerate?: boolean;
  locale?: string;
  onConfirm: (language: SpecLanguage, isForceRegenerate: boolean) => void;
  usage: UsageStatusResponse | null;
};

const open = ({
  analysisId,
  estimatedCost,
  isRegenerate = false,
  locale,
  onConfirm,
  usage,
}: OpenOptions) => {
  if (!store.isOpen) {
    // Determine default language: stored preference > locale > English
    let defaultLanguage: SpecLanguage = "English";
    const storedPreference = getStoredLanguagePreference(analysisId);
    if (storedPreference) {
      defaultLanguage = storedPreference;
    } else if (locale) {
      const localeLanguage = localeToSpecLanguage(locale);
      if (localeLanguage) {
        defaultLanguage = localeLanguage;
      }
    }

    store.isOpen = true;
    store.isRegenerate = isRegenerate;
    store.usage = usage;
    store.onConfirm = onConfirm;
    store.estimatedCost = estimatedCost ?? null;
    store.analysisId = analysisId;
    store.selectedLanguage = defaultLanguage;
    notifyListeners();
  }
};

const close = () => {
  if (store.isOpen) {
    store.isOpen = false;
    store.isRegenerate = false;
    store.usage = null;
    store.onConfirm = null;
    store.estimatedCost = null;
    store.analysisId = null;
    store.selectedLanguage = "English";
    notifyListeners();
  }
};

const setSelectedLanguage = (language: SpecLanguage) => {
  if (store.selectedLanguage !== language) {
    store.selectedLanguage = language;
    // Save preference for this analysis
    if (store.analysisId) {
      saveLanguagePreference(store.analysisId, language);
    }
    notifyListeners();
  }
};

const confirm = () => {
  if (store.onConfirm) {
    store.onConfirm(store.selectedLanguage, store.isRegenerate);
  }
  close();
};

export const useQuotaConfirmDialog = () => {
  const isOpen = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);
  const usage = useSyncExternalStore(subscribe, getUsageSnapshot, () => null);
  const onConfirm = useSyncExternalStore(subscribe, getOnConfirmSnapshot, () => null);
  const estimatedCost = useSyncExternalStore(subscribe, getEstimatedCostSnapshot, () => null);
  const selectedLanguage = useSyncExternalStore(
    subscribe,
    getSelectedLanguageSnapshot,
    () => "English" as SpecLanguage
  );
  const analysisId = useSyncExternalStore(subscribe, getAnalysisIdSnapshot, () => null);
  const isRegenerate = useSyncExternalStore(subscribe, getIsRegenerateSnapshot, () => false);

  const onOpenChange = (open: boolean) => {
    if (!open) {
      close();
    }
  };

  return {
    analysisId,
    close,
    confirm,
    estimatedCost,
    isOpen,
    isRegenerate,
    onConfirm,
    onOpenChange,
    open,
    selectedLanguage,
    setSelectedLanguage,
    specLanguages: SPEC_LANGUAGES,
    usage,
  };
};
