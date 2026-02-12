import type { Metadata } from "next";
import { notFound } from "next/navigation";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { isValidGitHubUrl } from "@/features/home";
import { locales } from "@/i18n/config";

import { AnalysisPage } from "./analysis-page";

export const dynamic = "force-dynamic";

type AnalyzePageProps = {
  params: Promise<{
    locale: string;
    owner: string;
    repo: string;
  }>;
};

const AnalyzePage = async ({ params }: AnalyzePageProps) => {
  const { locale, owner, repo } = await params;
  setRequestLocale(locale);

  const repoUrl = `https://github.com/${owner}/${repo}`;
  if (!isValidGitHubUrl(repoUrl)) {
    notFound();
  }

  return <AnalysisPage owner={owner} repo={repo} />;
};

export default AnalyzePage;

export const generateMetadata = async ({ params }: AnalyzePageProps): Promise<Metadata> => {
  const { locale, owner, repo } = await params;
  const t = await getTranslations({ locale, namespace: "metadata" });

  return {
    alternates: {
      canonical: `/${locale}/analyze/${owner}/${repo}`,
      languages: Object.fromEntries(
        locales.map((loc) => [loc, `/${loc}/analyze/${owner}/${repo}`])
      ),
    },
    description: t("analyzeDescription", { owner, repo }),
    title: t("analyzeTitle", { owner, repo }),
  };
};
