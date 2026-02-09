import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";

import { handleUnauthorizedError, isUnauthorizedError } from "@/lib/api/error-handler";

export const createQueryClient = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: false,
      },
    },
    mutationCache: new MutationCache({
      onError: (error) => {
        if (isUnauthorizedError(error)) {
          handleUnauthorizedError(queryClient);
        }
      },
    }),
    queryCache: new QueryCache({
      onError: (error) => {
        if (isUnauthorizedError(error)) {
          handleUnauthorizedError(queryClient);
        }
      },
    }),
  });

  return queryClient;
};
