export const STORAGE_KEY = "global-search-recent-items";
const MAX_RECENT_ITEMS = 10;

export type RecentItem = {
  fullName: string;
  id: string;
  name: string;
  owner: string;
  timestamp: number;
  type: "repository";
};

const isNonEmptyString = (value: unknown, maxLength: number): boolean =>
  typeof value === "string" && value.length > 0 && value.length <= maxLength;

const isValidRecentItem = (item: unknown): item is RecentItem => {
  if (typeof item !== "object" || item === null) return false;
  const obj = item as Record<string, unknown>;
  return (
    isNonEmptyString(obj.id, 256) &&
    isNonEmptyString(obj.fullName, 512) &&
    isNonEmptyString(obj.name, 256) &&
    isNonEmptyString(obj.owner, 256) &&
    typeof obj.timestamp === "number" &&
    obj.timestamp > 0 &&
    obj.type === "repository"
  );
};

export const loadRecentItems = (): RecentItem[] => {
  if (typeof window === "undefined") return [];

  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    if (!stored) return [];

    const parsed: unknown = JSON.parse(stored);
    if (!Array.isArray(parsed)) return [];

    return parsed.filter(isValidRecentItem).slice(0, MAX_RECENT_ITEMS);
  } catch {
    return [];
  }
};

export const saveRecentItems = (items: RecentItem[]): void => {
  if (typeof window === "undefined") return;

  try {
    const trimmed = items.slice(0, MAX_RECENT_ITEMS);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(trimmed));
  } catch {
    // Storage full or unavailable - silently ignore
  }
};

export const addRecentItem = (item: Omit<RecentItem, "timestamp" | "type">): RecentItem[] => {
  const existing = loadRecentItems();

  const filtered = existing.filter((i) => i.id !== item.id);

  const newItem: RecentItem = {
    ...item,
    timestamp: Date.now(),
    type: "repository",
  };

  const updated = [newItem, ...filtered].slice(0, MAX_RECENT_ITEMS);
  saveRecentItems(updated);
  return updated;
};

export const clearRecentItems = (): void => {
  if (typeof window === "undefined") return;

  try {
    localStorage.removeItem(STORAGE_KEY);
  } catch {
    // Silently ignore
  }
};
