"use client";

import { ReactNode } from "react";
import useScroll from "@/lib/hooks/use-scroll";

interface NavBarClientProps {
  children: ReactNode;
}

export default function NavBarClient({ children }: NavBarClientProps) {
  const scrolled = useScroll(50);

  return (
    <div
      className={`w-full ${
        scrolled
          ? "border-b border-gray-200 bg-white/50 backdrop-blur-xl"
          : "bg-white/0"
      } transition-all`}
    >
      {children}
    </div>
  );
}
