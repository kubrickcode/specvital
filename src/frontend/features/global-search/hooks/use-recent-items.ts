"use client";

import { useSyncExternalStore } from "react";

import {
  addRecentItem,
  clearRecentItems,
  loadRecentItems,
  STORAGE_KEY,
  type RecentItem,
} from "../lib/recent-items";

type Listener = () => void;

const listeners = new Set<Listener>();
let cachedItems: RecentItem[] | null = null;

const notifyListeners = () => {
  cachedItems = null;
  for (const listener of listeners) {
    listener();
  }
};

if (typeof window !== "undefined") {
  window.addEventListener("storage", (event) => {
    if (event.key === STORAGE_KEY) {
      notifyListeners();
    }
  });
}

const subscribe = (listener: Listener): (() => void) => {
  listeners.add(listener);
  return () => listeners.delete(listener);
};

const getSnapshot = (): RecentItem[] => {
  if (cachedItems === null) {
    cachedItems = loadRecentItems();
  }
  return cachedItems;
};

// Must be cached to avoid infinite loop in useSyncExternalStore
const EMPTY_ITEMS: RecentItem[] = [];
const getServerSnapshot = (): RecentItem[] => EMPTY_ITEMS;

type UseRecentItemsReturn = {
  addItem: (item: Omit<RecentItem, "timestamp" | "type">) => void;
  clearItems: () => void;
  recentItems: RecentItem[];
};

export const useRecentItems = (): UseRecentItemsReturn => {
  const recentItems = useSyncExternalStore(subscribe, getSnapshot, getServerSnapshot);

  const addItem = (item: Omit<RecentItem, "timestamp" | "type">) => {
    addRecentItem(item);
    notifyListeners();
  };

  const clearItems = () => {
    clearRecentItems();
    notifyListeners();
  };

  return {
    addItem,
    clearItems,
    recentItems,
  };
};
