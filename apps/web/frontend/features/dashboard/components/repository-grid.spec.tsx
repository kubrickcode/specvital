import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";

import { RepositoryGrid } from "./repository-grid";

describe("RepositoryGrid", () => {
  it("renders children inside a labeled list", () => {
    render(
      <RepositoryGrid ariaLabel="Repository list">
        <li>Repo A</li>
        <li>Repo B</li>
      </RepositoryGrid>
    );

    const list = screen.getByRole("list", { name: "Repository list" });

    expect(list).toBeInTheDocument();
    expect(screen.getByText("Repo A")).toBeInTheDocument();
    expect(screen.getByText("Repo B")).toBeInTheDocument();
  });

  it("sets aria-busy when loading", () => {
    render(
      <RepositoryGrid ariaLabel="Loading repositories" isLoading>
        <li>Skeleton</li>
      </RepositoryGrid>
    );

    expect(screen.getByRole("list")).toHaveAttribute("aria-busy", "true");
  });

  it("sets aria-busy false when not loading", () => {
    render(
      <RepositoryGrid ariaLabel="Repository list">
        <li>Repo</li>
      </RepositoryGrid>
    );

    expect(screen.getByRole("list")).toHaveAttribute("aria-busy", "false");
  });
});
