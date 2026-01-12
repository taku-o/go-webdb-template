import "./globals.css";
import cx from "classnames";
import { sfPro, inter } from "./fonts";
import { Suspense } from "react";
import { Analytics as VercelAnalytics } from "@vercel/analytics/react";
import Navbar from "@/components/layout/navbar";

export const metadata = {
  title: "Client2 App",
  description: "New client application",
  metadataBase: new URL("https://precedent.dev"),
};

export default async function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="ja">
      <body className={cx(sfPro.variable, inter.variable)}>
        <div className="pointer-events-none fixed h-screen w-full bg-gradient-to-br from-indigo-50 via-white to-cyan-100 -z-10" />
        <Suspense fallback="...">
          <Navbar />
        </Suspense>
        <main className="flex min-h-screen w-full flex-col items-center justify-center py-32">
          {children}
        </main>
        <VercelAnalytics />
      </body>
    </html>
  );
}
