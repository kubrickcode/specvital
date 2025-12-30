"use client";

import { Check, ChevronDown, Circle, CircleDashed, Crosshair, Filter, XCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import type { TestStatus } from "@/lib/api";
import { cn } from "@/lib/utils";

type StatusFilterProps = {
  onChange: (value: TestStatus[]) => void;
  value: TestStatus[];
};

const STATUS_OPTIONS = [
  {
    color: "text-green-600",
    icon: Check,
    key: "active",
  },
  {
    color: "text-purple-500",
    icon: Crosshair,
    key: "focused",
  },
  {
    color: "text-amber-500",
    icon: CircleDashed,
    key: "skipped",
  },
  {
    color: "text-blue-500",
    icon: Circle,
    key: "todo",
  },
  {
    color: "text-red-400",
    icon: XCircle,
    key: "xfail",
  },
] as const;

type StatusKey = (typeof STATUS_OPTIONS)[number]["key"];

const TRANSLATION_KEYS: Record<StatusKey, string> = {
  active: "statusActive",
  focused: "statusFocused",
  skipped: "statusSkipped",
  todo: "statusTodo",
  xfail: "statusXfail",
};

export const StatusFilter = ({ onChange, value }: StatusFilterProps) => {
  const t = useTranslations("analyze.filter");

  const handleToggle = useCallback(
    (statusKey: TestStatus) => {
      const isSelected = value.includes(statusKey);
      if (isSelected) {
        onChange(value.filter((v) => v !== statusKey));
      } else {
        onChange([...value, statusKey]);
      }
    },
    [onChange, value]
  );

  const selectedCount = value.length;

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          className={cn("gap-2 min-h-[44px] sm:min-h-0", selectedCount > 0 && "border-primary/50")}
          size="sm"
          variant="outline"
        >
          <Filter className="h-4 w-4" />
          <span>
            {t("statusFilter")}
            {selectedCount > 0 && ` (${selectedCount})`}
          </span>
          <ChevronDown className="h-3 w-3 opacity-50" />
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-48 p-2">
        <div className="flex flex-col gap-1">
          {STATUS_OPTIONS.map((option) => {
            const Icon = option.icon;
            const isChecked = value.includes(option.key);

            return (
              <label
                className={cn(
                  "flex items-center gap-3 px-2 py-1.5 rounded-md cursor-pointer",
                  "hover:bg-muted/50 transition-colors"
                )}
                key={option.key}
              >
                <Checkbox checked={isChecked} onCheckedChange={() => handleToggle(option.key)} />
                <Icon className={cn("h-4 w-4", option.color)} />
                <span className="text-sm">{t(TRANSLATION_KEYS[option.key])}</span>
              </label>
            );
          })}
        </div>
      </PopoverContent>
    </Popover>
  );
};
