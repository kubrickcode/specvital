import type { Framework } from "@/lib/api";
import { cn } from "@/lib/utils";

type FrameworkBadgeProps = {
  framework: Framework;
};

const FRAMEWORK_STYLES: Record<Framework, string> = {
  go: "bg-cyan-100 text-cyan-800",
  jest: "bg-red-100 text-red-800",
  playwright: "bg-purple-100 text-purple-800",
  vitest: "bg-green-100 text-green-800",
};

const capitalizeFramework = (framework: Framework): string => {
  return framework.charAt(0).toUpperCase() + framework.slice(1);
};

export const FrameworkBadge = ({ framework }: FrameworkBadgeProps) => {
  const colorClasses = FRAMEWORK_STYLES[framework];

  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-0.5 text-xs font-medium",
        colorClasses
      )}
    >
      {capitalizeFramework(framework)}
    </span>
  );
};
