export type HighlightRange = {
  end: number;
  start: number;
};

export const findHighlightRanges = (text: string, query: string): HighlightRange[] => {
  if (!query) return [];

  const ranges: HighlightRange[] = [];
  const lowerText = text.toLowerCase();
  const lowerQuery = query.toLowerCase();

  let startIndex = 0;
  let index = lowerText.indexOf(lowerQuery, startIndex);

  while (index !== -1) {
    ranges.push({ end: index + query.length, start: index });
    startIndex = index + 1;
    index = lowerText.indexOf(lowerQuery, startIndex);
  }

  return ranges;
};
