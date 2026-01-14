"use client";

import * as React from "react";

import { Popover, PopoverContent, PopoverTrigger } from "@/components/ui/popover";
import { Tooltip, TooltipContent, TooltipTrigger } from "@/components/ui/tooltip";
import { useMediaQuery } from "@/lib/hooks";
import { cn } from "@/lib/utils";

type ResponsiveTooltipProps = {
  children: React.ReactNode;
  className?: string;
  content: React.ReactNode;
  contentClassName?: string;
  open?: boolean;
  side?: "top" | "right" | "bottom" | "left";
  sideOffset?: number;
};

export const ResponsiveTooltip = ({
  children,
  className,
  content,
  contentClassName,
  open,
  side = "top",
  sideOffset = 4,
}: ResponsiveTooltipProps) => {
  const hasHover = useMediaQuery("(hover: hover) and (pointer: fine)");

  if (hasHover) {
    return (
      <Tooltip open={open}>
        <TooltipTrigger asChild className={className}>
          {children}
        </TooltipTrigger>
        <TooltipContent className={contentClassName} side={side} sideOffset={sideOffset}>
          {content}
        </TooltipContent>
      </Tooltip>
    );
  }

  return (
    <Popover open={open}>
      <PopoverTrigger asChild className={className}>
        {children}
      </PopoverTrigger>
      <PopoverContent
        className={cn(
          "w-auto max-w-[90vw] border-0 bg-foreground/95 px-3 py-1.5 text-xs text-background shadow-lg shadow-black/10 backdrop-blur-sm break-words",
          contentClassName
        )}
        side={side}
        sideOffset={sideOffset}
      >
        {content}
      </PopoverContent>
    </Popover>
  );
};
