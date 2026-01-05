import { Check, Circle, CircleDashed, Crosshair, XCircle } from "lucide-react";

import type { TestStatus } from "@/lib/api";
import { cn } from "@/lib/utils";

import type { ConvertedTestItem } from "../../types";

type SpecItemProps = {
  item: ConvertedTestItem;
};

const STATUS_CONFIG: Record<
  TestStatus,
  {
    color: string;
    icon: typeof Check;
    label: string;
  }
> = {
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
};

export const SpecItem = ({ item }: SpecItemProps) => {
  const config = STATUS_CONFIG[item.status];
  const Icon = config.icon;

  return (
    <div
      className={cn(
        "group flex items-start gap-3 px-3 py-2 rounded-md",
        "hover:bg-muted/50 transition-colors"
      )}
    >
      <Icon
        aria-label={config.label}
        className={cn("mt-0.5 h-4 w-4 flex-shrink-0", config.color)}
      />
      <div className="flex-1 min-w-0">
        <p className="text-sm">{item.convertedName}</p>
        <p className="text-xs text-muted-foreground truncate opacity-0 group-hover:opacity-100 transition-opacity">
          {item.originalName}
        </p>
      </div>
      <span className="text-xs text-muted-foreground font-mono flex-shrink-0">L:{item.line}</span>
    </div>
  );
};
