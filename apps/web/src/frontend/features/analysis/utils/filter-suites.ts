import type { TestCase, TestSuite } from "@/lib/api";

import type { FilterState } from "../types";

const matchesQuery = (test: TestCase, searchTerms: string[]): boolean => {
  if (searchTerms.length === 0) {
    return true;
  }
  const searchableText = `${test.name} ${test.filePath}`.toLowerCase();
  return searchTerms.every((term) => searchableText.includes(term));
};

const matchesStatus = (test: TestCase, statuses: FilterState["statuses"]): boolean => {
  if (statuses.length === 0) {
    return true;
  }
  return statuses.includes(test.status);
};

const matchesFramework = (test: TestCase, frameworks: string[]): boolean => {
  if (frameworks.length === 0) {
    return true;
  }
  return frameworks.includes(test.framework);
};

const matchesAllFilters = (test: TestCase, filter: FilterState, searchTerms: string[]): boolean => {
  return (
    matchesQuery(test, searchTerms) &&
    matchesStatus(test, filter.statuses) &&
    matchesFramework(test, filter.frameworks)
  );
};

export const filterSuites = (suites: TestSuite[], filter: FilterState): TestSuite[] => {
  const { frameworks, query, statuses } = filter;

  const hasNoFilters = !query.trim() && statuses.length === 0 && frameworks.length === 0;
  if (hasNoFilters) {
    return suites;
  }

  const searchTerms = query.toLowerCase().split(/\s+/).filter(Boolean);

  return suites.reduce<TestSuite[]>((acc, suite) => {
    const matchingTests = suite.tests.filter((test) =>
      matchesAllFilters(test, filter, searchTerms)
    );

    const filePathMatches =
      searchTerms.length > 0 &&
      searchTerms.every((term) => suite.filePath.toLowerCase().includes(term));
    const suiteNameMatches =
      suite.suiteName &&
      searchTerms.length > 0 &&
      searchTerms.every((term) => suite.suiteName.toLowerCase().includes(term));

    if (matchingTests.length > 0) {
      acc.push({ ...suite, tests: matchingTests });
    } else if (
      (filePathMatches || suiteNameMatches) &&
      statuses.length === 0 &&
      frameworks.length === 0
    ) {
      acc.push(suite);
    }

    return acc;
  }, []);
};
