import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it } from "vitest";

import { type DisplayUpdateStatus, UpdateStatusBadge } from "./update-status-badge";

const messages = {
  dashboard: {
    status: {
      analyzing: "Analyzing...",
      newCommits: "New commits",
      unknown: "Status unknown",
      upToDate: "Up to date",
    },
  },
};

const renderBadge = (status: DisplayUpdateStatus) => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <UpdateStatusBadge status={status} />
    </NextIntlClientProvider>
  );
};

describe("UpdateStatusBadge", () => {
  it("renders up-to-date status", () => {
    renderBadge("up-to-date");

    expect(screen.getByText("Up to date")).toBeInTheDocument();
  });

  it("renders new-commits status", () => {
    renderBadge("new-commits");

    expect(screen.getByText("New commits")).toBeInTheDocument();
  });

  it("renders analyzing status with spinner", () => {
    const { container } = renderBadge("analyzing");

    expect(screen.getByText("Analyzing...")).toBeInTheDocument();

    const spinnerIcon = container.querySelector(".animate-spin");
    expect(spinnerIcon).toBeInTheDocument();
  });

  it("renders unknown status", () => {
    renderBadge("unknown");

    expect(screen.getByText("Status unknown")).toBeInTheDocument();
  });
});
