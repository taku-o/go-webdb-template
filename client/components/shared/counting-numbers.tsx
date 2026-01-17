"use client";

import { useEffect, useRef, useState } from "react";

export default function CountingNumbers({
  value,
  className,
  start = 0,
  duration = 800,
}: {
  value: number;
  className: string;
  start?: number;
  duration?: number;
}) {
  const [count, setCount] = useState(start);
  const animationRef = useRef<number | null>(null);

  useEffect(() => {
    if (animationRef.current) {
      cancelAnimationFrame(animationRef.current);
    }

    let startTime: number | undefined;
    const animateCount = (timestamp: number) => {
      if (!startTime) startTime = timestamp;
      const timePassed = timestamp - startTime;
      const progress = timePassed / duration;
      const currentCount = easeOutQuad(progress, start, value, 1);
      if (currentCount >= value) {
        setCount(value);
        return;
      }
      setCount(currentCount);
      animationRef.current = requestAnimationFrame(animateCount);
    };
    animationRef.current = requestAnimationFrame(animateCount);

    return () => {
      if (animationRef.current) {
        cancelAnimationFrame(animationRef.current);
      }
    };
  }, [value, start, duration]);

  return <p className={className}>{Intl.NumberFormat().format(count)}</p>;
}

const easeOutQuad = (t: number, b: number, c: number, d: number) => {
  t = t > d ? d : t / d;
  return Math.round(-c * t * (t - 2) + b);
};
