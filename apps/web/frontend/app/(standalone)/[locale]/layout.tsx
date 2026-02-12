import { Geist_Mono } from "next/font/google";
import localFont from "next/font/local";
import { notFound } from "next/navigation";
import { NextIntlClientProvider } from "next-intl";
import { getMessages, setRequestLocale } from "next-intl/server";

import { MotionProvider } from "@/components/feedback";
import { ThemeProvider } from "@/components/theme";
import { isValidLocale, locales } from "@/i18n/config";
import "../../globals.css";

const pretendard = localFont({
  display: "swap",
  src: "../../../public/fonts/PretendardVariable.woff2",
  variable: "--font-pretendard",
  weight: "45 920",
});

const geistMono = Geist_Mono({
  subsets: ["latin"],
  variable: "--font-geist-mono",
});

type StandaloneLayoutProps = {
  children: React.ReactNode;
  params: Promise<{
    locale: string;
  }>;
};

const StandaloneLayout = async ({ children, params }: StandaloneLayoutProps) => {
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
        <NextIntlClientProvider messages={messages}>
          <ThemeProvider>
            <MotionProvider>{children}</MotionProvider>
          </ThemeProvider>
        </NextIntlClientProvider>
      </body>
    </html>
  );
};

export default StandaloneLayout;

export const generateStaticParams = () => {
  return locales.map((locale) => ({ locale }));
};
