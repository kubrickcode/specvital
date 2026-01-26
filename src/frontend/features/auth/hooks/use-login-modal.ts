"use client";

import { create } from "zustand";

type LoginModalStore = {
  close: () => void;
  isOpen: boolean;
  onOpenChange: (isOpen: boolean) => void;
  open: () => void;
};

const useLoginModalStore = create<LoginModalStore>((set) => ({
  close: () => set({ isOpen: false }),
  isOpen: false,
  onOpenChange: (isOpen) => set({ isOpen }),
  open: () => set({ isOpen: true }),
}));

export const useLoginModal = () => useLoginModalStore();
