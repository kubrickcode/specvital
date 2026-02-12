"use client";

import { ChevronDown, HelpCircle } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import { Collapsible, CollapsibleContent, CollapsibleTrigger } from "@/components/ui/collapsible";
import { cn } from "@/lib/utils";

import { STATUS_CONFIG, STATUS_ORDER } from "../constants/status-config";

export const StatusLegend = () => {
  const t = useTranslations("specView.statusLegend");
  const [isOpen, setIsOpen] = useState(false);

  return (
    <Collapsible onOpenChange={setIsOpen} open={isOpen}>
      <CollapsibleTrigger asChild>
        <Button
          className="gap-1.5 text-muted-foreground hover:text-foreground"
          size="sm"
          variant="ghost"
        >
          <HelpCircle className="h-4 w-4" />
          <span className="text-xs">{t("title")}</span>
          <ChevronDown
            className={cn("h-3.5 w-3.5 transition-transform duration-200", isOpen && "rotate-180")}
          />
        </Button>
      </CollapsibleTrigger>
      <CollapsibleContent className="overflow-hidden data-[state=closed]:animate-collapse-up data-[state=open]:animate-collapse-down">
        <div className="mt-3 p-3 bg-muted/50 rounded-lg border">
          <div className="grid gap-2">
            {STATUS_ORDER.map((status) => {
              const config = STATUS_CONFIG[status];
              const Icon = config.icon;

              return (
                <div className="flex items-center gap-3" key={status}>
                  <span
                    className={cn(
                      "flex items-center justify-center w-6 h-6 rounded-full flex-shrink-0",
                      "ring-1",
                      config.bgColor,
                      config.ringColor
                    )}
                  >
                    <Icon className={cn("h-3.5 w-3.5", config.color)} />
                  </span>
                  <div className="flex-1 min-w-0">
                    <span className="text-sm font-medium">{t(config.labelKey)}</span>
                    <p className="text-xs text-muted-foreground">{t(config.descriptionKey)}</p>
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </CollapsibleContent>
    </Collapsible>
  );
};
