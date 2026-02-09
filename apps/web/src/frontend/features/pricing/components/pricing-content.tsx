"use client";

import { PricingFaq } from "./pricing-faq";
import { PricingGrid } from "./pricing-grid";
import { PricingHeader } from "./pricing-header";
import { PromotionBanner } from "./promotion-banner";

export const PricingContent = () => {
  return (
    <div className="container mx-auto max-w-7xl px-4 py-12">
      <PromotionBanner />
      <PricingHeader />
      <PricingGrid />
      <PricingFaq />
    </div>
  );
};
