import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { AuthLoadingWrapper } from "./auth-loading-wrapper";

const mockUseAuth = vi.fn();

vi.mock("../hooks/use-auth", () => ({
  useAuth: () => mockUseAuth(),
}));

describe("AuthLoadingWrapper", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("shows loading fallback when auth is loading", () => {
    mockUseAuth.mockReturnValue({
      isLoading: true,
    });

    render(
      <AuthLoadingWrapper>
        <div>Content</div>
      </AuthLoadingWrapper>
    );

    expect(screen.queryByText("Content")).not.toBeInTheDocument();
    expect(screen.getByRole("main")).toBeInTheDocument();
  });

  it("shows loading fallback with custom message when provided", () => {
    mockUseAuth.mockReturnValue({
      isLoading: true,
    });

    render(
      <AuthLoadingWrapper loadingMessage="Loading...">
        <div>Content</div>
      </AuthLoadingWrapper>
    );

    expect(screen.getByText("Loading...")).toBeInTheDocument();
  });

  it("renders children when auth check completes", () => {
    mockUseAuth.mockReturnValue({
      isLoading: false,
    });

    render(
      <AuthLoadingWrapper>
        <div>Content</div>
      </AuthLoadingWrapper>
    );

    expect(screen.getByText("Content")).toBeInTheDocument();
  });
});
