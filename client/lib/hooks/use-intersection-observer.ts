import { useCallback, useRef, useState } from "react";

interface Args extends IntersectionObserverInit {
  freezeOnceVisible?: boolean;
}

interface UseIntersectionObserverReturn {
  entry: IntersectionObserverEntry | undefined;
  setElementRef: (node: Element | null) => void;
}

function useIntersectionObserver({
  threshold = 0,
  root = null,
  rootMargin = "0%",
  freezeOnceVisible = false,
}: Args): UseIntersectionObserverReturn {
  const [entry, setEntry] = useState<IntersectionObserverEntry>();
  const observerRef = useRef<IntersectionObserver | null>(null);

  const frozen = entry?.isIntersecting && freezeOnceVisible;

  const setElementRef = useCallback(
    (node: Element | null) => {
      // 既存のObserverをクリーンアップ
      if (observerRef.current) {
        observerRef.current.disconnect();
        observerRef.current = null;
      }

      const hasIOSupport = !!window.IntersectionObserver;

      if (!hasIOSupport || frozen || !node) return;

      const observerParams = { threshold, root, rootMargin };
      const observer = new IntersectionObserver(([entry]) => {
        setEntry(entry);
      }, observerParams);

      observer.observe(node);
      observerRef.current = observer;
    },
    [threshold, root, rootMargin, frozen]
  );

  return { entry, setElementRef };
}

export default useIntersectionObserver;
