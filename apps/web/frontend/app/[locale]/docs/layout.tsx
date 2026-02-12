import type { ReactNode } from "react";

type DocsLayoutProps = {
  children: ReactNode;
};

const DocsLayout = ({ children }: DocsLayoutProps) => {
  return <>{children}</>;
};

export default DocsLayout;
