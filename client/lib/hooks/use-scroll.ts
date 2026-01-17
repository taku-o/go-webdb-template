import { useCallback, useRef, useState } from "react";

export default function useScroll(threshold: number) {
  const [scrolled, setScrolled] = useState(false);
  const listenerRef = useRef<(() => void) | null>(null);

  const setScrollRef = useCallback(
    (node: HTMLElement | null) => {
      if (node && !listenerRef.current) {
        const onScroll = () => {
          setScrolled(window.pageYOffset > threshold);
        };
        window.addEventListener("scroll", onScroll);
        listenerRef.current = () => {
          window.removeEventListener("scroll", onScroll);
        };
      } else if (!node && listenerRef.current) {
        listenerRef.current();
        listenerRef.current = null;
      }
    },
    [threshold]
  );

  return { scrolled, setScrollRef };
}
