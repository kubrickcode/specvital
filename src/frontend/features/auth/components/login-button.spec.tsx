import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { LoginButton } from "./login-button";

const mockLogin = vi.fn();
let mockLoginPending = false;

vi.mock("../hooks", () => ({
  useAuth: () => ({
    login: mockLogin,
    loginPending: mockLoginPending,
  }),
}));

const messages = {
  auth: {
    login: "Sign in",
  },
};

const renderLoginButton = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <LoginButton />
    </NextIntlClientProvider>
  );
};

describe("LoginButton", () => {
  it("renders login button with GitHub icon", () => {
    renderLoginButton();

    const button = screen.getByRole("button", { name: /sign in/i });
    expect(button).toBeInTheDocument();
  });

  it("calls login when clicked", () => {
    renderLoginButton();

    const button = screen.getByRole("button", { name: /sign in/i });
    fireEvent.click(button);

    expect(mockLogin).toHaveBeenCalled();
  });

  it("is disabled when login is pending", () => {
    mockLoginPending = true;
    renderLoginButton();

    const button = screen.getByRole("button", { name: /sign in/i });
    expect(button).toBeDisabled();

    mockLoginPending = false;
  });
});
