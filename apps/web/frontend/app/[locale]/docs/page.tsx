import { redirect } from "next/navigation";

import { routing } from "@/i18n/routing";

export const dynamic = "force-static";

type DocsPageProps = {
  params: Promise<{
    locale: string;
  }>;
};

export const generateStaticParams = () => {
  return routing.locales.map((locale) => ({ locale }));
};

const DocsPage = async ({ params }: DocsPageProps) => {
  const { locale } = await params;
  redirect(`/${locale}/docs/test-detection`);
};

export default DocsPage;
