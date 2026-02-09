import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { SignInTrigger } from "./sign-in-trigger";

const mockOpen = vi.fn();

vi.mock("../hooks", () => ({
  useLoginModal: () => ({
    open: mockOpen,
  }),
}));

const messages = {
  auth: {
    login: "Sign in",
  },
};

const renderSignInTrigger = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <SignInTrigger />
    </NextIntlClientProvider>
  );
};

describe("SignInTrigger", () => {
  it("renders sign in button", () => {
    renderSignInTrigger();

    const button = screen.getByRole("button", { name: /sign in/i });
    expect(button).toBeInTheDocument();
  });

  it("opens login modal when clicked", () => {
    renderSignInTrigger();

    const button = screen.getByRole("button", { name: /sign in/i });
    fireEvent.click(button);

    expect(mockOpen).toHaveBeenCalled();
  });
});
