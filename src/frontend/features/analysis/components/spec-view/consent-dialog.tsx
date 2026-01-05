"use client";

import { useTranslations } from "next-intl";
import { useState } from "react";

import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";

import {
  type ConversionLanguage,
  CONVERSION_LANGUAGES,
  DEFAULT_CONVERSION_LANGUAGE,
  LANGUAGE_LABELS,
} from "../../types";

type ConsentDialogProps = {
  onCancel: () => void;
  onConsent: (language: ConversionLanguage) => void;
  open: boolean;
};

export const ConsentDialog = ({ onCancel, onConsent, open }: ConsentDialogProps) => {
  const t = useTranslations("analyze.specView");
  const [language, setLanguage] = useState<ConversionLanguage>(DEFAULT_CONVERSION_LANGUAGE);

  const handleConvert = () => {
    onConsent(language);
  };

  return (
    <Dialog onOpenChange={(isOpen) => !isOpen && onCancel()} open={open}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle>{t("consentTitle")}</DialogTitle>
          <DialogDescription>{t("consentDescription")}</DialogDescription>
        </DialogHeader>

        <div className="py-4">
          <div className="space-y-2">
            <Label htmlFor="language-select">{t("languageLabel")}</Label>
            <Select
              onValueChange={(value) => setLanguage(value as ConversionLanguage)}
              value={language}
            >
              <SelectTrigger className="w-full" id="language-select">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {CONVERSION_LANGUAGES.map((lang) => (
                  <SelectItem key={lang} value={lang}>
                    {LANGUAGE_LABELS[lang]}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
        </div>

        <DialogFooter>
          <Button onClick={onCancel} variant="outline">
            {t("cancel")}
          </Button>
          <Button onClick={handleConvert}>{t("convert")}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};
