import { Layers } from "lucide-react";
import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { QueueProcessingContent } from "@/features/docs";

export const dynamic = "force-static";

type QueueProcessingPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "docs.queueProcessing" });

  return {
    description: t("meta.description"),
    title: t("meta.title"),
  };
};

const QueueProcessingPage = async ({ params }: QueueProcessingPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations({ locale, namespace: "docs.queueProcessing" });

  return (
    <article>
      <header className="mb-10">
        <div className="mb-4 flex size-12 items-center justify-center rounded-xl bg-primary/10 text-primary">
          <Layers className="size-6" />
        </div>
        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
        <p className="mt-3 text-lg text-muted-foreground">{t("subtitle")}</p>
      </header>

      <QueueProcessingContent />
    </article>
  );
};

export default QueueProcessingPage;
