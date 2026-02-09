import type { Framework } from "@/lib/api";
import { getFrameworkBadgeClasses } from "@/lib/styles";
import { cn } from "@/lib/utils";

type FrameworkBadgeProps = {
  framework: Framework;
};

const capitalize = (str: string): string => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export const FrameworkBadge = ({ framework }: FrameworkBadgeProps) => {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
        getFrameworkBadgeClasses(framework)
      )}
    >
      {capitalize(framework)}
    </span>
  );
};
