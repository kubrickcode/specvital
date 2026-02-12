import Fuse, { type IFuseOptions } from "fuse.js";

export type RepositorySearchItem = {
  category: "bookmarks" | "community" | "repositories";
  fullName: string;
  id: string;
  isAnalyzedByMe: boolean;
  isBookmarked: boolean;
  name: string;
  owner: string;
};

export const FUSE_OPTIONS: IFuseOptions<RepositorySearchItem> = {
  findAllMatches: true,
  includeMatches: true,
  includeScore: true,
  keys: [
    { name: "fullName", weight: 2 },
    { name: "name", weight: 1.5 },
    { name: "owner", weight: 1 },
  ],
  minMatchCharLength: 1,
  shouldSort: true,
  threshold: 0.4,
};

export const createRepositoryFuse = (items: RepositorySearchItem[]): Fuse<RepositorySearchItem> =>
  new Fuse(items, FUSE_OPTIONS);
