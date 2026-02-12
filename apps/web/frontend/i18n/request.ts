import { getRequestConfig } from "next-intl/server";

import { isValidLocale } from "./config";
import { routing } from "./routing";

export default getRequestConfig(async ({ requestLocale }) => {
  let locale = await requestLocale;

  if (!locale || !isValidLocale(locale)) {
    locale = routing.defaultLocale;
  }

  return {
    locale,
    messages: (await import(`@/i18n/messages/${locale}.json`)).default,
  };
});
