import type { TestSuite } from "@/lib/api";
import { TestSuiteAccordion } from "./test-suite-accordion";

type TestListProps = {
  suites: TestSuite[];
};

export const TestList = ({ suites }: TestListProps) => {
  if (suites.length === 0) {
    return (
      <div className="rounded-lg border bg-card p-6 text-center text-muted-foreground">
        No test suites found in this repository.
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {suites.map((suite) => (
        <TestSuiteAccordion key={`${suite.filePath}-${suite.framework}`} suite={suite} />
      ))}
    </div>
  );
};
