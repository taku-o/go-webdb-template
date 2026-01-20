import Image from "next/image";
import Link from "next/link";
import { auth } from "@/auth";
import { AuthButtons } from "@/components/auth/auth-buttons";
import NavBarClient from "./navbar-client";

export default async function NavBar() {
  const session = await auth();

  return (
    <>
      <div className="fixed top-0 flex w-full justify-center z-30 transition-all">
        <NavBarClient>
          <div className="mx-5 flex h-16 w-full max-w-screen-xl items-center justify-between">
            <Link href="/" className="flex items-center font-display text-2xl">
              <Image
                src="/logo.png"
                alt="Precedent logo"
                width="30"
                height="30"
                className="mr-2 rounded-sm"
              ></Image>
              <p>Precedent</p>
            </Link>
            <div className="flex items-center gap-4">
              <AuthButtons user={session?.user ?? null} />
            </div>
          </div>
        </NavBarClient>
      </div>
    </>
  );
}
