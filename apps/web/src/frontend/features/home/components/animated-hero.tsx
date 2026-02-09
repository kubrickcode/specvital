"use client";

import { motion } from "motion/react";
import type { ReactNode } from "react";

import { createStaggerContainer, fadeInUp, useReducedMotion } from "@/lib/motion";

type AnimatedHeroProps = {
  card: ReactNode;
  headline: ReactNode;
  subheadline: ReactNode;
  trustBadges: ReactNode;
};

const heroStaggerContainer = createStaggerContainer(0.15, 0);

export const AnimatedHero = ({ card, headline, subheadline, trustBadges }: AnimatedHeroProps) => {
  const shouldReduceMotion = useReducedMotion();

  const containerVariants = shouldReduceMotion ? {} : heroStaggerContainer;
  const itemVariants = shouldReduceMotion ? {} : fadeInUp;

  return (
    <motion.div
      animate="visible"
      className="w-full max-w-2xl space-y-6 text-center"
      initial={shouldReduceMotion ? false : "hidden"}
      variants={containerVariants}
    >
      <motion.div className="space-y-2" variants={itemVariants}>
        {headline}
        {subheadline}
      </motion.div>

      <motion.div variants={itemVariants}>{card}</motion.div>

      <motion.div variants={itemVariants}>{trustBadges}</motion.div>
    </motion.div>
  );
};
