/**
 * Mock data for GitHub organizations API
 * Endpoint: /api/user/github/organizations
 */

import type { GitHubOrganization, GitHubOrganizationsResponse } from "./types";

export const mockOrganizations: GitHubOrganization[] = [
  {
    id: 1001,
    login: "test-org",
    avatarUrl: "https://avatars.githubusercontent.com/u/1001",
    description: "Test organization with app installed",
    accessStatus: "accessible",
  },
  {
    id: 1002,
    login: "another-org",
    avatarUrl: "https://avatars.githubusercontent.com/u/1002",
    description: "Another org with app installed",
    accessStatus: "accessible",
  },
  {
    id: 1003,
    login: "restricted-org",
    avatarUrl: "https://avatars.githubusercontent.com/u/1003",
    description: "Organization without GitHub App",
    accessStatus: "restricted",
  },
  {
    id: 1004,
    login: "pending-org",
    avatarUrl: "https://avatars.githubusercontent.com/u/1004",
    description: "Organization with pending installation",
    accessStatus: "pending",
  },
];

export const mockOrganizationsResponse: GitHubOrganizationsResponse = {
  data: mockOrganizations,
};

// Empty organizations
export const mockOrganizationsEmpty: GitHubOrganizationsResponse = {
  data: [],
};

// Only restricted organizations
export const mockOrganizationsAllRestricted: GitHubOrganizationsResponse = {
  data: mockOrganizations.filter((org) => org.accessStatus === "restricted"),
};
