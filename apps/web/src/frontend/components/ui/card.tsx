import { cva, type VariantProps } from "class-variance-authority";
import * as React from "react";

import { cn } from "@/lib/utils/index";

const cardVariants = cva("bg-card text-card-foreground flex flex-col gap-6 rounded-xl py-6", {
  defaultVariants: {
    depth: "default",
  },
  variants: {
    depth: {
      default: "border shadow-sm hover:shadow-md transition-shadow duration-200",
      elevated:
        "border border-border/50 shadow-xl shadow-black/15 hover:shadow-2xl hover:shadow-black/20 hover:-translate-y-px transition-all duration-200 motion-reduce:hover:translate-y-0 dark:shadow-black/40 dark:hover:shadow-black/50",
      floating:
        "border border-border/50 shadow-xl shadow-black/10 hover:shadow-2xl hover:shadow-black/15 hover:-translate-y-1 transition-all duration-200 motion-reduce:hover:translate-y-0 dark:shadow-black/40 dark:hover:shadow-black/50",
      raised:
        "border border-border/50 shadow-md shadow-black/8 hover:shadow-lg hover:shadow-black/12 transition-all duration-200 dark:shadow-black/25 dark:hover:shadow-black/35",
    },
  },
});

type CardProps = React.ComponentProps<"div"> & VariantProps<typeof cardVariants>;

function Card({ className, depth, ...props }: CardProps) {
  return <div className={cn(cardVariants({ depth }), className)} data-slot="card" {...props} />;
}

function CardHeader({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn(
        "@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6",
        className
      )}
      data-slot="card-header"
      {...props}
    />
  );
}

function CardTitle({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("leading-none font-semibold", className)}
      data-slot="card-title"
      {...props}
    />
  );
}

function CardDescription({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("text-muted-foreground text-sm", className)}
      data-slot="card-description"
      {...props}
    />
  );
}

function CardAction({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("col-start-2 row-span-2 row-start-1 self-start justify-self-end", className)}
      data-slot="card-action"
      {...props}
    />
  );
}

function CardContent({ className, ...props }: React.ComponentProps<"div">) {
  return <div className={cn("px-6", className)} data-slot="card-content" {...props} />;
}

function CardFooter({ className, ...props }: React.ComponentProps<"div">) {
  return (
    <div
      className={cn("flex items-center px-6 [.border-t]:pt-6", className)}
      data-slot="card-footer"
      {...props}
    />
  );
}

export { Card, CardHeader, CardFooter, CardTitle, CardAction, CardDescription, CardContent };
