type RepositoryGridProps = {
  ariaLabel?: string;
  children: React.ReactNode;
  isLoading?: boolean;
};

export const RepositoryGrid = ({ ariaLabel, children, isLoading = false }: RepositoryGridProps) => {
  return (
    <ul
      aria-busy={isLoading}
      aria-label={ariaLabel}
      className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 auto-rows-fr [&>li]:h-full"
    >
      {children}
    </ul>
  );
};
