import { cleanup, fireEvent, render, screen } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { QuotaConfirmDialog } from "./quota-confirm-dialog";

const mockClose = vi.fn();
const mockConfirm = vi.fn();
const mockOnOpenChange = vi.fn();
const mockSetSelectedLanguage = vi.fn();
const mockSetForceRegenerate = vi.fn();

const defaultDialogState = {
  analysisId: "analysis-123",
  close: mockClose,
  confirm: mockConfirm,
  estimatedCost: 10,
  forceRegenerate: false,
  isOpen: true,
  isRegenerate: false,
  isSameCommit: false,
  onOpenChange: mockOnOpenChange,
  regeneratingLanguage: null as string | null,
  selectedLanguage: "English",
  setForceRegenerate: mockSetForceRegenerate,
  setSelectedLanguage: mockSetSelectedLanguage,
  usage: {
    analysis: { limit: 100, percentage: 30, reserved: 0, used: 30 },
    plan: {
      analysisMonthlyLimit: 100,
      retentionDays: 90,
      specviewMonthlyLimit: 1000,
      tier: "pro" as const,
    },
    resetAt: "2026-03-01T00:00:00Z",
    specview: { limit: 1000, percentage: 30, reserved: 0, used: 300 },
  },
};

let dialogState = { ...defaultDialogState };

vi.mock("../hooks/use-quota-confirm-dialog", () => ({
  useQuotaConfirmDialog: () => dialogState,
}));

vi.mock("@tanstack/react-query", () => ({
  useQuery: () => ({ data: undefined, isError: false }),
}));

vi.mock("next-intl", () => ({
  useTranslations: (namespace: string) => {
    const allMessages: Record<string, Record<string, string>> = {
      "specView.generate": {
        analysisMode: "Analysis Mode",
        cacheCheckFailed: "Could not check cache availability.",
        fresh: "Fresh analysis",
        freshWarning: "Quota will be consumed for all tests",
        recommended: "Recommended",
        sameCommitCacheDisabled: "Cache is unavailable for the same commit",
        withCache: "Use cache",
        withCacheBenefit: "Reuses previous analysis to save quota",
      },
      "specView.quota": {
        processing: "processing",
      },
      "specView.quotaConfirm": {
        afterUsageShort: "After ~{after} / {limit}",
        cancel: "Cancel",
        current: "Current",
        dangerMessage: "You're almost at your monthly limit.",
        description: "This will use your AI Spec Doc generation quota.",
        estimatedCostNote: "Actual usage may differ from the estimate",
        exceededMessage: "You've reached your monthly quota limit.",
        generate: "Generate Document",
        limit: "Limit",
        regenerate: "Regenerate Document",
        regenerateDescription: "A new version will be generated.",
        regenerateTitle: "Regenerate Specification",
        specviewUsage: "Usage This Month",
        title: "Generate Specification",
        unit: "tests",
        unlimited: "Unlimited",
        used: "used",
        viewAccount: "View Account →",
        warningMessage: "You're running low on your monthly quota.",
        wouldExceedMessage: "This generation would exceed your monthly quota limit.",
      },
    };
    const messages = allMessages[namespace] || {};
    return (key: string, params?: Record<string, string | number>) => {
      let value = messages[key] || key;
      if (params) {
        for (const [paramKey, paramValue] of Object.entries(params)) {
          value = value.replace(`{${paramKey}}`, String(paramValue));
        }
      }
      return value;
    };
  },
}));

vi.mock("./language-combobox", () => ({
  LanguageCombobox: ({ value }: { onValueChange: (v: string) => void; value: string }) => (
    <div data-testid="language-combobox">{value}</div>
  ),
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({
    children,
    href,
    onClick,
    ...props
  }: {
    children: React.ReactNode;
    className?: string;
    href: string;
    onClick?: () => void;
  }) => (
    <a href={href} onClick={onClick} {...props}>
      {children}
    </a>
  ),
}));

describe("QuotaConfirmDialog", () => {
  beforeEach(() => {
    dialogState = { ...defaultDialogState };
    vi.clearAllMocks();
  });

  afterEach(() => {
    cleanup();
  });

  describe("dialog open/close", () => {
    it("renders dialog when isOpen is true", () => {
      render(<QuotaConfirmDialog />);

      expect(screen.getByText("Generate Specification")).toBeInTheDocument();
    });

    it("does not render dialog content when isOpen is false", () => {
      dialogState = { ...defaultDialogState, isOpen: false };

      render(<QuotaConfirmDialog />);

      expect(screen.queryByText("Generate Specification")).not.toBeInTheDocument();
    });

    it("calls close when cancel button is clicked", () => {
      render(<QuotaConfirmDialog />);

      fireEvent.click(screen.getByRole("button", { name: "Cancel" }));

      expect(mockClose).toHaveBeenCalledOnce();
    });
  });

  describe("generate vs regenerate mode", () => {
    it("shows generate title for initial generation", () => {
      render(<QuotaConfirmDialog />);

      expect(screen.getByText("Generate Specification")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Generate Document" })).toBeInTheDocument();
    });

    it("shows regenerate title for regeneration", () => {
      dialogState = { ...defaultDialogState, isRegenerate: true };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText("Regenerate Specification")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Regenerate Document" })).toBeInTheDocument();
    });
  });

  describe("usage display", () => {
    it("displays current usage count", () => {
      render(<QuotaConfirmDialog />);

      expect(screen.getByText("Usage This Month")).toBeInTheDocument();
      expect(screen.getByText(/Current/)).toBeInTheDocument();
      expect(screen.getByText(/300/)).toBeInTheDocument();
    });

    it("displays limit value", () => {
      render(<QuotaConfirmDialog />);

      expect(screen.getByText(/Limit.*1,000/)).toBeInTheDocument();
    });

    it("shows unlimited text for unlimited plans", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: null, percentage: null, reserved: 0, used: 150 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText(/150/)).toBeInTheDocument();
      expect(screen.getByText(/Unlimited/)).toBeInTheDocument();
    });

    it("shows reserved processing count when reserved is non-zero", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 30, reserved: 5, used: 300 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText(/5.*processing/)).toBeInTheDocument();
    });
  });

  describe("quota exceeded state", () => {
    it("disables generate button when quota is exceeded", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 100, reserved: 0, used: 1000 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByRole("button", { name: "Generate Document" })).toBeDisabled();
    });

    it("shows exceeded message when quota is at 100%", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 100, reserved: 0, used: 1000 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText("You've reached your monthly quota limit.")).toBeInTheDocument();
    });

    it("shows view account link when quota is exceeded", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 100, reserved: 0, used: 1000 },
        },
      };

      render(<QuotaConfirmDialog />);

      const viewAccountLink = screen.getByText("View Account →");
      expect(viewAccountLink).toBeInTheDocument();
      expect(viewAccountLink.closest("a")).toHaveAttribute("href", "/account");
    });

    it("calls close when view account link is clicked", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 100, reserved: 0, used: 1000 },
        },
      };

      render(<QuotaConfirmDialog />);

      fireEvent.click(screen.getByText("View Account →"));

      expect(mockClose).toHaveBeenCalledOnce();
    });
  });

  describe("would-exceed state", () => {
    it("disables generate button when generation would exceed limit", () => {
      dialogState = {
        ...defaultDialogState,
        estimatedCost: 200,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 90, reserved: 0, used: 900 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByRole("button", { name: "Generate Document" })).toBeDisabled();
    });
  });

  describe("warning level display", () => {
    it("shows warning message when usage is between 70-90%", () => {
      dialogState = {
        ...defaultDialogState,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 75, reserved: 0, used: 750 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText("You're running low on your monthly quota.")).toBeInTheDocument();
    });

    it("does not disable generate button at warning level", () => {
      dialogState = {
        ...defaultDialogState,
        estimatedCost: 10,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 75, reserved: 0, used: 750 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByRole("button", { name: "Generate Document" })).toBeEnabled();
    });
  });

  describe("danger level display", () => {
    it("shows danger message when usage is between 90-100%", () => {
      dialogState = {
        ...defaultDialogState,
        estimatedCost: 5,
        usage: {
          ...defaultDialogState.usage,
          specview: { limit: 1000, percentage: 92, reserved: 0, used: 920 },
        },
      };

      render(<QuotaConfirmDialog />);

      expect(screen.getByText("You're almost at your monthly limit.")).toBeInTheDocument();
    });
  });

  describe("confirm action", () => {
    it("calls confirm when generate button is clicked", () => {
      render(<QuotaConfirmDialog />);

      fireEvent.click(screen.getByRole("button", { name: "Generate Document" }));

      expect(mockConfirm).toHaveBeenCalledOnce();
    });

    it("enables generate button when quota is normal", () => {
      render(<QuotaConfirmDialog />);

      expect(screen.getByRole("button", { name: "Generate Document" })).toBeEnabled();
    });
  });

  describe("language selector", () => {
    it("renders language combobox with selected language", () => {
      render(<QuotaConfirmDialog />);

      const combobox = screen.getByTestId("language-combobox");
      expect(combobox).toHaveTextContent("English");
    });
  });

  describe("no usage data", () => {
    it("renders without usage section when usage is null", () => {
      dialogState = { ...defaultDialogState, usage: null };

      render(<QuotaConfirmDialog />);

      expect(screen.queryByText("Usage This Month")).not.toBeInTheDocument();
      expect(screen.getByRole("button", { name: "Generate Document" })).toBeEnabled();
    });
  });
});
