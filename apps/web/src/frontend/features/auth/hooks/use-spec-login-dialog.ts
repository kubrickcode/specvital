"use client";

import { create } from "zustand";

type SpecLoginDialogStore = {
  close: () => void;
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
  open: () => void;
};

const useSpecLoginDialogStore = create<SpecLoginDialogStore>((set) => ({
  close: () => set({ isOpen: false }),
  isOpen: false,
  onOpenChange: (isOpen) => set({ isOpen }),
  open: () => set({ isOpen: true }),
}));

export const useSpecLoginDialog = () => useSpecLoginDialogStore();
