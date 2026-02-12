import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { AvailableLanguageInfo, RepoVersionInfo, SpecDocument } from "../types";

import { ExecutiveSummary } from "./executive-summary";

vi.mock("./cache-stats-indicator", () => ({
  CacheStatsIndicator: () => <span data-testid="cache-stats">cache</span>,
}));

vi.mock("./spec-export-button", () => ({
  SpecExportButton: () => <button type="button">Export</button>,
}));

vi.mock("./status-legend", () => ({
  StatusLegend: () => <span data-testid="status-legend">legend</span>,
}));

// Radix DropdownMenu checks `event instanceof PointerEvent` which jsdom lacks
class MockPointerEvent extends Event {
  button: number;
  ctrlKey: boolean;
  pointerType: string;
  constructor(type: string, props: PointerEventInit & EventInit = {}) {
    super(type, props);
    this.button = props.button ?? 0;
    this.ctrlKey = props.ctrlKey ?? false;
    this.pointerType = props.pointerType ?? "mouse";
  }
}
window.PointerEvent = MockPointerEvent as never;

Object.defineProperty(window, "matchMedia", {
  value: vi.fn().mockImplementation((query: string) => ({
    addEventListener: vi.fn(),
    addListener: vi.fn(),
    dispatchEvent: vi.fn(),
    matches: false,
    media: query,
    onchange: null,
    removeEventListener: vi.fn(),
    removeListener: vi.fn(),
  })),
  writable: true,
});

const openDropdown = (trigger: HTMLElement) => {
  fireEvent.pointerDown(trigger, { button: 0, ctrlKey: false, pointerType: "mouse" });
};

const messages = {
  specView: {
    executiveSummary: {
      availableLanguages: "Available Languages",
      dateTooltip: "Creation date",
      generateNew: "Generate New",
      languageTooltip: "Document language",
      latestLabel: "Latest",
      noSummary: "No executive summary available.",
      regenerateAriaLabel: "Regenerate spec",
      regenerateTooltip: "Regenerate this spec document",
      switchLanguageTooltip: "Switch language",
      switchVersionTooltip: "Switch version",
      title: "Executive Summary",
      versionHistory: "Version History",
      versionLabel: "v{version}",
    },
    stats: {
      behaviors: "{count} behaviors",
      domains: "{count} domains",
      features: "{count} features",
    },
  },
};

const createDocument = (overrides: Partial<SpecDocument> = {}): SpecDocument => ({
  analysisId: "analysis-1",
  availableLanguages: [],
  createdAt: "2024-06-01T12:00:00Z",
  domains: [
    {
      classificationConfidence: 0.9,
      description: "Auth domain",
      features: [
        {
          behaviors: [
            {
              convertedDescription: "Login works",
              id: "b1",
              originalName: "should login",
              sortOrder: 1,
              testCaseId: "tc1",
            },
          ],
          description: "Login feature",
          id: "f1",
          name: "Login",
          sortOrder: 1,
        },
      ],
      id: "d1",
      name: "Authentication",
      sortOrder: 1,
    },
  ],
  executiveSummary: "This is an executive summary.",
  id: "doc-1",
  language: "English",
  version: 1,
  ...overrides,
});

const createVersions = (count: number): RepoVersionInfo[] =>
  Array.from({ length: count }, (_, i) => ({
    analysisId: `analysis-${i + 1}`,
    commitSha: `commit${String(i + 1).padStart(7, "0")}`,
    createdAt: `2024-0${i + 1}-01T12:00:00Z`,
    id: `doc-${i + 1}`,
    language: "English" as const,
    modelId: "gemini-2.0-flash",
    version: i + 1,
  }));

const renderExecutiveSummary = (
  props: Partial<React.ComponentProps<typeof ExecutiveSummary>> = {}
) => {
  const defaultProps: React.ComponentProps<typeof ExecutiveSummary> = {
    document: createDocument(),
    ...props,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <ExecutiveSummary {...defaultProps} />
    </NextIntlClientProvider>
  );
};

describe("ExecutiveSummary", () => {
  describe("Version History", () => {
    it("renders dropdown when multiple versions exist", () => {
      const versions = createVersions(3);
      renderExecutiveSummary({
        commitSha: "commit0000001",
        onVersionSwitch: vi.fn(),
        versions,
      });

      // The version trigger button should have data-slot="tooltip-trigger" and contain chevron
      const buttons = screen.getAllByRole("button");
      const versionTrigger = buttons.find(
        (b) => b.textContent?.includes("Jun 1") && b.getAttribute("aria-haspopup") === "menu"
      );
      expect(versionTrigger).toBeDefined();
    });

    it("renders static button without dropdown when only one version exists", () => {
      const versions = createVersions(1);
      renderExecutiveSummary({
        onVersionSwitch: vi.fn(),
        versions,
      });

      // With a single version, no dropdown trigger, just a static button
      const buttons = screen.getAllByRole("button");
      const versionButton = buttons.find((b) => b.textContent?.includes("Jun 1"));
      expect(versionButton).toBeDefined();
      expect(versionButton).not.toHaveAttribute("aria-haspopup", "menu");
    });

    it("shows all versions in dropdown list when opened", () => {
      const versions = createVersions(3);
      renderExecutiveSummary({
        commitSha: "commit0000001",
        latestDocumentId: "doc-1",
        onVersionSwitch: vi.fn(),
        versions,
      });

      const trigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("Jun 1")
        )!;
      openDropdown(trigger);

      // Version History section header
      expect(screen.getByText("Version History")).toBeInTheDocument();

      // 3 versions produce 3 commit groups in the dropdown, plus 1 in the trigger button = 4 total
      const commitCodeElements = screen.getAllByText("commit0");
      expect(commitCodeElements.length).toBeGreaterThanOrEqual(3);

      // All 3 version items should be rendered as menuitems
      expect(screen.getAllByRole("menuitem")).toHaveLength(3);
    });

    it("calls onVersionSwitch when a different version is selected", () => {
      const versions = createVersions(2);
      const onVersionSwitch = vi.fn();
      renderExecutiveSummary({
        commitSha: "commit0000001",
        document: createDocument({ id: "doc-1" }),
        latestDocumentId: "doc-1",
        onVersionSwitch,
        versions,
      });

      const trigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("Jun 1")
        )!;
      openDropdown(trigger);

      const menuItems = screen.getAllByRole("menuitem");
      // Click the version that is not current
      const otherItem = menuItems.find((item) => !item.classList.contains("bg-muted"));
      expect(otherItem).toBeDefined();
      fireEvent.click(otherItem!);

      expect(onVersionSwitch).toHaveBeenCalledWith("doc-2");
    });

    it("marks current version with check icon and active styling", () => {
      const versions = createVersions(2);
      renderExecutiveSummary({
        commitSha: "commit0000001",
        document: createDocument({ id: "doc-1" }),
        latestDocumentId: "doc-1",
        onVersionSwitch: vi.fn(),
        versions,
      });

      const trigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("Jun 1")
        )!;
      openDropdown(trigger);

      const menuItems = screen.getAllByRole("menuitem");
      const currentItem = menuItems.find((item) => item.classList.contains("bg-muted"));
      expect(currentItem).toBeDefined();
      // Current version should have a check icon SVG
      expect(currentItem!.querySelector("svg")).toBeInTheDocument();
    });
  });

  describe("Regeneration", () => {
    it("renders regenerate button when onRegenerate is provided", () => {
      renderExecutiveSummary({
        onRegenerate: vi.fn(),
      });

      expect(screen.getByRole("button", { name: "Regenerate spec" })).toBeInTheDocument();
    });

    it("calls onRegenerate when regenerate button is clicked", () => {
      const onRegenerate = vi.fn();
      renderExecutiveSummary({ onRegenerate });

      fireEvent.click(screen.getByRole("button", { name: "Regenerate spec" }));

      expect(onRegenerate).toHaveBeenCalledOnce();
    });
  });

  describe("Language Dropdown extended", () => {
    const availableLanguages: AvailableLanguageInfo[] = [
      {
        createdAt: "2024-05-01T10:00:00Z",
        hasPreviousSpec: false,
        language: "English",
        latestVersion: 2,
      },
      {
        createdAt: "2024-05-15T10:00:00Z",
        hasPreviousSpec: true,
        language: "Korean",
        latestVersion: 1,
      },
    ];

    it("shows Available Languages and Generate New two-tier structure", () => {
      renderExecutiveSummary({
        document: createDocument({ availableLanguages, language: "English" }),
        onGenerateNewLanguage: vi.fn(),
        onLanguageSwitch: vi.fn(),
      });

      const langTrigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("English")
        )!;
      openDropdown(langTrigger);

      expect(screen.getByText("Available Languages")).toBeInTheDocument();
      expect(screen.getByText("Generate New")).toBeInTheDocument();
    });

    it("shows check icon and version info for current language", () => {
      renderExecutiveSummary({
        document: createDocument({ availableLanguages, language: "English" }),
        onGenerateNewLanguage: vi.fn(),
        onLanguageSwitch: vi.fn(),
      });

      const langTrigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("English")
        )!;
      openDropdown(langTrigger);

      // Current language (English) should have bg-muted styling and check icon
      const menuItems = screen.getAllByRole("menuitem");
      const englishItem = menuItems.find(
        (item) => item.textContent?.includes("English") && item.classList.contains("bg-muted")
      );
      expect(englishItem).toBeDefined();
      expect(englishItem!.querySelector("svg")).toBeInTheDocument();

      // Should show version info (v2 for English)
      expect(screen.getByText("v2")).toBeInTheDocument();
    });

    it("shows version and date info for available languages", () => {
      renderExecutiveSummary({
        document: createDocument({ availableLanguages, language: "English" }),
        onGenerateNewLanguage: vi.fn(),
        onLanguageSwitch: vi.fn(),
      });

      const langTrigger = screen
        .getAllByRole("button")
        .find(
          (b) => b.getAttribute("aria-haspopup") === "menu" && b.textContent?.includes("English")
        )!;
      openDropdown(langTrigger);

      // Version info for both languages
      expect(screen.getByText("v2")).toBeInTheDocument();
      expect(screen.getByText("v1")).toBeInTheDocument();

      // Date info - use getAllByText since "May 1" might appear in trigger area too
      const may1Elements = screen.getAllByText((_content, element) => {
        return element?.textContent === "May 1" || false;
      });
      expect(may1Elements.length).toBeGreaterThanOrEqual(1);
      expect(screen.getByText("May 15")).toBeInTheDocument();
    });
  });
});
