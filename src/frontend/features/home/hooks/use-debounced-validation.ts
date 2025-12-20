import { debounce } from "es-toolkit";
import { useEffect, useMemo, useState } from "react";

type ValidationState = "idle" | "valid" | "invalid";

type UseDebouncedValidationOptions = {
  delay?: number;
  minLength?: number;
};

export const useDebouncedValidation = (
  value: string,
  validate: (value: string) => boolean,
  options: UseDebouncedValidationOptions = {}
): ValidationState => {
  const { delay = 500, minLength = 1 } = options;
  const [validationState, setValidationState] = useState<ValidationState>("idle");

  const debouncedValidate = useMemo(
    () =>
      debounce((val: string) => {
        setValidationState(validate(val) ? "valid" : "invalid");
      }, delay),
    [validate, delay]
  );

  useEffect(() => {
    if (value.length < minLength) {
      setValidationState("idle");
      debouncedValidate.cancel();
      return;
    }

    debouncedValidate(value);

    return () => debouncedValidate.cancel();
  }, [value, debouncedValidate, minLength]);

  return validationState;
};
