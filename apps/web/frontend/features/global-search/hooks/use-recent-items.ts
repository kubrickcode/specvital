"use client";

import { create } from "zustand";

import {
  addRecentItem,
  clearRecentItems,
  loadRecentItems,
  STORAGE_KEY,
  type RecentItem,
} from "../lib/recent-items";

type RecentItemsStore = {
  addItem: (item: Omit<RecentItem, "timestamp" | "type">) => void;
  clearItems: () => void;
  recentItems: RecentItem[];
};

const useRecentItemsStore = create<RecentItemsStore>((set) => ({
  addItem: (item) => {
    addRecentItem(item);
    set({ recentItems: loadRecentItems() });
  },
  clearItems: () => {
    clearRecentItems();
    set({ recentItems: [] });
  },
  recentItems: typeof window !== "undefined" ? loadRecentItems() : [],
}));

// Sync across tabs via storage event
if (typeof window !== "undefined") {
  window.addEventListener("storage", (event) => {
    if (event.key === STORAGE_KEY) {
      useRecentItemsStore.setState({ recentItems: loadRecentItems() });
    }
  });
}

export const useRecentItems = (): RecentItemsStore => useRecentItemsStore();
