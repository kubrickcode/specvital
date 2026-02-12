import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import type { SpecDocument } from "../types";

import { TocSidebar } from "./toc-sidebar";

const mockScrollToSection = vi.fn();

vi.mock("../contexts", () => ({
  useDocumentNavigation: () => ({
    activeSection: null,
    scrollToSection: mockScrollToSection,
  }),
}));

vi.mock("@/lib/hooks", () => ({
  useTruncateDetection: () => ({ isTruncated: false, ref: { current: null } }),
}));

const messages = {
  specView: {
    stats: {
      behaviors: "{count} behaviors",
      domains: "{count} domains",
      features: "{count} features",
    },
    toc: {
      description: "Table of contents for the spec document",
      filterMatchCount: "{matched} / {total} matched",
      openToc: "Open table of contents",
      title: "Table of Contents",
    },
  },
};

const createDocument = (): SpecDocument => ({
  analysisId: "analysis-1",
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
            {
              convertedDescription: "Logout works",
              id: "b2",
              originalName: "should logout",
              sortOrder: 2,
              testCaseId: "tc2",
            },
          ],
          description: "Login feature",
          id: "f1",
          name: "Login",
          sortOrder: 1,
        },
        {
          behaviors: [
            {
              convertedDescription: "Register works",
              id: "b3",
              originalName: "should register",
              sortOrder: 1,
              testCaseId: "tc3",
            },
          ],
          description: "Registration feature",
          id: "f2",
          name: "Registration",
          sortOrder: 2,
        },
      ],
      id: "d1",
      name: "Authentication",
      sortOrder: 1,
    },
    {
      classificationConfidence: 0.85,
      description: "Payment domain",
      features: [
        {
          behaviors: [
            {
              convertedDescription: "Checkout works",
              id: "b4",
              originalName: "should checkout",
              sortOrder: 1,
              testCaseId: "tc4",
            },
          ],
          description: "Checkout feature",
          id: "f3",
          name: "Checkout",
          sortOrder: 1,
        },
      ],
      id: "d2",
      name: "Payment",
      sortOrder: 2,
    },
  ],
  executiveSummary: "Summary",
  id: "doc-1",
  language: "English",
  version: 1,
});

const renderTocSidebar = (props: Partial<React.ComponentProps<typeof TocSidebar>> = {}) => {
  const defaultProps = {
    document: createDocument(),
    ...props,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <TocSidebar {...defaultProps} />
    </NextIntlClientProvider>
  );
};

describe("TocSidebar", () => {
  it("displays domain and feature items", () => {
    renderTocSidebar();

    // Domain names should be visible
    expect(screen.getByText("Authentication")).toBeInTheDocument();
    expect(screen.getByText("Payment")).toBeInTheDocument();
  });

  it("calls scrollToSection when a domain item is clicked", () => {
    mockScrollToSection.mockClear();
    renderTocSidebar();

    // Click the Authentication domain
    fireEvent.click(screen.getByText("Authentication"));

    expect(mockScrollToSection).toHaveBeenCalledWith("domain-d1");
  });
});
