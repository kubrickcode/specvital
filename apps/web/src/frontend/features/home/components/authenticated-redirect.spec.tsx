import { render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";

import { AuthenticatedRedirect } from "./authenticated-redirect";

const mockUseAuth = vi.fn();
const mockReplace = vi.fn();

vi.mock("@/features/auth", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/i18n/navigation", () => ({
  useRouter: () => ({
    replace: mockReplace,
  }),
}));

describe("AuthenticatedRedirect", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("without children (redirect-only mode)", () => {
    it("returns null when not loading and showLoading is false", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: false,
      });

      const { container } = render(<AuthenticatedRedirect />);

      expect(container.firstChild).toBeNull();
    });

    it("shows loading fallback when isLoading and showLoading is true", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: true,
      });

      render(<AuthenticatedRedirect showLoading />);

      expect(screen.getByRole("main")).toBeInTheDocument();
    });

    it("returns null when isLoading and showLoading is false", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: true,
      });

      const { container } = render(<AuthenticatedRedirect />);

      expect(container.firstChild).toBeNull();
    });

    it("redirects to dashboard when authenticated", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: true,
        isLoading: false,
      });

      render(<AuthenticatedRedirect />);

      expect(mockReplace).toHaveBeenCalledWith("/dashboard");
    });
  });

  describe("with children (wrapper mode)", () => {
    it("shows loading fallback when auth is loading", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: true,
      });

      render(
        <AuthenticatedRedirect>
          <div>Home Content</div>
        </AuthenticatedRedirect>
      );

      expect(screen.queryByText("Home Content")).not.toBeInTheDocument();
      expect(screen.getByRole("main")).toBeInTheDocument();
    });

    it("shows loading fallback when authenticated (during redirect)", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: true,
        isLoading: false,
      });

      render(
        <AuthenticatedRedirect>
          <div>Home Content</div>
        </AuthenticatedRedirect>
      );

      expect(screen.queryByText("Home Content")).not.toBeInTheDocument();
      expect(screen.getByRole("main")).toBeInTheDocument();
    });

    it("renders children when not authenticated", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: false,
        isLoading: false,
      });

      render(
        <AuthenticatedRedirect>
          <div>Home Content</div>
        </AuthenticatedRedirect>
      );

      expect(screen.getByText("Home Content")).toBeInTheDocument();
    });

    it("redirects to dashboard when authenticated", () => {
      mockUseAuth.mockReturnValue({
        isAuthenticated: true,
        isLoading: false,
      });

      render(
        <AuthenticatedRedirect>
          <div>Home Content</div>
        </AuthenticatedRedirect>
      );

      expect(mockReplace).toHaveBeenCalledWith("/dashboard");
    });
  });
});
