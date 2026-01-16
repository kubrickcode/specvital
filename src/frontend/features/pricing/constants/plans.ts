import type { PricingPlan } from "../types";

export const PRICING_PLANS: PricingPlan[] = [
  {
    cta: "getStarted",
    features: [
      { label: "specview", value: 5000 },
      { label: "analysis", value: 50 },
      { label: "retention", value: 30 },
    ],
    monthlyPrice: 0,
    tier: "free",
  },
  {
    cta: "startFree",
    features: [
      { label: "specview", value: 100000 },
      { label: "analysis", value: 1000 },
      { label: "retention", value: 180 },
    ],
    highlighted: true,
    monthlyPrice: 15,
    tier: "pro",
  },
  {
    cta: "startFree",
    features: [
      { label: "specview", value: 500000 },
      { label: "analysis", value: 5000 },
      { label: "retention", value: 365 },
    ],
    monthlyPrice: 59,
    tier: "pro_plus",
  },
  {
    cta: "contactUs",
    features: [
      { label: "specview", value: null },
      { label: "analysis", value: null },
      { label: "retention", value: null },
    ],
    monthlyPrice: null,
    tier: "enterprise",
  },
];

export const FAQ_ITEMS = [
  "specview",
  "analysis",
  "promotionFree",
  "paymentLive",
  "retention",
] as const;

export type FaqKey = (typeof FAQ_ITEMS)[number];
