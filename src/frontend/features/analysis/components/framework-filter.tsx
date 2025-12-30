"use client";

import { Filter } from "lucide-react";
import { useTranslations } from "next-intl";
import { useCallback } from "react";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { cn } from "@/lib/utils";

const COLOR_PALETTE = [
  "bg-red-500",
  "bg-orange-500",
  "bg-amber-500",
  "bg-yellow-500",
  "bg-lime-500",
  "bg-green-500",
  "bg-emerald-500",
  "bg-teal-500",
  "bg-cyan-500",
  "bg-sky-500",
  "bg-blue-500",
  "bg-indigo-500",
  "bg-violet-500",
  "bg-purple-500",
  "bg-fuchsia-500",
  "bg-pink-500",
  "bg-rose-500",
  "bg-slate-500",
  "bg-zinc-500",
  "bg-stone-500",
];

const hashToIndex = (str: string): number => {
  let hash = 5381;
  for (let i = 0; i < str.length; i++) {
    hash = (hash * 33) ^ str.charCodeAt(i);
  }
  return (hash >>> 0) % COLOR_PALETTE.length;
};

const getFrameworkColor = (framework: string): string => {
  const colorIndex = hashToIndex(framework.toLowerCase());
  return COLOR_PALETTE[colorIndex];
};

const capitalize = (str: string): string => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

type FrameworkFilterProps = {
  availableFrameworks: string[];
  onChange: (value: string[]) => void;
  value: string[];
};

export const FrameworkFilter = ({ availableFrameworks, onChange, value }: FrameworkFilterProps) => {
  const t = useTranslations("analyze.filter");

  const handleToggle = useCallback(
    (framework: string) => {
      const isSelected = value.includes(framework);
      if (isSelected) {
        onChange(value.filter((f) => f !== framework));
      } else {
        onChange([...value, framework]);
      }
    },
    [onChange, value]
  );

  if (availableFrameworks.length <= 1) {
    return null;
  }

  const selectedCount = value.length;
  const buttonLabel =
    selectedCount > 0 ? `${t("frameworkFilter")} (${selectedCount})` : t("frameworkFilter");

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          className={cn(selectedCount > 0 && "border-primary/50 bg-primary/5")}
          size="default"
          variant="outline"
        >
          <Filter className="h-4 w-4" />
          {buttonLabel}
        </Button>
      </PopoverTrigger>
      <PopoverContent align="start" className="w-56 p-2">
        <div className="flex flex-col gap-1">
          {availableFrameworks.map((framework) => {
            const isChecked = value.includes(framework);
            const colorClass = getFrameworkColor(framework);

            return (
              <label
                className="flex cursor-pointer items-center gap-3 rounded-md px-2 py-1.5 hover:bg-accent transition-colors"
                key={framework}
              >
                <Checkbox checked={isChecked} onCheckedChange={() => handleToggle(framework)} />
                <span className={cn("h-2.5 w-2.5 rounded-full shrink-0", colorClass)} />
                <span className="text-sm">{capitalize(framework)}</span>
              </label>
            );
          })}
        </div>
      </PopoverContent>
    </Popover>
  );
};
