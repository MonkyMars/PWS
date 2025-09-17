import { QueryClient } from "@tanstack/react-query";

/**
 * Creates and configures a new QueryClient instance for TanStack Query
 * with optimized settings for the ELO application
 */
export function createQueryClient(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        // Cache data for 5 minutes by default
        staleTime: 1000 * 60 * 5,
        // Keep data in cache for 10 minutes
        gcTime: 1000 * 60 * 10,
        // Retry failed requests twice
        retry: 2,
        // Don't refetch on window focus for better UX
        refetchOnWindowFocus: false,
      },
      mutations: {
        // Retry failed mutations once
        retry: 1,
      },
    },
  });
}

let queryClient: QueryClient | undefined;

/**
 * Gets the global QueryClient instance, creating one if it doesn't exist
 */
export function getQueryClient(): QueryClient {
  if (!queryClient) {
    queryClient = createQueryClient();
  }
  return queryClient;
}
