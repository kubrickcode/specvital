"use client";

import type { FuseResultMatch } from "fuse.js";
import { Bookmark, FolderGit2, Star } from "lucide-react";

import { CommandItem } from "@/components/ui/command";

import type { ScoredResult } from "../lib/score-results";

type RepositorySearchItemProps = {
  onSelect: () => void;
  result: ScoredResult;
};

export const RepositorySearchItem = ({ onSelect, result }: RepositorySearchItemProps) => {
  const { item, matches } = result;

  const renderHighlightedText = (text: string, fieldName: string) => {
    const fieldMatches = matches?.filter((m: FuseResultMatch) => m.key === fieldName);
    if (!fieldMatches?.length) {
      return text;
    }

    const indices = fieldMatches.flatMap((m: FuseResultMatch) => m.indices);
    if (!indices.length) {
      return text;
    }

    const sortedIndices = [...indices].sort((a, b) => a[0] - b[0]);
    const mergedIndices: Array<[number, number]> = [];
    for (const [start, end] of sortedIndices) {
      const last = mergedIndices[mergedIndices.length - 1];
      if (last && start <= last[1] + 1) {
        last[1] = Math.max(last[1], end);
      } else {
        mergedIndices.push([start, end]);
      }
    }

    const parts: Array<{ highlighted: boolean; text: string }> = [];
    let lastIndex = 0;

    for (const [start, end] of mergedIndices) {
      if (start > lastIndex) {
        parts.push({ highlighted: false, text: text.slice(lastIndex, start) });
      }
      parts.push({ highlighted: true, text: text.slice(start, end + 1) });
      lastIndex = end + 1;
    }

    if (lastIndex < text.length) {
      parts.push({ highlighted: false, text: text.slice(lastIndex) });
    }

    return parts.map((part, index) =>
      part.highlighted ? (
        <mark className="bg-primary/20 text-foreground rounded-sm px-0.5" key={index}>
          {part.text}
        </mark>
      ) : (
        <span key={index}>{part.text}</span>
      )
    );
  };

  const getCategoryIcon = () => {
    if (item.isAnalyzedByMe) {
      return <Star className="size-3 text-yellow-500 shrink-0" />;
    }
    if (item.isBookmarked) {
      return <Bookmark className="size-3 text-blue-500 shrink-0" />;
    }
    return null;
  };

  return (
    <CommandItem className="flex items-center gap-2" onSelect={onSelect} value={`repo-${item.id}`}>
      <FolderGit2 className="size-4 text-muted-foreground shrink-0" />
      <span className="flex-1 truncate">{renderHighlightedText(item.fullName, "fullName")}</span>
      {getCategoryIcon()}
    </CommandItem>
  );
};
