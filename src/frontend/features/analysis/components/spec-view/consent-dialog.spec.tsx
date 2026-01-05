import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { ConsentDialog } from "./consent-dialog";

const messages = {
  analyze: {
    specView: {
      cancel: "Cancel",
      consentDescription: "AI will convert test names to natural language.",
      consentTitle: "AI Conversion",
      convert: "Convert",
      languageLabel: "Language",
    },
  },
};

const renderWithIntl = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

describe("ConsentDialog", () => {
  it("should render dialog when open", () => {
    const onCancel = vi.fn();
    const onConsent = vi.fn();

    renderWithIntl(<ConsentDialog onCancel={onCancel} onConsent={onConsent} open={true} />);

    expect(screen.getByText("AI Conversion")).toBeInTheDocument();
    expect(screen.getByText("AI will convert test names to natural language.")).toBeInTheDocument();
  });

  it("should not render dialog when closed", () => {
    const onCancel = vi.fn();
    const onConsent = vi.fn();

    renderWithIntl(<ConsentDialog onCancel={onCancel} onConsent={onConsent} open={false} />);

    expect(screen.queryByText("AI Conversion")).not.toBeInTheDocument();
  });

  it("should call onCancel when cancel button is clicked", () => {
    const onCancel = vi.fn();
    const onConsent = vi.fn();

    renderWithIntl(<ConsentDialog onCancel={onCancel} onConsent={onConsent} open={true} />);

    fireEvent.click(screen.getByText("Cancel"));
    expect(onCancel).toHaveBeenCalled();
  });

  it("should call onConsent with default language when convert button is clicked", () => {
    const onCancel = vi.fn();
    const onConsent = vi.fn();

    renderWithIntl(<ConsentDialog onCancel={onCancel} onConsent={onConsent} open={true} />);

    fireEvent.click(screen.getByText("Convert"));
    expect(onConsent).toHaveBeenCalledWith("en");
  });

  it("should show language selector", () => {
    const onCancel = vi.fn();
    const onConsent = vi.fn();

    renderWithIntl(<ConsentDialog onCancel={onCancel} onConsent={onConsent} open={true} />);

    expect(screen.getByText("Language")).toBeInTheDocument();
    expect(screen.getByText("English")).toBeInTheDocument();
  });
});
