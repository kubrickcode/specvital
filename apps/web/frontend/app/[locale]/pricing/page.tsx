import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { PricingContent } from "@/features/pricing";

export const dynamic = "force-static";

type PricingPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

const PricingPage = async ({ params }: PricingPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  return (
    <main id="main-content">
      <PricingContent />
    </main>
  );
};

export default PricingPage;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "pricing" });

  return {
    description: t("header.subtitle"),
    title: t("header.title"),
  };
};
