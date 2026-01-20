import type { SpecLanguage } from "../types";

export const SPEC_LANGUAGES: SpecLanguage[] = [
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

export const isValidSpecLanguage = (value: string): value is SpecLanguage => {
  return SPEC_LANGUAGES.includes(value as SpecLanguage);
};
