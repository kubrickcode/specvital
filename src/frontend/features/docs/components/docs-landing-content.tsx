"use client";

import { BookOpen, CreditCard, GitBranch, Layers, Zap } from "lucide-react";
import { useTranslations } from "next-intl";

import { DocsTopicCard } from "./docs-topic-card";

const TOPICS = [
  { icon: Layers, key: "testDetection" },
  { icon: CreditCard, key: "usageBilling" },
  { icon: GitBranch, key: "githubAccess" },
  { icon: Zap, key: "queueProcessing" },
  { icon: BookOpen, key: "specviewGeneration" },
] as const;

export const DocsLandingContent = () => {
  const t = useTranslations("docs");

  return (
    <div className="container mx-auto max-w-5xl px-4 py-12">
      <header className="mb-12 text-center">
        <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">{t("landing.title")}</h1>
        <p className="mx-auto mt-4 max-w-2xl text-muted-foreground">{t("landing.description")}</p>
      </header>

      <section aria-labelledby="how-it-works-heading">
        <h2 className="mb-6 text-xl font-semibold" id="how-it-works-heading">
          {t("landing.howItWorksTitle")}
        </h2>
        <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
          {TOPICS.map(({ icon, key }) => (
            <DocsTopicCard
              description={t(`topics.${key}.description`)}
              href={`/docs/how-it-works/${t(`topics.${key}.slug`)}`}
              icon={icon}
              key={key}
              title={t(`topics.${key}.title`)}
            />
          ))}
        </div>
      </section>
    </div>
  );
};
