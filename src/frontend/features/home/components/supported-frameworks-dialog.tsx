"use client";

import {
  BirdIcon,
  BracesIcon,
  CodeIcon,
  CoffeeIcon,
  CpuIcon,
  DatabaseIcon,
  DiamondIcon,
  GemIcon,
  HashIcon,
  TerminalIcon,
  ZapIcon,
} from "lucide-react";
import { useTranslations } from "next-intl";
import type { ReactNode } from "react";

import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";

type SupportedFrameworksDialogProps = {
  onOpenChange: (open: boolean) => void;
  open: boolean;
};

type FrameworkCategory = {
  category: string;
  frameworks: string[];
  icon: ReactNode;
};

const FRAMEWORK_CATEGORIES: FrameworkCategory[] = [
  {
    category: "javascript",
    frameworks: ["Jest", "Vitest", "Mocha", "Playwright", "Cypress"],
    icon: <BracesIcon className="size-4" />,
  },
  {
    category: "python",
    frameworks: ["pytest", "unittest"],
    icon: <TerminalIcon className="size-4" />,
  },
  {
    category: "go",
    frameworks: ["Go Testing"],
    icon: <ZapIcon className="size-4" />,
  },
  {
    category: "java",
    frameworks: ["JUnit 4", "JUnit 5", "TestNG"],
    icon: <CoffeeIcon className="size-4" />,
  },
  {
    category: "kotlin",
    frameworks: ["Kotest"],
    icon: <DiamondIcon className="size-4" />,
  },
  {
    category: "csharp",
    frameworks: ["NUnit", "xUnit", "MSTest"],
    icon: <HashIcon className="size-4" />,
  },
  {
    category: "ruby",
    frameworks: ["RSpec", "Minitest"],
    icon: <GemIcon className="size-4" />,
  },
  {
    category: "php",
    frameworks: ["PHPUnit"],
    icon: <CodeIcon className="size-4" />,
  },
  {
    category: "rust",
    frameworks: ["Cargo Test"],
    icon: <CpuIcon className="size-4" />,
  },
  {
    category: "cpp",
    frameworks: ["Google Test"],
    icon: <DatabaseIcon className="size-4" />,
  },
  {
    category: "swift",
    frameworks: ["XCTest"],
    icon: <BirdIcon className="size-4" />,
  },
];

export const SupportedFrameworksDialog = ({
  onOpenChange,
  open,
}: SupportedFrameworksDialogProps) => {
  const t = useTranslations("home.frameworks");

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <div className="max-h-[60vh] space-y-4 overflow-y-auto pr-2">
          {FRAMEWORK_CATEGORIES.map((category) => (
            <div key={category.category}>
              <div className="mb-2 flex items-center gap-2">
                <span className="text-muted-foreground">{category.icon}</span>
                <h3 className="text-sm font-medium">{t(category.category)}</h3>
              </div>
              <div className="flex flex-wrap gap-2">
                {category.frameworks.map((framework) => (
                  <Badge className="bg-black/5 text-foreground dark:bg-white/10" key={framework}>
                    {framework}
                  </Badge>
                ))}
              </div>
            </div>
          ))}
        </div>
      </DialogContent>
    </Dialog>
  );
};
