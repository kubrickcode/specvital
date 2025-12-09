import { QueryClient } from "@tanstack/react-query";

export const createQueryClient = () =>
  new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: false,
        staleTime: 1000 * 60, // 1 minute
      },
    },
  });
