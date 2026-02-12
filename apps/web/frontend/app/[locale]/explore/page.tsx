import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { ExploreContent } from "@/features/explore";

export const dynamic = "force-dynamic";

type ExplorePageProps = {
  params: Promise<{
    locale: string;
  }>;
};

const ExplorePage = async ({ params }: ExplorePageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations("explore");

  return (
    <main className="container mx-auto px-4 py-8">
      <h1 className="mb-6 text-2xl font-bold tracking-tight">{t("title")}</h1>
      <ExploreContent />
    </main>
  );
};

export default ExplorePage;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "explore" });

  return {
    title: t("title"),
  };
};
