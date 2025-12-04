import { Check, Circle, CircleDashed } from "lucide-react";
import type { TestCase } from "@/lib/api";
import { cn } from "@/lib/utils";

type TestItemProps = {
  test: TestCase;
};

const STATUS_CONFIG = {
  active: {
    icon: Check,
    color: "text-green-600",
    label: "Active test",
  },
  skipped: {
    icon: CircleDashed,
    color: "text-amber-500",
    label: "Skipped test",
  },
  todo: {
    icon: Circle,
    color: "text-blue-500",
    label: "Todo test",
  },
} as const;

export const TestItem = ({ test }: TestItemProps) => {
  const config = STATUS_CONFIG[test.status];
  const Icon = config.icon;

  return (
    <div
      className={cn(
        "flex items-center gap-3 px-3 py-2 rounded-md",
        "hover:bg-muted/50 transition-colors"
      )}
    >
      <Icon
        className={cn("h-4 w-4 flex-shrink-0", config.color)}
        aria-label={config.label}
      />
      <span className="flex-1 text-sm truncate">{test.name}</span>
      <span className="text-xs text-muted-foreground font-mono flex-shrink-0">
        L:{test.line}
      </span>
    </div>
  );
};
