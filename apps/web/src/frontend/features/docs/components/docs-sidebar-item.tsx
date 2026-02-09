"use client";

import type { LucideIcon } from "lucide-react";

import { Link, usePathname } from "@/i18n/navigation";
import { cn } from "@/lib/utils";

type DocsSidebarItemProps = {
  href: string;
  icon: LucideIcon;
  title: string;
};

export const DocsSidebarItem = ({ href, icon: Icon, title }: DocsSidebarItemProps) => {
  const pathname = usePathname();
  const isActive = pathname === href;

  return (
    <Link
      className={cn(
        "flex items-center gap-3 rounded-md px-3 py-2 text-sm transition-colors",
        isActive
          ? "bg-gradient-to-r from-foreground/[0.10] to-transparent font-medium text-foreground"
          : "hover:bg-muted/50 text-muted-foreground"
      )}
      href={href}
    >
      <Icon className="size-4 flex-shrink-0" />
      <span>{title}</span>
    </Link>
  );
};
