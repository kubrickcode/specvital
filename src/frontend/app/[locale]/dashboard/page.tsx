import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";
import { Suspense } from "react";

import { LoadingFallback } from "@/components/feedback";
import { RequireAuth } from "@/features/auth";
import { DashboardContent } from "@/features/dashboard";

export const dynamic = "force-dynamic";

type DashboardPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

const DashboardPage = async ({ params }: DashboardPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations("dashboard");

  return (
    <Suspense fallback={<LoadingFallback />}>
      <RequireAuth>
        <main className="container mx-auto px-4 py-8">
          <h1 className="mb-6 text-2xl font-bold tracking-tight">{t("title")}</h1>
          <DashboardContent />
        </main>
      </RequireAuth>
    </Suspense>
  );
};

export default DashboardPage;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "dashboard" });

  return {
    title: t("title"),
  };
};
