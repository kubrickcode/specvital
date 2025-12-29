import type { TestSuite } from "@/lib/api";

import type { FilterState } from "../types";

export const filterSuites = (suites: TestSuite[], filter: FilterState): TestSuite[] => {
  const { query } = filter;

  if (!query.trim()) {
    return suites;
  }

  const searchTerms = query.toLowerCase().split(/\s+/).filter(Boolean);

  return suites.reduce<TestSuite[]>((acc, suite) => {
    const matchingTests = suite.tests.filter((test) => {
      const searchableText = `${test.name} ${test.filePath}`.toLowerCase();
      return searchTerms.every((term) => searchableText.includes(term));
    });

    const filePathMatches = searchTerms.every((term) =>
      suite.filePath.toLowerCase().includes(term)
    );
    const suiteNameMatches =
      suite.suiteName && searchTerms.every((term) => suite.suiteName.toLowerCase().includes(term));

    if (matchingTests.length > 0) {
      acc.push({ ...suite, tests: matchingTests });
    } else if (filePathMatches || suiteNameMatches) {
      acc.push(suite);
    }

    return acc;
  }, []);
};
