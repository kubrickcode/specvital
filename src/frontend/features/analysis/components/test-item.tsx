import { Check, Circle, CircleDashed, Crosshair, XCircle } from "lucide-react";

import type { TestCase } from "@/lib/api";
import { cn } from "@/lib/utils";

type TestItemProps = {
  test: TestCase;
};

const STATUS_CONFIG = {
  active: {
    color: "text-green-600",
    icon: Check,
    label: "Active test",
  },
  focused: {
    color: "text-purple-500",
    icon: Crosshair,
    label: "Focused test",
  },
  skipped: {
    color: "text-amber-500",
    icon: CircleDashed,
    label: "Skipped test",
  },
  todo: {
    color: "text-blue-500",
    icon: Circle,
    label: "Todo test",
  },
  xfail: {
    color: "text-red-400",
    icon: XCircle,
    label: "Expected failure",
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
      <Icon aria-label={config.label} className={cn("h-4 w-4 flex-shrink-0", config.color)} />
      <span className="flex-1 text-sm truncate">{test.name}</span>
      <span className="text-xs text-muted-foreground font-mono flex-shrink-0">L:{test.line}</span>
    </div>
  );
};
