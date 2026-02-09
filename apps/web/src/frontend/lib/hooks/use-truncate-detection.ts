import { useEffect, useRef, useState } from "react";

export const useTruncateDetection = <T extends HTMLElement>() => {
  const ref = useRef<T>(null);
  const [isTruncated, setIsTruncated] = useState(false);

  useEffect(() => {
    const checkTruncation = () => {
      const element = ref.current;
      if (!element) return;
      const truncated = element.scrollWidth > element.clientWidth;
      setIsTruncated(truncated);
    };

    // Delay check for layout stabilization
    const timeoutId = setTimeout(checkTruncation, 100);

    // Recheck on window resize
    window.addEventListener("resize", checkTruncation);

    return () => {
      clearTimeout(timeoutId);
      window.removeEventListener("resize", checkTruncation);
    };
  }, []);

  return { isTruncated, ref };
};
