"use client";

import { BookOpen, CreditCard, Layers, Lightbulb, Menu } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { ScrollArea } from "@/components/ui/scroll-area";
import {
  Sheet,
  SheetContent,
  SheetDescription,
  SheetHeader,
  SheetTitle,
  SheetTrigger,
} from "@/components/ui/sheet";

import { DocsSidebarItem } from "./docs-sidebar-item";

const SIDEBAR_ITEMS = [
  { icon: Layers, key: "testDetection" },
  { icon: CreditCard, key: "usageBilling" },
  { icon: BookOpen, key: "specviewGeneration" },
  { icon: Lightbulb, key: "writingGuide" },
] as const;

type SidebarContentProps = {
  t: ReturnType<typeof useTranslations<"docs">>;
};

const SidebarContent = ({ t }: SidebarContentProps) => {
  return (
    <div className="flex h-full flex-col">
      <div className="border-b px-4 py-3">
        <h2 className="text-sm font-semibold">{t("sidebar.title")}</h2>
      </div>
      <ScrollArea className="flex-1">
        <nav aria-label={t("sidebar.ariaLabel")} className="space-y-1 p-2">
          {SIDEBAR_ITEMS.map(({ icon, key }) => (
            <DocsSidebarItem
              href={`/docs/${t(`topics.${key}.slug`)}`}
              icon={icon}
              key={key}
              title={t(`topics.${key}.title`)}
            />
          ))}
        </nav>
      </ScrollArea>
    </div>
  );
};

export const DocsSidebar = () => {
  const t = useTranslations("docs");

  return (
    <>
      {/* Desktop sidebar */}
      <aside className="hidden w-64 flex-shrink-0 lg:block">
        <div className="sticky top-20 h-[calc(100vh-6rem)] overflow-hidden rounded-lg border bg-card">
          <SidebarContent t={t} />
        </div>
      </aside>

      {/* Mobile sidebar (Sheet) */}
      <div className="lg:hidden">
        <Sheet>
          <SheetTrigger asChild>
            <Button
              className="fixed bottom-20 right-4 z-40 h-12 w-12 rounded-full shadow-lg"
              size="icon"
            >
              <Menu className="size-5" />
              <span className="sr-only">{t("sidebar.openMenu")}</span>
            </Button>
          </SheetTrigger>
          <SheetContent className="w-72 p-0" side="left">
            <SheetHeader className="sr-only">
              <SheetTitle>{t("sidebar.title")}</SheetTitle>
              <SheetDescription>{t("sidebar.description")}</SheetDescription>
            </SheetHeader>
            <SidebarContent t={t} />
          </SheetContent>
        </Sheet>
      </div>
    </>
  );
};
