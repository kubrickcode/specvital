import type { Metadata } from "next";
import { getTranslations, setRequestLocale } from "next-intl/server";

import { AccountContent } from "@/features/account";

export const dynamic = "force-dynamic";

type AccountPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

const AccountPage = async ({ params }: AccountPageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations("account");

  return (
    <main className="container mx-auto max-w-4xl px-4 py-8">
      <h1 className="mb-6 text-2xl font-bold tracking-tight">{t("title")}</h1>
      <AccountContent />
    </main>
  );
};

export default AccountPage;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "account" });

  return {
    title: t("title"),
  };
};
