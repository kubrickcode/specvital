"use client";

import { parseAsBoolean, useQueryState } from "nuqs";

const bookmarkFilterParser = parseAsBoolean.withDefault(false);

export const useBookmarkFilter = () => {
  const [bookmarkOnly, setBookmarkOnly] = useQueryState("bookmarked", bookmarkFilterParser);

  return { bookmarkOnly, setBookmarkOnly } as const;
};
