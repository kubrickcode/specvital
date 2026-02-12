import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { DataViewToggle } from "./data-view-toggle";

const messages = {
  analyze: {
    viewMode: {
      list: "List view",
      tree: "Tree view",
    },
  },
};

const renderDataViewToggle = (props: Partial<React.ComponentProps<typeof DataViewToggle>> = {}) => {
  const defaultProps = {
    onChange: vi.fn(),
    value: "list" as const,
    ...props,
  };

  return render(
    <NextIntlClientProvider locale="en" messages={messages} timeZone="UTC">
      <DataViewToggle {...defaultProps} />
    </NextIntlClientProvider>
  );
};

describe("DataViewToggle", () => {
  it("renders both list and tree toggle options in Tests tab", () => {
    renderDataViewToggle();

    expect(screen.getByRole("radio", { name: "List view" })).toBeInTheDocument();
    expect(screen.getByRole("radio", { name: "Tree view" })).toBeInTheDocument();
  });
});
