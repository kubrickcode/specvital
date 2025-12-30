import { getTranslations, setRequestLocale } from "next-intl/server";
import { Suspense } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { OAuthErrorHandler } from "@/features/auth";
import { AnimatedHero, TrustBadgesWithDialog, UrlInputForm } from "@/features/home";

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
    <div className="flex flex-1 flex-col items-center justify-center bg-[radial-gradient(ellipse_80%_50%_at_50%_40%,var(--hero-gradient-center),var(--hero-gradient-edge)_70%)] p-8">
      <Suspense fallback={null}>
        <OAuthErrorHandler />
      </Suspense>

      <AnimatedHero
        card={
          <Card className="mx-auto w-full max-w-xl" depth="elevated">
            <CardContent>
              <UrlInputForm />
            </CardContent>
          </Card>
        }
        headline={
          <h1 className="text-3xl font-bold tracking-tight sm:text-4xl">{t("headline")}</h1>
        }
        subheadline={<p className="text-muted-foreground">{t("subheadline")}</p>}
        trustBadges={<TrustBadgesWithDialog />}
      />
    </div>
  );
};

export default HomePage;
