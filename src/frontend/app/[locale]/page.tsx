import { getTranslations, setRequestLocale } from "next-intl/server";
import { Suspense } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { OAuthErrorHandler } from "@/features/auth";
import { TrustBadgesWithDialog, UrlInputForm } from "@/features/home";

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
    <main className="flex min-h-screen flex-col items-center justify-center bg-[radial-gradient(ellipse_80%_50%_at_50%_40%,var(--hero-gradient-center),var(--hero-gradient-edge)_70%)] p-8">
      <Suspense fallback={null}>
        <OAuthErrorHandler />
      </Suspense>

      <div className="w-full max-w-2xl space-y-6 text-center">
        <div className="space-y-2">
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">{t("headline")}</h1>
          <p className="text-muted-foreground">{t("subheadline")}</p>
        </div>

        <Card className="mx-auto w-full max-w-xl" depth="elevated">
          <CardContent>
            <UrlInputForm />
          </CardContent>
        </Card>

        <TrustBadgesWithDialog />
      </div>
    </main>
  );
};

export default HomePage;
