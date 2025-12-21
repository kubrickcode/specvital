import { getTranslations, setRequestLocale } from "next-intl/server";
import { Suspense } from "react";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { OAuthErrorHandler } from "@/features/auth";
import { TrustBadges, UrlInputForm } from "@/features/home";

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
      <Suspense fallback={null}>
        <OAuthErrorHandler />
      </Suspense>

      <div className="w-full max-w-2xl space-y-8 text-center">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">{t("title")}</h1>
          <p className="text-xl font-medium">{t("tagline")}</p>
          <p className="mx-auto max-w-xl text-base text-muted-foreground">{t("subtitle")}</p>
        </div>

        <Card className="mx-auto w-full max-w-xl">
          <CardHeader className="text-center">
            <CardTitle className="text-lg">{t("inputLabel")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <UrlInputForm />
            <p className="text-center text-sm text-muted-foreground">{t("helperText")}</p>
          </CardContent>
        </Card>

        <div className="mx-auto w-full max-w-xl">
          <TrustBadges />
        </div>
      </div>
    </main>
  );
};

export default HomePage;
