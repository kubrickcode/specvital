import type en from "@/i18n/messages/en.json";

export type Messages = typeof en;

declare global {
  interface IntlMessages extends Messages {}
}
