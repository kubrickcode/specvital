import { debounce } from "es-toolkit";
import { useEffect, useRef, useState } from "react";

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

  const debouncedValidateRef = useRef(
    debounce((val: string) => {
      setValidationState(validate(val) ? "valid" : "invalid");
    }, delay)
  );

  useEffect(() => {
    debouncedValidateRef.current = debounce((val: string) => {
      setValidationState(validate(val) ? "valid" : "invalid");
    }, delay);
  }, [validate, delay]);

  useEffect(() => {
    if (value.length < minLength) {
      setValidationState("idle");
      debouncedValidateRef.current.cancel();
      return;
    }

    debouncedValidateRef.current(value);

    return () => debouncedValidateRef.current.cancel();
  }, [value, minLength]);

  return validationState;
};
