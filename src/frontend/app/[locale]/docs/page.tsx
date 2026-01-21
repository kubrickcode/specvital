import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { DocsLandingContent } from "@/features/docs";

export const dynamic = "force-static";

type DocsPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

const DocsPage = async ({ params }: DocsPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  return (
    <main id="main-content">
      <DocsLandingContent />
    </main>
  );
};

export default DocsPage;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "docs" });

  return {
    description: t("landing.description"),
    title: t("landing.title"),
  };
};
