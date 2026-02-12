import type { HighlightRange } from "../utils/highlight";

type HighlightedTextProps = {
  ranges: HighlightRange[];
  text: string;
};

export const HighlightedText = ({ ranges, text }: HighlightedTextProps) => {
  if (ranges.length === 0) {
    return <>{text}</>;
  }

  const parts: React.ReactNode[] = [];
  let lastIndex = 0;

  for (const range of ranges) {
    if (range.start > lastIndex) {
      parts.push(text.slice(lastIndex, range.start));
    }
    parts.push(
      <mark className="bg-yellow-200 dark:bg-yellow-800 rounded-sm px-0.5" key={range.start}>
        {text.slice(range.start, range.end)}
      </mark>
    );
    lastIndex = range.end;
  }

  if (lastIndex < text.length) {
    parts.push(text.slice(lastIndex));
  }

  return <>{parts}</>;
};
