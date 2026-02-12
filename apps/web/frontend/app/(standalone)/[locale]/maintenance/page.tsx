import { setRequestLocale } from "next-intl/server";

import { MaintenancePage } from "@/components/feedback";

export const dynamic = "force-static";

type MaintenancePageRouteProps = {
  params: Promise<{
    locale: string;
  }>;
};

export default async function MaintenancePageRoute({ params }: MaintenancePageRouteProps) {
  const { locale } = await params;
  setRequestLocale(locale);

  return <MaintenancePage />;
}
