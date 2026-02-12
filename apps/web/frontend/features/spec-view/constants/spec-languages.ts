import type { SpecLanguage } from "../types";

/**
 * Language information with native name for display
 */
export type SpecLanguageInfo = {
  code: SpecLanguage;
  nativeName: string;
};

/**
 * Language data with native names for improved UX
 * Format: { code: "English", nativeName: "English" }
 */
export const SPEC_LANGUAGE_INFO: SpecLanguageInfo[] = [
  { code: "Arabic", nativeName: "العربية" },
  { code: "Chinese", nativeName: "中文" },
  { code: "Czech", nativeName: "Čeština" },
  { code: "Danish", nativeName: "Dansk" },
  { code: "Dutch", nativeName: "Nederlands" },
  { code: "English", nativeName: "English" },
  { code: "Finnish", nativeName: "Suomi" },
  { code: "French", nativeName: "Français" },
  { code: "German", nativeName: "Deutsch" },
  { code: "Greek", nativeName: "Ελληνικά" },
  { code: "Hindi", nativeName: "हिन्दी" },
  { code: "Indonesian", nativeName: "Bahasa Indonesia" },
  { code: "Italian", nativeName: "Italiano" },
  { code: "Japanese", nativeName: "日本語" },
  { code: "Korean", nativeName: "한국어" },
  { code: "Polish", nativeName: "Polski" },
  { code: "Portuguese", nativeName: "Português" },
  { code: "Russian", nativeName: "Русский" },
  { code: "Spanish", nativeName: "Español" },
  { code: "Swedish", nativeName: "Svenska" },
  { code: "Thai", nativeName: "ไทย" },
  { code: "Turkish", nativeName: "Türkçe" },
  { code: "Ukrainian", nativeName: "Українська" },
  { code: "Vietnamese", nativeName: "Tiếng Việt" },
];

/**
 * Backward compatible array of language codes
 * @deprecated Prefer SPEC_LANGUAGE_INFO for display purposes
 */
export const SPEC_LANGUAGES: SpecLanguage[] = SPEC_LANGUAGE_INFO.map((lang) => lang.code);

/**
 * Type guard for valid spec language
 */
export const isValidSpecLanguage = (value: string): value is SpecLanguage => {
  return SPEC_LANGUAGES.includes(value as SpecLanguage);
};

/**
 * Get language info by code
 */
export const getLanguageInfo = (code: SpecLanguage): SpecLanguageInfo | undefined => {
  return SPEC_LANGUAGE_INFO.find((lang) => lang.code === code);
};

/**
 * Get display label for a language (native name + English name)
 * Example: "한국어 (Korean)"
 */
export const getLanguageDisplayLabel = (code: SpecLanguage): string => {
  const info = getLanguageInfo(code);
  if (!info) return code;
  if (info.nativeName === code) return code;
  return `${info.nativeName} (${code})`;
};
