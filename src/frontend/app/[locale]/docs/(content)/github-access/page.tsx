import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { GithubAccessContent } from "@/features/docs";

export const dynamic = "force-static";

type GithubAccessPageProps = {
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
  const t = await getTranslations({ locale, namespace: "docs.githubAccess" });

  return {
    description: t("meta.description"),
    title: t("meta.title"),
  };
};

const GithubAccessPage = async ({ params }: GithubAccessPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations({ locale, namespace: "docs.githubAccess" });

  return (
    <article>
      <header className="mb-10">
        <h1 className="text-3xl font-bold tracking-tight">{t("title")}</h1>
        <p className="mt-3 text-lg text-muted-foreground">{t("subtitle")}</p>
      </header>

      <GithubAccessContent />
    </article>
  );
};

export default GithubAccessPage;
