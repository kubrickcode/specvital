/**
 * Mock data for repository stats API
 * Endpoint: /api/repositories/stats
 */

import type { RepositoryStatsResponse } from "./types";

export const mockStatsNormal: RepositoryStatsResponse = {
  totalRepositories: 25,
  totalTests: 1543,
};

export const mockStatsEmpty: RepositoryStatsResponse = {
  totalRepositories: 0,
  totalTests: 0,
};

export const mockStatsLarge: RepositoryStatsResponse = {
  totalRepositories: 150,
  totalTests: 25000,
};
