import { cn } from "@/lib/utils";

type StatsCardProps = {
  active: number;
  label: string;
  skipped: number;
  todo: number;
  total: number;
};

export const StatsCard = ({ active, label, skipped, todo, total }: StatsCardProps) => {
  return (
    <div className={cn("rounded-lg border bg-card p-6 shadow-xs")}>
      <h3 className="text-lg font-semibold mb-4">{label}</h3>
      <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-foreground">{total}</span>
          <span className="text-sm text-muted-foreground">Total</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-green-600">{active}</span>
          <span className="text-sm text-muted-foreground">Active</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-amber-600">{skipped}</span>
          <span className="text-sm text-muted-foreground">Skipped</span>
        </div>
        <div className="flex flex-col">
          <span className="text-3xl font-bold text-blue-600">{todo}</span>
          <span className="text-sm text-muted-foreground">Todo</span>
        </div>
      </div>
    </div>
  );
};
