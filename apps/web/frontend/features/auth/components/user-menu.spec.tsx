import { fireEvent, render, screen } from "@testing-library/react";
import { NextIntlClientProvider } from "next-intl";
import { describe, expect, it, vi, beforeEach } from "vitest";

import messages from "@/i18n/messages/en.json";

import { UserMenu } from "./user-menu";

const mockLogout = vi.fn();
const mockUseAuth = vi.fn();

vi.mock("../hooks", () => ({
  useAuth: () => mockUseAuth(),
}));

vi.mock("@/i18n/navigation", () => ({
  Link: ({ children, ...props }: React.PropsWithChildren<{ href: string }>) => (
    <a {...props}>{children}</a>
  ),
}));

vi.mock("@/lib/background-tasks", () => ({
  TaskBadge: () => null,
  TasksDropdownSection: () => null,
}));

// Radix DropdownMenu does not open in jsdom due to missing Pointer Events API.
// Mock the UI component to render content directly for behavior testing.
vi.mock("@/components/ui/dropdown-menu", () => {
  const DropdownMenu = ({ children }: React.PropsWithChildren) => <div>{children}</div>;
  const DropdownMenuTrigger = ({
    asChild,
    children,
  }: React.PropsWithChildren<{ asChild?: boolean }>) => (
    <div data-testid="dropdown-trigger">{asChild ? children : <button>{children}</button>}</div>
  );
  const DropdownMenuContent = ({ children }: React.PropsWithChildren) => (
    <div data-testid="dropdown-content" role="menu">
      {children}
    </div>
  );
  const DropdownMenuLabel = ({
    children,
    className,
  }: React.PropsWithChildren<{ className?: string }>) => (
    <div className={className}>{children}</div>
  );
  const DropdownMenuSeparator = () => <hr />;
  const DropdownMenuItem = ({
    children,
    onClick,
    disabled,
    asChild,
  }: React.PropsWithChildren<{
    asChild?: boolean;
    disabled?: boolean;
    onClick?: () => void;
  }>) => (
    <div data-disabled={disabled || undefined} onClick={onClick} role="menuitem">
      {children}
    </div>
  );

  return {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuLabel,
    DropdownMenuSeparator,
    DropdownMenuTrigger,
  };
});

const renderWithProvider = (ui: React.ReactElement) =>
  render(
    <NextIntlClientProvider locale="en" messages={messages}>
      {ui}
    </NextIntlClientProvider>
  );

const authenticatedUser = {
  logout: mockLogout,
  logoutPending: false,
  user: {
    avatarUrl: "https://avatars.githubusercontent.com/u/1?v=4",
    login: "testuser",
    name: "Test User",
  },
};

describe("UserMenu", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockUseAuth.mockReturnValue(authenticatedUser);
  });

  it("renders nothing when user is null", () => {
    mockUseAuth.mockReturnValue({
      logout: vi.fn(),
      logoutPending: false,
      user: null,
    });

    const { container } = renderWithProvider(<UserMenu />);

    expect(container.firstChild).toBeNull();
  });

  it("renders avatar image with alt text matching login", () => {
    renderWithProvider(<UserMenu />);

    const img = screen.getByAltText("testuser");
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute("src", "https://avatars.githubusercontent.com/u/1?v=4");
  });

  it("renders screen-reader label for the trigger button", () => {
    renderWithProvider(<UserMenu />);

    expect(screen.getByText("Open user menu")).toBeInTheDocument();
  });

  it("displays username with @ prefix in the dropdown", () => {
    renderWithProvider(<UserMenu />);

    expect(screen.getByText("@testuser")).toBeInTheDocument();
    expect(screen.getByText("Test User")).toBeInTheDocument();
  });

  it("displays Account link and Sign out menu item", () => {
    renderWithProvider(<UserMenu />);

    expect(screen.getByText("Account")).toBeInTheDocument();
    expect(screen.getByText("Sign out")).toBeInTheDocument();
  });

  it("renders Account as a link to /account", () => {
    renderWithProvider(<UserMenu />);

    const accountLink = screen.getByText("Account").closest("a");
    expect(accountLink).toHaveAttribute("href", "/account");
  });

  it("calls logout when Sign out is clicked", () => {
    renderWithProvider(<UserMenu />);

    const signOutItem = screen.getByText("Sign out").closest("[role='menuitem']")!;
    fireEvent.click(signOutItem);

    expect(mockLogout).toHaveBeenCalledOnce();
  });
});
