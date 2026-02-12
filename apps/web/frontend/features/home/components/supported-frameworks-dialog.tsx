"use client";

import { useTranslations } from "next-intl";

import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { FRAMEWORK_CATEGORIES } from "@/lib/constants/frameworks";

type SupportedFrameworksDialogProps = {
  onOpenChange: (open: boolean) => void;
  open: boolean;
};

export const SupportedFrameworksDialog = ({
  onOpenChange,
  open,
}: SupportedFrameworksDialogProps) => {
  const t = useTranslations("home.frameworks");

  return (
    <Dialog onOpenChange={onOpenChange} open={open}>
      <DialogContent className="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>{t("title")}</DialogTitle>
          <DialogDescription>{t("description")}</DialogDescription>
        </DialogHeader>
        <div className="max-h-[60vh] space-y-4 overflow-y-auto pr-2">
          {FRAMEWORK_CATEGORIES.map((category) => {
            const IconComponent = category.icon;
            return (
              <div key={category.category}>
                <div className="mb-2 flex items-center gap-2">
                  <span className="text-muted-foreground">
                    <IconComponent className="size-4" />
                  </span>
                  <h3 className="text-sm font-medium">{t(category.category)}</h3>
                </div>
                <div className="flex flex-wrap gap-2">
                  {category.frameworks.map((framework) => (
                    <Badge className="bg-black/5 text-foreground dark:bg-white/10" key={framework}>
                      {framework}
                    </Badge>
                  ))}
                </div>
              </div>
            );
          })}
        </div>
      </DialogContent>
    </Dialog>
  );
};
