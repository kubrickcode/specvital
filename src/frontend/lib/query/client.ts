import { MutationCache, QueryCache, QueryClient } from "@tanstack/react-query";

import { handleUnauthorizedError, isAuthQuery, isUnauthorizedError } from "@/lib/api/error-handler";

export const createQueryClient = () => {
  const queryClient = new QueryClient({
    defaultOptions: {
      queries: {
        refetchOnWindowFocus: false,
        retry: false,
        staleTime: 1000 * 60, // 1 minute
      },
    },
    mutationCache: new MutationCache({
      onError: (error, _variables, _context, mutation) => {
        if (
          isUnauthorizedError(error) &&
          mutation.options.mutationKey &&
          isAuthQuery(mutation.options.mutationKey)
        ) {
          handleUnauthorizedError(queryClient);
        }
      },
    }),
    queryCache: new QueryCache({
      onError: (error, query) => {
        if (isUnauthorizedError(error) && isAuthQuery(query.queryKey)) {
          handleUnauthorizedError(queryClient);
        }
      },
    }),
  });

  return queryClient;
};
