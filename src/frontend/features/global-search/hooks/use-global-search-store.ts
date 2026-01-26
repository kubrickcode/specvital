"use client";

import { create } from "zustand";

type GlobalSearchStore = {
  close: () => void;
  isOpen: boolean;
  open: () => void;
  toggle: () => void;
};

const useGlobalSearchStoreInternal = create<GlobalSearchStore>((set, get) => ({
  close: () => set({ isOpen: false }),
  isOpen: false,
  open: () => set({ isOpen: true }),
  toggle: () => set({ isOpen: !get().isOpen }),
}));

export const globalSearchStore = {
  close: () => useGlobalSearchStoreInternal.getState().close(),
  open: () => useGlobalSearchStoreInternal.getState().open(),
  toggle: () => useGlobalSearchStoreInternal.getState().toggle(),
};

export const useGlobalSearchStore = () => useGlobalSearchStoreInternal();
