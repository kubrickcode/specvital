"use client";

import { AlertCircle, CheckCircle2, Clock, Loader2, Sparkles, Zap } from "lucide-react";
import { useEffect, useState } from "react";

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Progress } from "@/components/ui/progress";
import { cn } from "@/lib/utils";

import type { SpecGenerationStatusEnum } from "../types";

type GenerationProgressProps = {
  status: SpecGenerationStatusEnum;
};

type Step = {
  description: string;
  icon: typeof Sparkles;
  status: "completed" | "current" | "pending";
  title: string;
};

const GENERATION_STEPS = [
  {
    description: "Preparing test data for analysis",
    icon: Zap,
    title: "Initialization",
  },
  {
    description: "Identifying domains and features",
    icon: Sparkles,
    title: "Domain Classification",
  },
  {
    description: "Converting test names to specifications",
    icon: Sparkles,
    title: "Behavior Conversion",
  },
  {
    description: "Creating executive summary",
    icon: Sparkles,
    title: "Summary Generation",
  },
];

const getStepsForStatus = (status: SpecGenerationStatusEnum): Step[] => {
  const statusToStep: Record<SpecGenerationStatusEnum, number> = {
    completed: 4,
    failed: -1,
    not_found: -1,
    pending: 0,
    running: 1,
  };

  const currentStep = statusToStep[status] ?? -1;

  return GENERATION_STEPS.map((step, index) => ({
    ...step,
    status:
      index < currentStep
        ? ("completed" as const)
        : index === currentStep
          ? ("current" as const)
          : ("pending" as const),
  }));
};

const getProgressPercent = (status: SpecGenerationStatusEnum): number => {
  const statusToProgress: Record<SpecGenerationStatusEnum, number> = {
    completed: 100,
    failed: 0,
    not_found: 0,
    pending: 5,
    running: 40,
  };
  return statusToProgress[status] ?? 0;
};

const StepIcon = ({ step }: { step: Step }) => {
  const Icon = step.icon;

  if (step.status === "completed") {
    return <CheckCircle2 className="h-5 w-5 text-green-500" />;
  }

  if (step.status === "current") {
    return <Loader2 className="h-5 w-5 text-primary animate-spin" />;
  }

  return <Icon className="h-5 w-5 text-muted-foreground" />;
};

export const GenerationProgress = ({ status }: GenerationProgressProps) => {
  const [progress, setProgress] = useState(0);
  const steps = getStepsForStatus(status);
  const targetProgress = getProgressPercent(status);

  useEffect(() => {
    const timer = setTimeout(() => setProgress(targetProgress), 100);
    return () => clearTimeout(timer);
  }, [targetProgress]);

  if (status === "failed") {
    return (
      <Card className="border-destructive/50">
        <CardHeader>
          <div className="flex items-center gap-3">
            <div className="rounded-full bg-destructive/10 p-2">
              <AlertCircle className="h-5 w-5 text-destructive" />
            </div>
            <div>
              <CardTitle className="text-lg">Generation Failed</CardTitle>
              <CardDescription>An error occurred during document generation</CardDescription>
            </div>
          </div>
        </CardHeader>
      </Card>
    );
  }

  const isPending = status === "pending";
  const isRunning = status === "running";

  return (
    <Card aria-live="polite" role="status">
      <CardHeader>
        <div className="flex items-center gap-3">
          <div className="rounded-full bg-primary/10 p-2">
            {isPending ? (
              <Clock className="h-5 w-5 text-primary" />
            ) : (
              <Loader2 className="h-5 w-5 text-primary animate-spin" />
            )}
          </div>
          <div className="flex-1">
            <CardTitle className="text-lg">
              {isPending ? "Queued for Generation" : "Generating Document"}
            </CardTitle>
            <CardDescription>
              {isPending
                ? "Your request is queued and will start processing soon"
                : "AI is analyzing your test cases"}
            </CardDescription>
          </div>
        </div>
      </CardHeader>

      <CardContent className="space-y-6">
        <div className="space-y-2">
          <Progress aria-label={`Generation progress: ${progress}%`} value={progress} />
          <p className="text-xs text-muted-foreground text-right">{progress}%</p>
        </div>

        {isRunning && (
          <div aria-label="Generation steps" className="space-y-3">
            {steps.map((step, index) => (
              <div
                className={cn(
                  "flex items-start gap-3 p-3 rounded-lg transition-colors",
                  step.status === "current" && "bg-primary/5",
                  step.status === "pending" && "opacity-50"
                )}
                key={index}
              >
                <StepIcon step={step} />
                <div className="flex-1 min-w-0">
                  <p
                    className={cn(
                      "text-sm font-medium",
                      step.status === "completed" && "text-muted-foreground"
                    )}
                  >
                    {step.title}
                  </p>
                  <p className="text-xs text-muted-foreground">{step.description}</p>
                </div>
              </div>
            ))}
          </div>
        )}
      </CardContent>
    </Card>
  );
};
