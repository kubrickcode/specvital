"use client";

import { Filter, Star } from "lucide-react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Drawer,
  DrawerClose,
  DrawerContent,
  DrawerFooter,
  DrawerHeader,
  DrawerTitle,
  DrawerTrigger,
} from "@/components/ui/drawer";
import { Label } from "@/components/ui/label";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Toggle } from "@/components/ui/toggle";

import {
  OWNERSHIP_FILTER_ICONS,
  OWNERSHIP_FILTER_OPTIONS,
  type OwnershipFilter,
  useOwnershipFilter,
} from "../hooks/use-ownership-filter";
import { useStarredFilter } from "../hooks/use-starred-filter";

export const MobileFilterDrawer = () => {
  const t = useTranslations("dashboard.filter");
  const tOwnership = useTranslations("dashboard.filter.ownership");

  const { ownershipFilter, setOwnershipFilter } = useOwnershipFilter();
  const { setStarredOnly, starredOnly } = useStarredFilter();

  const [isOpen, setIsOpen] = useState(false);
  const [localOwnership, setLocalOwnership] = useState<OwnershipFilter>(ownershipFilter);
  const [localStarred, setLocalStarred] = useState(starredOnly);

  const ownershipLabels: Record<OwnershipFilter, string> = {
    all: tOwnership("all"),
    mine: tOwnership("mine"),
  };

  const activeFilterCount = calculateActiveFilterCount(ownershipFilter, starredOnly);

  const handleOpenChange = (open: boolean) => {
    if (open) {
      setLocalOwnership(ownershipFilter);
      setLocalStarred(starredOnly);
    }
    setIsOpen(open);
  };

  const handleApply = () => {
    setOwnershipFilter(localOwnership === "all" ? null : localOwnership);
    setStarredOnly(localStarred ? true : null);
    setIsOpen(false);
  };

  const handleStarredToggle = (pressed: boolean) => {
    setLocalStarred(pressed);
  };

  return (
    <Drawer onOpenChange={handleOpenChange} open={isOpen}>
      <DrawerTrigger asChild>
        <Button aria-label={t("filterButton")} className="h-11 gap-2" variant="outline">
          <Filter aria-hidden="true" className="size-4" />
          <span>{t("filterButton")}</span>
          {activeFilterCount > 0 && (
            <span className="ml-1 rounded-full bg-primary px-1.5 py-0.5 text-xs text-primary-foreground">
              {activeFilterCount}
            </span>
          )}
        </Button>
      </DrawerTrigger>
      <DrawerContent>
        <DrawerHeader>
          <DrawerTitle>{t("filterButton")}</DrawerTitle>
        </DrawerHeader>

        <div className="flex flex-col gap-6 px-4">
          <div className="flex flex-col gap-3">
            <Label className="text-sm font-medium">{tOwnership("label")}</Label>
            <RadioGroup
              onValueChange={(value) => setLocalOwnership(value as OwnershipFilter)}
              value={localOwnership}
            >
              {OWNERSHIP_FILTER_OPTIONS.map((option) => {
                const Icon = OWNERSHIP_FILTER_ICONS[option];
                return (
                  <div className="flex items-center gap-3" key={option}>
                    <RadioGroupItem id={`ownership-${option}`} value={option} />
                    <Label
                      className="flex flex-1 cursor-pointer items-center gap-2"
                      htmlFor={`ownership-${option}`}
                    >
                      <Icon aria-hidden="true" className="size-4" />
                      {ownershipLabels[option]}
                    </Label>
                  </div>
                );
              })}
            </RadioGroup>
          </div>

          <div className="flex flex-col gap-3">
            <Label className="text-sm font-medium">{t("starredLabel")}</Label>
            <Toggle
              aria-label={t("starredLabel")}
              className="h-11 w-full justify-start gap-2 px-3"
              onPressedChange={handleStarredToggle}
              pressed={localStarred}
              variant="outline"
            >
              <Star
                aria-hidden="true"
                className={localStarred ? "fill-yellow-400 text-yellow-400" : ""}
              />
              <span>{t("starred")}</span>
            </Toggle>
          </div>
        </div>

        <DrawerFooter>
          <DrawerClose asChild>
            <Button onClick={handleApply}>{t("apply")}</Button>
          </DrawerClose>
        </DrawerFooter>
      </DrawerContent>
    </Drawer>
  );
};

const calculateActiveFilterCount = (ownership: OwnershipFilter, starred: boolean): number => {
  let count = 0;
  if (ownership !== "all") count++;
  if (starred) count++;
  return count;
};
