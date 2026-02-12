import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { RequireAuth } from "./require-auth";

const mockUseAuth = vi.fn();
const mockReplace = vi.fn();

vi.mock("../hooks/use-auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/i18n/navigation", () => ({
  useRouter: () => ({
    replace: mockReplace,
  }),
}));

describe("RequireAuth", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows loading fallback when auth is loading", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
    });

    render(
      <RequireAuth>
        <div>Protected Content</div>
      </RequireAuth>
    );

    expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
    expect(screen.getByRole("main")).toBeInTheDocument();
  });

  it("shows loading fallback with custom message when provided", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: true,
    });

    render(
      <RequireAuth loadingMessage="Checking authentication...">
        <div>Protected Content</div>
      </RequireAuth>
    );

    expect(screen.getByText("Checking authentication...")).toBeInTheDocument();
  });

  it("renders children when authenticated", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: true,
      isLoading: false,
    });

    render(
      <RequireAuth>
        <div>Protected Content</div>
      </RequireAuth>
    );

    expect(screen.getByText("Protected Content")).toBeInTheDocument();
  });

  it("redirects to home when not authenticated", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
    });

    render(
      <RequireAuth>
        <div>Protected Content</div>
      </RequireAuth>
    );

    expect(mockReplace).toHaveBeenCalledWith("/");
  });

  it("shows loading fallback while redirecting unauthenticated user", () => {
    mockUseAuth.mockReturnValue({
      isAuthenticated: false,
      isLoading: false,
    });

    render(
      <RequireAuth>
        <div>Protected Content</div>
      </RequireAuth>
    );

    expect(screen.queryByText("Protected Content")).not.toBeInTheDocument();
    expect(screen.getByRole("main")).toBeInTheDocument();
  });
});
