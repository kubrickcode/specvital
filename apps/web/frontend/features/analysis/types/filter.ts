import type { TestStatus } from "@/lib/api";

export type FilterState = {
  frameworks: string[];
  query: string;
  statuses: TestStatus[];
};
