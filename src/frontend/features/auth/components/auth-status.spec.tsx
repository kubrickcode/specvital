import { render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi } from "vitest";

import { AuthStatus } from "./auth-status";

const mockUseAuth = vi.fn();

vi.mock("../hooks", () => ({
  useAuth: () => mockUseAuth(),
}));

const messages = {
  auth: {
    login: "Sign in",
    logout: "Sign out",
    openUserMenu: "Open user menu",
  },
};

const renderAuthStatus = () => {
  return render(
    <NextIntlClientProvider locale="en" messages={messages}>
      <AuthStatus />
    </NextIntlClientProvider>
  );
};

describe("AuthStatus", () => {
  it("shows loading spinner when loading", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
      user: null,
    });

    renderAuthStatus();

    expect(screen.queryByRole("button")).not.toBeInTheDocument();
  });

  it("shows login button when not authenticated", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
      login: vi.fn(),
      loginPending: false,
      user: null,
    });

    renderAuthStatus();

    expect(screen.getByRole("button", { name: /sign in/i })).toBeInTheDocument();
  });

  it("shows user menu when authenticated", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
      logout: vi.fn(),
      logoutPending: false,
      user: {
        avatarUrl: "https://example.com/avatar.png",
        id: "1",
        login: "testuser",
        name: "Test User",
      },
    });

    renderAuthStatus();

    expect(screen.getByRole("button", { name: /open user menu/i })).toBeInTheDocument();
  });
});
