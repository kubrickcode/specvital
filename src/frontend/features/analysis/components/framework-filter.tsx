"use client";

import { Filter } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Checkbox } from "@/components/ui/checkbox";
import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { getFrameworkSolidColor } from "@/lib/styles";
import { cn } from "@/lib/utils";

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

  const handleToggle = (framework: string) => {
    const isSelected = value.includes(framework);
    if (isSelected) {
      onChange(value.filter((f) => f !== framework));
    } else {
      onChange([...value, framework]);
    }
  };

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

            return (
              <label
                className="flex cursor-pointer items-center gap-3 rounded-md px-2 py-1.5 hover:bg-accent transition-colors"
                key={framework}
              >
                <Checkbox checked={isChecked} onCheckedChange={() => handleToggle(framework)} />
                <span
                  className="h-2.5 w-2.5 rounded-full shrink-0"
                  style={{ backgroundColor: getFrameworkSolidColor(framework) }}
                />
                <span className="text-sm">{capitalize(framework)}</span>
              </label>
            );
          })}
        </div>
      </PopoverContent>
    </Popover>
  );
};
