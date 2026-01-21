import type { ReactNode } from "react";

import { DocsSidebar } from "@/features/docs";

type HowItWorksLayoutProps = {
  children: ReactNode;
};

const HowItWorksLayout = ({ children }: HowItWorksLayoutProps) => {
  return (
    <div className="container mx-auto flex max-w-6xl gap-8 px-4 py-8">
      <DocsSidebar />
      <main className="min-w-0 flex-1">{children}</main>
    </div>
  );
};

export default HowItWorksLayout;
