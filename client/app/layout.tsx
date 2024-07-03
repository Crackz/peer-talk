import AuthContextProvider from "@/modules/auth-provider";
import "./globals.css";
import type { Metadata } from "next";
import { Inter } from "next/font/google";
const inter = Inter({ subsets: ["latin"] });
import Image from "next/image";
import Link from "next/link";

export const metadata: Metadata = {
  title: "Peer Talk",
  description: "video chat online",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <AuthContextProvider>
          <header id="nav">
            <div className="nav--list">
              <button id="members__button">
                <svg
                  width="24"
                  height="24"
                  xmlns="http://www.w3.org/2000/svg"
                  fillRule="evenodd"
                  clipRule="evenodd"
                >
                  <path
                    d="M24 18v1h-24v-1h24zm0-6v1h-24v-1h24zm0-6v1h-24v-1h24z"
                    fill="#ede0e0"
                  />
                  <path d="M24 19h-24v-1h24v1zm0-6h-24v-1h24v1zm0-6h-24v-1h24v1z" />
                </svg>
              </button>
              <Link href="/">
                <h3 id="logo">
                  <Image src="/logo.png" alt="Logo" width={150} height={150} />
                </h3>
              </Link>
            </div>

            <div id="nav__links">
              <button id="chat__button">
                <svg
                  width="24"
                  height="24"
                  xmlns="http://www.w3.org/2000/svg"
                  fillRule="evenodd"
                  fill="#ede0e0"
                  clipRule="evenodd"
                >
                  <path d="M24 20h-3v4l-5.333-4h-7.667v-4h2v2h6.333l2.667 2v-2h3v-8.001h-2v-2h4v12.001zm-15.667-6l-5.333 4v-4h-3v-14.001l18 .001v14h-9.667zm-6.333-2h3v2l2.667-2h8.333v-10l-14-.001v10.001z" />
                </svg>
              </button>
              <Link className="nav__link" id="create__room__btn" href="/lobby">
                Create Room
              </Link>
            </div>
          </header>
          {children}
        </AuthContextProvider>
      </body>
    </html>
  );
}
