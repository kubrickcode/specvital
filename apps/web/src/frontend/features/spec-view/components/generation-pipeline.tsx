"use client";

import { Check, X } from "lucide-react";
import { m } from "motion/react";
import { useTranslations } from "next-intl";
import { Fragment } from "react";

import { PulseRing } from "@/components/feedback/pulse-ring";
import { ShimmerBar } from "@/components/feedback/shimmer-bar";
import { useReducedMotion } from "@/lib/motion";
import { cn } from "@/lib/utils";

import type { SpecGenerationStatusEnum } from "../types";

type PipelineStep = {
  label: string;
  number: number;
  status: "completed" | "active" | "failed" | "upcoming";
};

type GenerationPipelineProps = {
  className?: string;
  status: SpecGenerationStatusEnum | null;
};

const getSteps = (
  status: SpecGenerationStatusEnum | null,
  t: (key: string) => string
): PipelineStep[] => {
  if (!status) {
    return [
      { label: t("pipeline.step1"), number: 1, status: "upcoming" },
      { label: t("pipeline.step2"), number: 2, status: "upcoming" },
      { label: t("pipeline.step3"), number: 3, status: "upcoming" },
    ];
  }

  if (status === "pending" || status === "not_found") {
    return [
      { label: t("pipeline.step1"), number: 1, status: "active" },
      { label: t("pipeline.step2"), number: 2, status: "upcoming" },
      { label: t("pipeline.step3"), number: 3, status: "upcoming" },
    ];
  }

  if (status === "running") {
    return [
      { label: t("pipeline.step1"), number: 1, status: "completed" },
      { label: t("pipeline.step2"), number: 2, status: "active" },
      { label: t("pipeline.step3"), number: 3, status: "upcoming" },
    ];
  }

  if (status === "completed") {
    return [
      { label: t("pipeline.step1"), number: 1, status: "completed" },
      { label: t("pipeline.step2"), number: 2, status: "completed" },
      { label: t("pipeline.step3"), number: 3, status: "completed" },
    ];
  }

  if (status === "failed") {
    return [
      { label: t("pipeline.step1"), number: 1, status: "completed" },
      { label: t("pipeline.step2"), number: 2, status: "failed" },
      { label: t("pipeline.step3"), number: 3, status: "upcoming" },
    ];
  }

  return [
    { label: t("pipeline.step1"), number: 1, status: "upcoming" },
    { label: t("pipeline.step2"), number: 2, status: "upcoming" },
    { label: t("pipeline.step3"), number: 3, status: "upcoming" },
  ];
};

type PipelineNodeIconProps = {
  isCurrent: boolean;
  step: PipelineStep;
};

const PipelineNodeIcon = ({ isCurrent, step }: PipelineNodeIconProps) => {
  const shouldReduceMotion = useReducedMotion();

  const nodeContent = (() => {
    if (step.status === "completed") {
      return (
        <div className="flex size-10 items-center justify-center rounded-full bg-primary text-primary-foreground">
          <Check className="size-5" />
        </div>
      );
    }

    if (step.status === "failed") {
      return (
        <div className="flex size-10 items-center justify-center rounded-full border-2 border-destructive bg-destructive/10 text-destructive">
          <X className="size-5" />
        </div>
      );
    }

    if (step.status === "active") {
      return (
        <div className="relative flex size-10 items-center justify-center">
          <div className="absolute inset-0 flex items-center justify-center">
            <PulseRing className="text-[hsl(var(--ai-primary))]" size="md" />
          </div>
          <div className="flex size-10 items-center justify-center rounded-full bg-[hsl(var(--ai-primary))] text-white">
            <span className="text-sm font-semibold">{step.number}</span>
          </div>
        </div>
      );
    }

    return (
      <div className="flex size-10 items-center justify-center rounded-full bg-muted text-muted-foreground">
        <span className="text-sm font-medium">{step.number}</span>
      </div>
    );
  })();

  return (
    <div aria-current={isCurrent ? "step" : undefined} className="shrink-0" role="listitem">
      {shouldReduceMotion || step.status !== "active" ? (
        nodeContent
      ) : (
        <m.div
          animate={{ opacity: 1, scale: 1 }}
          initial={{ opacity: 0, scale: 0.95 }}
          transition={{ duration: 0.4, ease: "easeOut" }}
        >
          {nodeContent}
        </m.div>
      )}
    </div>
  );
};

type PipelineNodeLabelProps = {
  step: PipelineStep;
};

const PipelineNodeLabel = ({ step }: PipelineNodeLabelProps) => {
  return (
    <span
      className={cn(
        "w-10 text-center text-xs font-medium",
        step.status === "active" && "text-foreground",
        step.status === "completed" && "text-muted-foreground",
        step.status === "failed" && "text-destructive",
        step.status === "upcoming" && "text-muted-foreground"
      )}
    >
      {step.label}
    </span>
  );
};

type PipelineConnectorProps = {
  isActive: boolean;
  isCompleted: boolean;
};

const PipelineConnector = ({ isActive, isCompleted }: PipelineConnectorProps) => {
  if (isCompleted) {
    return <div className="mx-3 h-0.5 flex-1 rounded-full bg-primary" />;
  }

  if (isActive) {
    return (
      <div className="mx-3 h-0.5 flex-1 overflow-hidden rounded-full bg-muted">
        <ShimmerBar color="var(--ai-primary)" duration={2} height="xs" />
      </div>
    );
  }

  return <div className="mx-3 h-0.5 flex-1 rounded-full bg-muted" />;
};

export const GenerationPipeline = ({ className, status }: GenerationPipelineProps) => {
  const t = useTranslations("specView.generationProgress");
  const steps = getSteps(status, t);

  return (
    <div
      aria-label="Spec generation progress steps"
      className={cn("flex w-full flex-col gap-2", className)}
      role="list"
    >
      <div className="flex items-center">
        {steps.map((step, index) => (
          <Fragment key={step.number}>
            <PipelineNodeIcon isCurrent={step.status === "active"} step={step} />
            {index < steps.length - 1 && (
              <PipelineConnector
                isActive={step.status === "active"}
                isCompleted={step.status === "completed"}
              />
            )}
          </Fragment>
        ))}
      </div>
      <div className="flex items-start">
        {steps.map((step, index) => (
          <Fragment key={step.number}>
            <PipelineNodeLabel step={step} />
            {index < steps.length - 1 && <div className="mx-3 flex-1" />}
          </Fragment>
        ))}
      </div>
    </div>
  );
};
