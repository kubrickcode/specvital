"use client";

import { ChevronDown } from "lucide-react";
import { AnimatePresence, motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useState } from "react";

import { cn } from "@/lib/utils";

import { FAQ_ITEMS } from "../constants/plans";

type FaqItemProps = {
  answer: string;
  isOpen: boolean;
  onToggle: () => void;
  question: string;
};

const FaqItem = ({ answer, isOpen, onToggle, question }: FaqItemProps) => {
  return (
    <div className="border-b last:border-b-0">
      <button
        aria-expanded={isOpen}
        className={cn(
          "flex w-full items-center justify-between py-4 text-left",
          "transition-colors hover:text-foreground",
          isOpen ? "text-foreground" : "text-muted-foreground"
        )}
        onClick={onToggle}
        type="button"
      >
        <span className="pr-4 font-medium">{question}</span>
        <ChevronDown
          className={cn(
            "size-5 shrink-0 transition-transform duration-200",
            isOpen && "rotate-180"
          )}
        />
      </button>
      <AnimatePresence initial={false}>
        {isOpen && (
          <motion.div
            animate={{ height: "auto", opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            initial={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
          >
            <p className="pb-4 text-sm text-muted-foreground">{answer}</p>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export const PricingFaq = () => {
  const t = useTranslations("pricing.faq");
  const [openIndex, setOpenIndex] = useState<number | null>(null);

  const toggleItem = (index: number) => {
    setOpenIndex((prev) => (prev === index ? null : index));
  };

  return (
    <section aria-labelledby="faq-heading" className="mt-16">
      <h2 className="mb-8 text-center text-2xl font-bold" id="faq-heading">
        {t("title")}
      </h2>
      <div className="mx-auto max-w-2xl rounded-lg border bg-card px-4">
        {FAQ_ITEMS.map((key, index) => (
          <FaqItem
            answer={t(`items.${key}.answer`)}
            isOpen={openIndex === index}
            key={key}
            onToggle={() => toggleItem(index)}
            question={t(`items.${key}.question`)}
          />
        ))}
      </div>
    </section>
  );
};
