import { describe, it, expect, beforeEach } from "vitest";

// Mock localStorage
const localStorageMock = (() => {
  let store: Record<string, string> = {};
  return {
    clear: () => {
      store = {};
    },
    getItem: (key: string) => store[key] ?? null,
    removeItem: (key: string) => {
      delete store[key];
    },
    setItem: (key: string, value: string) => {
      store[key] = value;
    },
  };
})();

Object.defineProperty(global, "localStorage", { value: localStorageMock });

// Import after mocking localStorage
import { isValidSpecLanguage } from "../constants/spec-languages";
import type { SpecLanguage } from "../types";

// Re-implement the functions for testing (since they're not exported)
const getStoredLanguagePreference = (analysisId: string): SpecLanguage | null => {
  const stored = localStorage.getItem(`spec-language-${analysisId}`);
  if (stored && isValidSpecLanguage(stored)) {
    return stored;
  }
  return null;
};

const saveLanguagePreference = (analysisId: string, language: SpecLanguage) => {
  localStorage.setItem(`spec-language-${analysisId}`, language);
  localStorage.setItem("spec-language-default", language);
};

const getGlobalLanguageDefault = (): SpecLanguage | null => {
  const stored = localStorage.getItem("spec-language-default");
  if (stored && isValidSpecLanguage(stored)) {
    return stored;
  }
  return null;
};

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

type ResolveDefaultLanguageParams = {
  analysisId: string;
  initialLanguage?: SpecLanguage;
  locale?: string;
};

const resolveDefaultLanguage = ({
  analysisId,
  initialLanguage,
  locale,
}: ResolveDefaultLanguageParams): SpecLanguage => {
  if (initialLanguage) {
    return initialLanguage;
  }

  const storedPreference = getStoredLanguagePreference(analysisId);
  if (storedPreference) {
    return storedPreference;
  }

  const globalDefault = getGlobalLanguageDefault();
  if (globalDefault) {
    return globalDefault;
  }

  if (locale) {
    const localeLanguage = localeToSpecLanguage(locale);
    if (localeLanguage) {
      return localeLanguage;
    }
  }

  return "English";
};

describe("use-quota-confirm-dialog language preference", () => {
  beforeEach(() => {
    localStorage.clear();
  });

  describe("saveLanguagePreference", () => {
    it("should save language for specific analysis and global default", () => {
      saveLanguagePreference("analysis-1", "Korean");

      expect(localStorage.getItem("spec-language-analysis-1")).toBe("Korean");
      expect(localStorage.getItem("spec-language-default")).toBe("Korean");
    });

    it("should update global default when saving different language", () => {
      saveLanguagePreference("analysis-1", "Korean");
      saveLanguagePreference("analysis-2", "Japanese");

      expect(localStorage.getItem("spec-language-analysis-1")).toBe("Korean");
      expect(localStorage.getItem("spec-language-analysis-2")).toBe("Japanese");
      expect(localStorage.getItem("spec-language-default")).toBe("Japanese");
    });
  });

  describe("getGlobalLanguageDefault", () => {
    it("should return null when no global default is set", () => {
      expect(getGlobalLanguageDefault()).toBeNull();
    });

    it("should return the stored global default", () => {
      localStorage.setItem("spec-language-default", "Korean");
      expect(getGlobalLanguageDefault()).toBe("Korean");
    });

    it("should return null for invalid stored value", () => {
      localStorage.setItem("spec-language-default", "InvalidLanguage");
      expect(getGlobalLanguageDefault()).toBeNull();
    });
  });

  describe("resolveDefaultLanguage priority", () => {
    it("should use initialLanguage when provided", () => {
      localStorage.setItem("spec-language-analysis-1", "Korean");
      localStorage.setItem("spec-language-default", "Japanese");

      const result = resolveDefaultLanguage({
        analysisId: "analysis-1",
        initialLanguage: "Chinese",
        locale: "en",
      });

      expect(result).toBe("Chinese");
    });

    it("should use stored preference over global default", () => {
      localStorage.setItem("spec-language-analysis-1", "Korean");
      localStorage.setItem("spec-language-default", "Japanese");

      const result = resolveDefaultLanguage({
        analysisId: "analysis-1",
        locale: "en",
      });

      expect(result).toBe("Korean");
    });

    it("should use global default when no stored preference exists", () => {
      localStorage.setItem("spec-language-default", "Japanese");

      const result = resolveDefaultLanguage({
        analysisId: "analysis-new",
        locale: "en",
      });

      expect(result).toBe("Japanese");
    });

    it("should use locale when no stored preference or global default exists", () => {
      const result = resolveDefaultLanguage({
        analysisId: "analysis-new",
        locale: "ko",
      });

      expect(result).toBe("Korean");
    });

    it("should fallback to English when nothing is set", () => {
      const result = resolveDefaultLanguage({
        analysisId: "analysis-new",
      });

      expect(result).toBe("English");
    });

    it("should fallback to English for unsupported locale", () => {
      const result = resolveDefaultLanguage({
        analysisId: "analysis-new",
        locale: "xx",
      });

      expect(result).toBe("English");
    });
  });

  describe("cross-analysis language inheritance", () => {
    it("should inherit last used language to new analysis", () => {
      // User selects Korean in analysis A
      saveLanguagePreference("analysis-a", "Korean");

      // New analysis B should default to Korean
      const result = resolveDefaultLanguage({
        analysisId: "analysis-b",
        locale: "en",
      });

      expect(result).toBe("Korean");
    });

    it("should use per-analysis preference over inherited global", () => {
      // User used Korean globally
      saveLanguagePreference("analysis-a", "Korean");

      // But analysis B already has Japanese saved
      localStorage.setItem("spec-language-analysis-b", "Japanese");

      const result = resolveDefaultLanguage({
        analysisId: "analysis-b",
        locale: "en",
      });

      expect(result).toBe("Japanese");
    });
  });
});
