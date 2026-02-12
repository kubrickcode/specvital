"use client";

import { Clock, FolderGit2 } from "lucide-react";

import { CommandItem } from "@/components/ui/command";

import type { RecentItem } from "../lib/recent-items";

type RecentSearchItemProps = {
  item: RecentItem;
  onSelect: () => void;
};

export const RecentSearchItem = ({ item, onSelect }: RecentSearchItemProps) => {
  return (
    <CommandItem
      className="flex items-center gap-2"
      onSelect={onSelect}
      value={`recent-${item.id}`}
    >
      <FolderGit2 className="size-4 text-muted-foreground shrink-0" />
      <span className="flex-1 truncate">{item.fullName}</span>
      <Clock aria-hidden="true" className="size-3 text-muted-foreground shrink-0" />
    </CommandItem>
  );
};
