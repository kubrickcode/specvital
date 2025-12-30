import type { Transition, Variants } from "motion/react";

export const defaultTransition: Transition = {
  damping: 30,
  stiffness: 300,
  type: "spring",
};

export const easeOutTransition: Transition = {
  duration: 0.3,
  ease: "easeOut",
  type: "tween",
};

export const fadeIn: Variants = {
  hidden: { opacity: 0 },
  visible: { opacity: 1 },
};

export const fadeInUp: Variants = {
  hidden: { opacity: 0, y: 20 },
  visible: {
    opacity: 1,
    transition: defaultTransition,
    y: 0,
  },
};

export const fadeInDown: Variants = {
  hidden: { opacity: 0, y: -20 },
  visible: {
    opacity: 1,
    transition: defaultTransition,
    y: 0,
  },
};

export const scaleIn: Variants = {
  hidden: { opacity: 0, scale: 0.95 },
  visible: {
    opacity: 1,
    scale: 1,
    transition: defaultTransition,
  },
};

export const staggerContainer: Variants = {
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      delayChildren: 0.1,
      staggerChildren: 0.1,
    },
  },
};

export const createStaggerContainer = (staggerDuration = 0.1, delayDuration = 0.1): Variants => ({
  hidden: { opacity: 0 },
  visible: {
    opacity: 1,
    transition: {
      delayChildren: delayDuration,
      staggerChildren: staggerDuration,
    },
  },
});

export const staggerItem: Variants = {
  hidden: { opacity: 0, y: 10 },
  visible: {
    opacity: 1,
    transition: easeOutTransition,
    y: 0,
  },
};
