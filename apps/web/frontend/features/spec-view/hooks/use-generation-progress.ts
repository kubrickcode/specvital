"use client";

import { create } from "zustand";

type OpenOptions = {
  analysisId: string;
  onSwitchToBackground?: () => void;
  onViewDocument?: () => void;
};

type GenerationProgressStore = {
  analysisId: string | null;
  bringToForeground: () => void;
  close: () => void;
  isInBackground: boolean;
  isOpen: boolean;
  onSwitchToBackground: (() => void) | null;
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
      onSwitchToBackground: null,
      onViewDocument: null,
    });
  },
  isInBackground: false,
  isOpen: false,
  onSwitchToBackground: null,
  onViewDocument: null,
  open: ({ analysisId, onSwitchToBackground, onViewDocument }: OpenOptions) => {
    set({
      analysisId,
      isInBackground: false,
      isOpen: true,
      onSwitchToBackground: onSwitchToBackground ?? null,
      onViewDocument: onViewDocument ?? null,
    });
  },
  switchToBackground: () => {
    const { onSwitchToBackground } = get();
    set({ isInBackground: true, isOpen: false });
    onSwitchToBackground?.();
  },
}));

export const useGenerationProgress = () => useGenerationProgressStore();
