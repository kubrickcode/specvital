"use client";

import { FolderGit2, TestTube2 } from "lucide-react";
import { motion } from "motion/react";
import { useTranslations } from "next-intl";
import { useEffect, useState } from "react";

import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import { easeOutTransition, staggerContainer, staggerItem, useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import { useRepositoryStats } from "../hooks";

type AnimatedCounterProps = {
  duration?: number;
  shouldAnimate: boolean;
  value: number;
};

const AnimatedCounter = ({ duration = 1000, shouldAnimate, value }: AnimatedCounterProps) => {
  const [displayValue, setDisplayValue] = useState(shouldAnimate ? 0 : value);
  const shouldReduceMotion = useReducedMotion();

  useEffect(() => {
    if (!shouldAnimate || shouldReduceMotion) {
      setDisplayValue(value);
      return;
    }

    let startTime: number;
    let animationFrame: number;

    const animate = (currentTime: number) => {
      if (!startTime) startTime = currentTime;
      const progress = Math.min((currentTime - startTime) / duration, 1);
      const easeOut = 1 - Math.pow(1 - progress, 3);
      setDisplayValue(Math.floor(easeOut * value));

      if (progress < 1) {
        animationFrame = requestAnimationFrame(animate);
      }
    };

    animationFrame = requestAnimationFrame(animate);

    return () => {
      if (animationFrame) {
        cancelAnimationFrame(animationFrame);
      }
    };
  }, [duration, shouldAnimate, shouldReduceMotion, value]);

  return <span>{displayValue.toLocaleString()}</span>;
};

type StatCardProps = {
  icon: React.ReactNode;
  isVisible: boolean;
  label: string;
  value: number;
};

type IconVariant = "blue" | "emerald";

type StatCardPropsWithVariant = StatCardProps & {
  iconVariant: IconVariant;
};

const iconVariantStyles: Record<IconVariant, string> = {
  blue: "bg-gradient-to-br from-blue-500/20 to-blue-500/5 text-blue-500 ring-1 ring-blue-500/10",
  emerald:
    "bg-gradient-to-br from-emerald-500/20 to-emerald-500/5 text-emerald-500 ring-1 ring-emerald-500/10",
};

const StatCard = ({ icon, iconVariant, isVisible, label, value }: StatCardPropsWithVariant) => (
  <motion.div variants={staggerItem}>
    <Card className="h-full" depth="raised">
      <CardContent className="flex items-center gap-4">
        <div
          className={cn(
            "flex size-12 shrink-0 items-center justify-center rounded-xl",
            iconVariantStyles[iconVariant]
          )}
        >
          {icon}
        </div>
        <div className="flex flex-col">
          <motion.span
            animate={isVisible ? { opacity: 1 } : { opacity: 0 }}
            className="text-3xl font-bold tabular-nums tracking-tight"
            initial={{ opacity: 0 }}
            transition={easeOutTransition}
          >
            <AnimatedCounter shouldAnimate={isVisible} value={value} />
          </motion.span>
          <span className="text-sm text-muted-foreground">{label}</span>
        </div>
      </CardContent>
    </Card>
  </motion.div>
);

const StatCardSkeleton = () => (
  <Card className="h-full" depth="raised">
    <CardContent className="flex items-center gap-4">
      <Skeleton className="size-12 rounded-xl" />
      <div className="flex flex-col gap-2">
        <Skeleton className="h-8 w-16" />
        <Skeleton className="h-4 w-24" />
      </div>
    </CardContent>
  </Card>
);

export const SummarySection = () => {
  const t = useTranslations("dashboard.summary");
  const { data, isLoading } = useRepositoryStats();

  if (isLoading) {
    return (
      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3" role="status">
        <StatCardSkeleton />
        <StatCardSkeleton />
        <span className="sr-only">{t("loading")}</span>
      </div>
    );
  }

  if (!data) {
    return null;
  }

  return (
    <motion.div
      animate="visible"
      className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3"
      initial="hidden"
      variants={staggerContainer}
    >
      <StatCard
        icon={<TestTube2 aria-hidden="true" className="size-6" />}
        iconVariant="emerald"
        isVisible
        label={t("totalTests")}
        value={data.totalTests}
      />
      <StatCard
        icon={<FolderGit2 aria-hidden="true" className="size-6" />}
        iconVariant="blue"
        isVisible
        label={t("activeRepos")}
        value={data.totalRepositories}
      />
    </motion.div>
  );
};
