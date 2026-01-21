import type { LucideIcon } from "lucide-react";
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

export type FrameworkCategoryId =
  | "javascript"
  | "python"
  | "go"
  | "java"
  | "kotlin"
  | "csharp"
  | "ruby"
  | "php"
  | "rust"
  | "cpp"
  | "swift";

export type FrameworkCategory = {
  category: FrameworkCategoryId;
  frameworks: string[];
  icon: LucideIcon;
};

export const FRAMEWORK_CATEGORIES: FrameworkCategory[] = [
  {
    category: "javascript",
    frameworks: ["Jest", "Vitest", "Mocha", "Playwright", "Cypress"],
    icon: BracesIcon,
  },
  {
    category: "python",
    frameworks: ["pytest", "unittest"],
    icon: TerminalIcon,
  },
  {
    category: "go",
    frameworks: ["Go Testing"],
    icon: ZapIcon,
  },
  {
    category: "java",
    frameworks: ["JUnit 4", "JUnit 5", "TestNG"],
    icon: CoffeeIcon,
  },
  {
    category: "kotlin",
    frameworks: ["Kotest"],
    icon: DiamondIcon,
  },
  {
    category: "csharp",
    frameworks: ["NUnit", "xUnit", "MSTest"],
    icon: HashIcon,
  },
  {
    category: "ruby",
    frameworks: ["RSpec", "Minitest"],
    icon: GemIcon,
  },
  {
    category: "php",
    frameworks: ["PHPUnit"],
    icon: CodeIcon,
  },
  {
    category: "rust",
    frameworks: ["Cargo Test"],
    icon: CpuIcon,
  },
  {
    category: "cpp",
    frameworks: ["Google Test"],
    icon: DatabaseIcon,
  },
  {
    category: "swift",
    frameworks: ["XCTest"],
    icon: BirdIcon,
  },
];

export const TOTAL_FRAMEWORK_COUNT = FRAMEWORK_CATEGORIES.reduce(
  (count, category) => count + category.frameworks.length,
  0
);

export const HIGHLIGHTED_FRAMEWORKS = ["Jest", "Vitest", "pytest", "JUnit 5", "Go Testing"];
