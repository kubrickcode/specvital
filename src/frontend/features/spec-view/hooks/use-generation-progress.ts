"use client";

import { create } from "zustand";

type OpenOptions = {
  analysisId: string;
  onViewDocument?: () => void;
};

type GenerationProgressStore = {
  analysisId: string | null;
  bringToForeground: () => void;
  close: () => void;
  isInBackground: boolean;
  isOpen: boolean;
  onViewDocument: (() => void) | null;
  open: (options: OpenOptions) => void;
  switchToBackground: () => void;
};

const useGenerationProgressStore = create<GenerationProgressStore>((set, get) => ({
  analysisId: null,
  bringToForeground: () => {
    if (get().isInBackground) {
      set({ isInBackground: false, isOpen: true });
    }
  },
  close: () => {
    set({
      analysisId: null,
      isInBackground: false,
      isOpen: false,
      onViewDocument: null,
    });
  },
  isInBackground: false,
  isOpen: false,
  onViewDocument: null,
  open: ({ analysisId, onViewDocument }: OpenOptions) => {
    set({
      analysisId,
      isInBackground: false,
      isOpen: true,
      onViewDocument: onViewDocument ?? null,
    });
  },
  switchToBackground: () => {
    set({ isInBackground: true, isOpen: false });
  },
}));

export const useGenerationProgress = () => useGenerationProgressStore();
