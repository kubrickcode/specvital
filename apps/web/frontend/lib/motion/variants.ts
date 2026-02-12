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

export const slideInUp: Variants = {
  hidden: { opacity: 0, y: "100%" },
  visible: {
    opacity: 1,
    transition: easeOutTransition,
    y: 0,
  },
};

export const staggerItem: Variants = {
  hidden: { opacity: 0, y: 10 },
  visible: {
    opacity: 1,
    transition: easeOutTransition,
    y: 0,
  },
};

export const collapseTransition: Transition = {
  duration: 0.2,
  ease: [0.4, 0, 0.2, 1],
  type: "tween",
};

export const expandCollapse: Variants = {
  collapsed: {
    height: 0,
    opacity: 0,
    overflow: "hidden",
  },
  expanded: {
    height: "auto",
    opacity: 1,
    overflow: "visible",
    transition: collapseTransition,
  },
};
