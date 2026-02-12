import { cleanup, fireEvent, render, screen, waitFor } from "@testing-library/react";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { LanguageCombobox } from "./language-combobox";

vi.mock("next-intl", () => ({
  useTranslations: () => (key: string) => {
    const messages: Record<string, string> = {
      noLanguageFound: "No language found.",
      searchLanguage: "Search languages...",
      selectLanguage: "Select language",
    };
    return messages[key] || key;
  },
}));

describe("LanguageCombobox", () => {
  const mockOnValueChange = vi.fn();

  beforeEach(() => {
    mockOnValueChange.mockClear();
  });

  afterEach(() => {
    cleanup();
  });

  it("should render with selected language display label", () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="Korean" />);
    expect(screen.getByRole("combobox")).toHaveTextContent("한국어 (Korean)");
  });

  it("should open dropdown when clicked", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="English" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search languages...")).toBeInTheDocument();
    });
  });

  it("should filter languages by English name", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="English" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search languages...")).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText("Search languages...");
    fireEvent.change(input, { target: { value: "kor" } });

    await waitFor(() => {
      const options = screen.getAllByRole("option");
      expect(options.length).toBe(1);
      expect(options[0]).toHaveTextContent("한국어 (Korean)");
    });
  });

  it("should filter languages by native name", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="English" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search languages...")).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText("Search languages...");
    fireEvent.change(input, { target: { value: "한국" } });

    await waitFor(() => {
      const options = screen.getAllByRole("option");
      expect(options.length).toBe(1);
      expect(options[0]).toHaveTextContent("한국어 (Korean)");
    });
  });

  it("should call onValueChange when language is selected", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="English" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search languages...")).toBeInTheDocument();
    });

    const koreanOption = screen.getByText("한국어 (Korean)");
    fireEvent.click(koreanOption);

    expect(mockOnValueChange).toHaveBeenCalledWith("Korean");
  });

  it("should show check icon for selected language", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="Korean" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      const koreanOption = screen.getByRole("option", { name: /한국어 \(Korean\)/i });
      expect(koreanOption.querySelector("svg.opacity-100")).toBeInTheDocument();
    });
  });

  it("should show empty message when no match found", async () => {
    render(<LanguageCombobox onValueChange={mockOnValueChange} value="English" />);

    fireEvent.click(screen.getByRole("combobox"));

    await waitFor(() => {
      expect(screen.getByPlaceholderText("Search languages...")).toBeInTheDocument();
    });

    const input = screen.getByPlaceholderText("Search languages...");
    fireEvent.change(input, { target: { value: "xyz" } });

    await waitFor(() => {
      expect(screen.getByText("No language found.")).toBeInTheDocument();
    });
  });
});
