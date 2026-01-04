export type StatusCounts = {
  active: number;
  skipped: number;
  todo: number;
};

export const calculateStatusCounts = (tests: { status: string }[]): StatusCounts => {
  const counts: StatusCounts = { active: 0, skipped: 0, todo: 0 };

  for (const test of tests) {
    switch (test.status) {
      case "active":
      case "focused":
        counts.active++;
        break;
      case "skipped":
      case "xfail":
        counts.skipped++;
        break;
      case "todo":
        counts.todo++;
        break;
    }
  }

  return counts;
};
