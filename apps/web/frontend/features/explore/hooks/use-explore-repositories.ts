import {
  paginatedRepositoriesKeys,
  usePaginatedRepositories,
  type PaginatedRepositoriesOptions,
} from "@/features/dashboard";

type UseExploreRepositoriesOptions = Omit<PaginatedRepositoriesOptions, "view">;

export const useExploreRepositories = (options: UseExploreRepositoriesOptions = {}) => {
  return usePaginatedRepositories({
    ...options,
    view: "community",
  });
};

export const exploreRepositoriesKeys = {
  list: (options: UseExploreRepositoriesOptions) =>
    paginatedRepositoriesKeys.list({ ...options, view: "community" }),
};
