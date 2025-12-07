import type { Metadata } from "next";
import { Suspense } from "react";
import { notFound } from "next/navigation";
import { getTranslations, setRequestLocale } from "next-intl/server";
import { locales } from "@/i18n/config";
import { isValidGitHubUrl } from "@/features/home";
import { AnalysisContent } from "@/features/analysis";
import { fetchAnalysis } from "@/lib/api";
import Loading from "./loading";

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

  const dataPromise = fetchAnalysis(owner, repo);

  return (
    <Suspense fallback={<Loading />}>
      <AnalysisContent dataPromise={dataPromise} />
    </Suspense>
  );
};

export default AnalyzePage;

export const generateMetadata = async ({ params }: AnalyzePageProps): Promise<Metadata> => {
  const { locale, owner, repo } = await params;
  const t = await getTranslations({ locale, namespace: "metadata" });

  return {
    title: t("analyzeTitle", { owner, repo }),
    description: t("analyzeDescription", { owner, repo }),
    alternates: {
      canonical: `/${locale}/analyze/${owner}/${repo}`,
      languages: Object.fromEntries(
        locales.map((loc) => [loc, `/${loc}/analyze/${owner}/${repo}`])
      ),
    },
  };
};
