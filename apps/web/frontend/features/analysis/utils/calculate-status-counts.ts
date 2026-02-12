export type StatusCounts = {
  active: number;
  focused: number;
  skipped: number;
  todo: number;
  xfail: number;
};

export const calculateStatusCounts = (tests: { status: string }[]): StatusCounts => {
  const counts: StatusCounts = { active: 0, focused: 0, skipped: 0, todo: 0, xfail: 0 };

  for (const test of tests) {
    switch (test.status) {
      case "active":
        counts.active++;
        break;
      case "focused":
        counts.focused++;
        break;
      case "skipped":
        counts.skipped++;
        break;
      case "todo":
        counts.todo++;
        break;
      case "xfail":
        counts.xfail++;
        break;
    }
  }

  return counts;
};
