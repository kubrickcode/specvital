import { describe, expect, it } from "vitest";

import {
  getLanguageDisplayLabel,
  getLanguageInfo,
  isValidSpecLanguage,
  SPEC_LANGUAGE_INFO,
  SPEC_LANGUAGES,
} from "./spec-languages";

describe("spec-languages", () => {
  describe("SPEC_LANGUAGE_INFO", () => {
    it("should have 24 languages", () => {
      expect(SPEC_LANGUAGE_INFO).toHaveLength(24);
    });

    it("should have unique codes", () => {
      const codes = SPEC_LANGUAGE_INFO.map((lang) => lang.code);
      const uniqueCodes = new Set(codes);
      expect(uniqueCodes.size).toBe(codes.length);
    });

    it("should have native names for all languages", () => {
      for (const lang of SPEC_LANGUAGE_INFO) {
        expect(lang.nativeName).toBeTruthy();
        expect(lang.nativeName.length).toBeGreaterThan(0);
      }
    });
  });

  describe("SPEC_LANGUAGES", () => {
    it("should be derived from SPEC_LANGUAGE_INFO", () => {
      expect(SPEC_LANGUAGES).toEqual(SPEC_LANGUAGE_INFO.map((lang) => lang.code));
    });

    it("should maintain backward compatibility with 24 languages", () => {
      expect(SPEC_LANGUAGES).toHaveLength(24);
      expect(SPEC_LANGUAGES).toContain("English");
      expect(SPEC_LANGUAGES).toContain("Korean");
      expect(SPEC_LANGUAGES).toContain("Japanese");
    });
  });

  describe("isValidSpecLanguage", () => {
    it("should return true for valid languages", () => {
      expect(isValidSpecLanguage("English")).toBe(true);
      expect(isValidSpecLanguage("Korean")).toBe(true);
      expect(isValidSpecLanguage("Japanese")).toBe(true);
    });

    it("should return false for invalid languages", () => {
      expect(isValidSpecLanguage("InvalidLanguage")).toBe(false);
      expect(isValidSpecLanguage("")).toBe(false);
      expect(isValidSpecLanguage("english")).toBe(false); // case sensitive
    });
  });

  describe("getLanguageInfo", () => {
    it("should return language info for valid code", () => {
      const korean = getLanguageInfo("Korean");
      expect(korean).toEqual({ code: "Korean", nativeName: "한국어" });

      const japanese = getLanguageInfo("Japanese");
      expect(japanese).toEqual({ code: "Japanese", nativeName: "日本語" });
    });

    it("should return undefined for invalid code", () => {
      // @ts-expect-error - testing invalid input
      expect(getLanguageInfo("InvalidLanguage")).toBeUndefined();
    });
  });

  describe("getLanguageDisplayLabel", () => {
    it("should return native name with English in parentheses", () => {
      expect(getLanguageDisplayLabel("Korean")).toBe("한국어 (Korean)");
      expect(getLanguageDisplayLabel("Japanese")).toBe("日本語 (Japanese)");
      expect(getLanguageDisplayLabel("Chinese")).toBe("中文 (Chinese)");
    });

    it("should return just the code when native name equals code", () => {
      expect(getLanguageDisplayLabel("English")).toBe("English");
    });

    it("should handle special characters in native names", () => {
      expect(getLanguageDisplayLabel("Arabic")).toBe("العربية (Arabic)");
      expect(getLanguageDisplayLabel("Thai")).toBe("ไทย (Thai)");
      expect(getLanguageDisplayLabel("Greek")).toBe("Ελληνικά (Greek)");
    });
  });
});
