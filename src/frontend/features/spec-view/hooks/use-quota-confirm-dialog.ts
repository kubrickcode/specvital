"use client";

import { useSyncExternalStore } from "react";

import type { UsageStatusResponse } from "@/features/account/api/usage";

import { SPEC_LANGUAGES, isValidSpecLanguage } from "../constants/spec-languages";
import type { SpecLanguage } from "../types";

type QuotaConfirmDialogStore = {
  analysisId: string | null;
  estimatedCost: number | null;
  forceRegenerate: boolean;
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
  forceRegenerate: false,
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

const getForceRegenerateSnapshot = () => store.forceRegenerate;

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

type OpenOptions = {
  analysisId: string;
  estimatedCost?: number;
  initialLanguage?: SpecLanguage;
  isRegenerate?: boolean;
  locale?: string;
  onConfirm: (language: SpecLanguage, isForceRegenerate: boolean) => void;
  usage: UsageStatusResponse | null;
};

const open = ({
  analysisId,
  estimatedCost,
  initialLanguage,
  isRegenerate = false,
  locale,
  onConfirm,
  usage,
}: OpenOptions) => {
  if (!store.isOpen) {
    // Determine default language: initialLanguage > stored preference > global default > locale > English
    let defaultLanguage: SpecLanguage = "English";
    if (initialLanguage) {
      defaultLanguage = initialLanguage;
    } else {
      const storedPreference = getStoredLanguagePreference(analysisId);
      if (storedPreference) {
        defaultLanguage = storedPreference;
      } else {
        const globalDefault = getGlobalLanguageDefault();
        if (globalDefault) {
          defaultLanguage = globalDefault;
        } else if (locale) {
          const localeLanguage = localeToSpecLanguage(locale);
          if (localeLanguage) {
            defaultLanguage = localeLanguage;
          }
        }
      }
    }

    store.isOpen = true;
    store.isRegenerate = isRegenerate;
    store.forceRegenerate = isRegenerate; // Regenerate always forces, new generation defaults to cache
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
    store.forceRegenerate = false;
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

const setForceRegenerate = (value: boolean) => {
  if (store.forceRegenerate !== value) {
    store.forceRegenerate = value;
    notifyListeners();
  }
};

const confirm = () => {
  // Capture values before close() resets them
  const callback = store.onConfirm;
  const language = store.selectedLanguage;
  const forceRegen = store.forceRegenerate;
  close();
  callback?.(language, forceRegen);
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
  const forceRegenerate = useSyncExternalStore(subscribe, getForceRegenerateSnapshot, () => false);

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
    forceRegenerate,
    isOpen,
    isRegenerate,
    onConfirm,
    onOpenChange,
    open,
    selectedLanguage,
    setForceRegenerate,
    setSelectedLanguage,
    specLanguages: SPEC_LANGUAGES,
    usage,
  };
};
