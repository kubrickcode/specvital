"use client";

import { useState } from "react";

import { SupportedFrameworksDialog } from "./supported-frameworks-dialog";
import { TrustBadges } from "./trust-badges";

export const TrustBadgesWithDialog = () => {
  const [isDialogOpen, setIsDialogOpen] = useState(false);

  return (
    <>
      <TrustBadges onFrameworkClick={() => setIsDialogOpen(true)} />
      <SupportedFrameworksDialog onOpenChange={setIsDialogOpen} open={isDialogOpen} />
    </>
  );
};
