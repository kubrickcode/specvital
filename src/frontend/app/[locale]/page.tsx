import { getTranslations, setRequestLocale } from "next-intl/server";
import { UrlInputForm } from "@/components/url-input-form";

export const dynamic = "force-static";

type HomePageProps = {
  params: Promise<{
    locale: string;
  }>;
};

export const HomePage = async ({ params }: HomePageProps) => {
  const { locale } = await params;
  setRequestLocale(locale);

  const t = await getTranslations("home");

  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="w-full max-w-2xl space-y-8 text-center">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">{t("title")}</h1>
          <p className="text-lg text-muted-foreground">{t("subtitle")}</p>
        </div>

        <UrlInputForm />

        <p className="text-sm text-muted-foreground">{t("helperText")}</p>
      </div>
    </main>
  );
};

export default HomePage;
