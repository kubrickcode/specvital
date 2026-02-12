import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { beforeEach, describe, expect, it, vi } from "vitest";

import { LoginRequiredState } from "./login-required-state";

const mockOpen = vi.fn();

vi.mock("@/features/auth", () => ({
  useLoginModal: () => ({
    open: mockOpen,
  }),
}));

const messages = {
  explore: {
    loginRequired: {
      myReposDescription:
        "Connect your GitHub account to see and analyze your personal repositories.",
      myReposTitle: "Sign in to view your repositories",
      organizationsDescription:
        "Connect your GitHub account to access your organization repositories.",
      organizationsTitle: "Sign in to view organizations",
      signIn: "Sign in to continue",
    },
  },
};

const renderLoginRequired = (props: { descriptionKey: string; titleKey: string }) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <LoginRequiredState {...props} />
    </NextIntlClientProvider>
  );

describe("LoginRequiredState", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("displays title and description for My Repos", () => {
    renderLoginRequired({
      descriptionKey: "myReposDescription",
      titleKey: "myReposTitle",
    });

    expect(screen.getByText("Sign in to view your repositories")).toBeInTheDocument();
    expect(
      screen.getByText("Connect your GitHub account to see and analyze your personal repositories.")
    ).toBeInTheDocument();
  });

  it("displays title and description for Organizations", () => {
    renderLoginRequired({
      descriptionKey: "organizationsDescription",
      titleKey: "organizationsTitle",
    });

    expect(screen.getByText("Sign in to view organizations")).toBeInTheDocument();
    expect(
      screen.getByText("Connect your GitHub account to access your organization repositories.")
    ).toBeInTheDocument();
  });

  it("renders sign-in button that opens login modal", () => {
    renderLoginRequired({
      descriptionKey: "myReposDescription",
      titleKey: "myReposTitle",
    });

    const signInButton = screen.getByRole("button", { name: /sign in to continue/i });
    expect(signInButton).toBeInTheDocument();

    fireEvent.click(signInButton);

    expect(mockOpen).toHaveBeenCalledOnce();
  });
});
