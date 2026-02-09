/**
 * Task ID generation utilities for spec generation background tasks.
 * Centralizes task ID pattern to prevent inconsistencies across components.
 */

export const getSpecGenerationTaskId = (analysisId: string): string =>
  `spec-generation-${analysisId}`;
