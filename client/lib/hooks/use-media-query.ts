import { useCallback, useRef, useState } from "react";

export default function useMediaQuery() {
  const [device, setDevice] = useState<"mobile" | "tablet" | "desktop" | null>(
    null,
  );
  const [dimensions, setDimensions] = useState<{
    width: number;
    height: number;
  } | null>(null);
  const listenerRef = useRef<(() => void) | null>(null);

  const checkDevice = useCallback(() => {
    if (window.matchMedia("(max-width: 640px)").matches) {
      setDevice("mobile");
    } else if (
      window.matchMedia("(min-width: 641px) and (max-width: 1024px)").matches
    ) {
      setDevice("tablet");
    } else {
      setDevice("desktop");
    }
    setDimensions({ width: window.innerWidth, height: window.innerHeight });
  }, []);

  const setMediaQueryRef = useCallback(
    (node: HTMLElement | null) => {
      if (node && !listenerRef.current) {
        checkDevice();
        window.addEventListener("resize", checkDevice);
        listenerRef.current = () => {
          window.removeEventListener("resize", checkDevice);
        };
      } else if (!node && listenerRef.current) {
        listenerRef.current();
        listenerRef.current = null;
      }
    },
    [checkDevice]
  );

  return {
    device,
    width: dimensions?.width,
    height: dimensions?.height,
    isMobile: device === "mobile",
    isTablet: device === "tablet",
    isDesktop: device === "desktop",
    setMediaQueryRef,
  };
}
