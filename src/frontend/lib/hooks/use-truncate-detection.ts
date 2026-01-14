import { useEffect, useRef, useState } from "react";

export const useTruncateDetection = <T extends HTMLElement>() => {
  const ref = useRef<T>(null);
  const [isTruncated, setIsTruncated] = useState(false);

  const checkTruncation = () => {
    const element = ref.current;
    if (element) {
      setIsTruncated(element.scrollWidth > element.clientWidth);
    }
  };

  useEffect(() => {
    checkTruncation();

    window.addEventListener("resize", checkTruncation);
    return () => window.removeEventListener("resize", checkTruncation);
  }, []);

  return { isTruncated, ref };
};
