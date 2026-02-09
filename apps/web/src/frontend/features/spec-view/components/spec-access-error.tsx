"use client";

import { Github, Loader2, LogIn, ShieldX } from "lucide-react";
import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { useAuth } from "@/features/auth";
import { Link } from "@/i18n/navigation";
import { ROUTES } from "@/lib/routes";

type SpecAccessErrorProps = {
  type: "unauthorized" | "forbidden";
};

export const SpecAccessError = ({ type }: SpecAccessErrorProps) => {
  const t = useTranslations("specView.accessError");
  const { login, loginPending } = useAuth();

  if (type === "unauthorized") {
    return (
      <Card className="border-dashed">
        <CardHeader className="text-center pb-2">
          <div className="mx-auto mb-4 rounded-full bg-primary/10 p-4">
            <LogIn className="size-8 text-primary" />
          </div>
          <CardTitle>{t("unauthorized.title")}</CardTitle>
          <CardDescription className="max-w-sm mx-auto">
            {t("unauthorized.description")}
          </CardDescription>
        </CardHeader>
        <CardContent className="flex flex-col items-center gap-4 pt-2">
          <Button disabled={loginPending} onClick={login} size="lg" variant="cta">
            {loginPending ? (
              <Loader2 className="mr-2 size-5 animate-spin" />
            ) : (
              <Github className="mr-2 size-5" />
            )}
            {t("unauthorized.loginButton")}
          </Button>
        </CardContent>
      </Card>
    );
  }

  return (
    <Card className="border-dashed">
      <CardHeader className="text-center pb-2">
        <div className="mx-auto mb-4 rounded-full bg-destructive/10 p-4">
          <ShieldX className="size-8 text-destructive" />
        </div>
        <CardTitle>{t("forbidden.title")}</CardTitle>
        <CardDescription className="max-w-sm mx-auto">{t("forbidden.description")}</CardDescription>
      </CardHeader>
      <CardContent className="flex flex-col items-center gap-4 pt-2">
        <Button asChild size="lg">
          <Link href={ROUTES.DASHBOARD}>{t("forbidden.dashboardButton")}</Link>
        </Button>
      </CardContent>
    </Card>
  );
};
