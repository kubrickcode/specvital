import type { Metadata, Viewport } from "next";
import { Geist_Mono } from "next/font/google";
import localFont from "next/font/local";
import { notFound } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, getTranslations, setRequestLocale } from "next-intl/server";
import { NuqsAdapter } from "nuqs/adapters/next/app";

import { Header, MobileBottomBar } from "@/components/layout";
import { ThemeProvider } from "@/components/theme";
import { Toaster } from "@/components/ui/sonner";
import { isValidLocale, locales } from "@/i18n/config";
import { QueryProvider } from "@/lib/query";
import "../globals.css";

const pretendard = localFont({
  display: "swap",
  src: "../../public/fonts/PretendardVariable.woff2",
  variable: "--font-pretendard",
  weight: "45 920",
});

const geistMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-geist-mono",
});

type LocaleLayoutProps = {
  children: React.ReactNode;
  params: Promise<{
    locale: string;
  }>;
};

const LocaleLayout = async ({ children, params }: LocaleLayoutProps) => {
  const { locale } = await params;

  if (!isValidLocale(locale)) {
    notFound();
  }

  setRequestLocale(locale);

  const messages = await getMessages();

  return (
    <html lang={locale} suppressHydrationWarning>
      <body
        className={`${pretendard.variable} ${geistMono.variable} flex min-h-dvh flex-col font-sans antialiased`}
      >
        <NuqsAdapter>
          <QueryProvider>
            <NextIntlClientProvider messages={messages}>
              <ThemeProvider>
                <Header />
                <main className="flex flex-1 flex-col pb-16 md:pb-0" id="main-content">
                  {children}
                </main>
                <MobileBottomBar />
                <Toaster richColors />
              </ThemeProvider>
            </NextIntlClientProvider>
          </QueryProvider>
        </NuqsAdapter>
      </body>
    </html>
  );
};

export default LocaleLayout;

export const generateMetadata = async ({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> => {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: "metadata" });

  return {
    description: t("defaultDescription"),
    metadataBase: new URL(process.env.NEXT_PUBLIC_SITE_URL || "http://localhost:5173"),
    title: t("defaultTitle"),
  };
};

export const generateStaticParams = () => {
  return locales.map((locale) => ({ locale }));
};

export const viewport: Viewport = {
  viewportFit: "cover",
};
