"use client";

import { QueryClientProvider } from "@tanstack/react-query";
import { useState } from "react";

import { createQueryClient } from "./client";

type QueryProviderProps = {
  children: React.ReactNode;
};

export const QueryProvider = ({ children }: QueryProviderProps) => {
  const [queryClient] = useState(() => createQueryClient());

  return <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>;
};
