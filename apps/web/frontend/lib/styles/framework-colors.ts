type FrameworkColorConfig = {
  badge: string;
  solid: string;
};

const FRAMEWORK_COLORS: Record<string, FrameworkColorConfig> = {
  // JavaScript/TypeScript
  cypress: { badge: "bg-emerald-100 text-emerald-800", solid: "#1b9e77" },
  jest: { badge: "bg-red-100 text-red-800", solid: "#c21325" },
  mocha: { badge: "bg-amber-100 text-amber-800", solid: "#8d6748" },
  playwright: { badge: "bg-green-100 text-green-700", solid: "#2ead33" },
  vitest: { badge: "bg-green-100 text-green-800", solid: "#6da13f" },

  // Python
  pytest: { badge: "bg-blue-100 text-blue-800", solid: "#0a9edc" },
  unittest: { badge: "bg-sky-100 text-sky-800", solid: "#3572a5" },

  // Go
  go: { badge: "bg-cyan-100 text-cyan-800", solid: "#00add8" },
  "go-testing": { badge: "bg-cyan-100 text-cyan-800", solid: "#00add8" },

  // Java
  junit: { badge: "bg-orange-100 text-orange-800", solid: "#25a162" },
  junit4: { badge: "bg-orange-100 text-orange-800", solid: "#25a162" },
  junit5: { badge: "bg-orange-100 text-orange-800", solid: "#25a162" },
  testng: { badge: "bg-yellow-100 text-yellow-800", solid: "#cd6600" },

  // Kotlin
  kotest: { badge: "bg-purple-100 text-purple-800", solid: "#7f52ff" },

  // C#
  mstest: { badge: "bg-violet-100 text-violet-800", solid: "#512bd4" },
  nunit: { badge: "bg-violet-100 text-violet-800", solid: "#512bd4" },
  xunit: { badge: "bg-indigo-100 text-indigo-800", solid: "#68217a" },

  // Ruby
  minitest: { badge: "bg-red-100 text-red-700", solid: "#cc342d" },
  rspec: { badge: "bg-rose-100 text-rose-800", solid: "#cc342d" },

  // PHP
  phpunit: { badge: "bg-sky-100 text-sky-800", solid: "#4f5b93" },

  // Rust
  cargo: { badge: "bg-orange-100 text-orange-700", solid: "#f74c00" },
  "cargo-test": { badge: "bg-orange-100 text-orange-700", solid: "#f74c00" },

  // C++
  "google-test": { badge: "bg-blue-100 text-blue-800", solid: "#4285f4" },
  googletest: { badge: "bg-blue-100 text-blue-800", solid: "#4285f4" },

  // Swift
  xctest: { badge: "bg-orange-100 text-orange-800", solid: "#f05138" },
};

const FALLBACK_BADGE_PALETTE = [
  "bg-red-100 text-red-800",
  "bg-orange-100 text-orange-800",
  "bg-amber-100 text-amber-800",
  "bg-yellow-100 text-yellow-800",
  "bg-lime-100 text-lime-800",
  "bg-green-100 text-green-800",
  "bg-emerald-100 text-emerald-800",
  "bg-teal-100 text-teal-800",
  "bg-cyan-100 text-cyan-800",
  "bg-sky-100 text-sky-800",
  "bg-blue-100 text-blue-800",
  "bg-indigo-100 text-indigo-800",
  "bg-violet-100 text-violet-800",
  "bg-purple-100 text-purple-800",
  "bg-fuchsia-100 text-fuchsia-800",
  "bg-pink-100 text-pink-800",
  "bg-rose-100 text-rose-800",
  "bg-slate-100 text-slate-800",
  "bg-zinc-100 text-zinc-800",
  "bg-stone-100 text-stone-800",
];

const FALLBACK_SOLID_PALETTE = [
  "#ef4444",
  "#f97316",
  "#f59e0b",
  "#eab308",
  "#84cc16",
  "#22c55e",
  "#10b981",
  "#14b8a6",
  "#06b6d4",
  "#0ea5e9",
  "#3b82f6",
  "#6366f1",
  "#8b5cf6",
  "#a855f7",
  "#d946ef",
  "#ec4899",
  "#f43f5e",
  "#64748b",
  "#71717a",
  "#78716c",
];

const hashString = (str: string): number => {
  let hash = 5381;
  for (let i = 0; i < str.length; i++) {
    hash = (hash * 33) ^ str.charCodeAt(i);
  }
  return hash >>> 0;
};

export const getFrameworkBadgeClasses = (framework: string): string => {
  const normalized = framework.toLowerCase();
  const config = FRAMEWORK_COLORS[normalized];
  if (config) {
    return config.badge;
  }
  const index = hashString(normalized) % FALLBACK_BADGE_PALETTE.length;
  return FALLBACK_BADGE_PALETTE[index];
};

export const getFrameworkSolidColor = (framework: string): string => {
  const normalized = framework.toLowerCase();
  const config = FRAMEWORK_COLORS[normalized];
  if (config) {
    return config.solid;
  }
  const index = hashString(normalized) % FALLBACK_SOLID_PALETTE.length;
  return FALLBACK_SOLID_PALETTE[index];
};
