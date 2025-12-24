import type { Metadata } from "next";
import { Noto_Sans_JP, Potta_One } from "next/font/google";
import "./globals.css";

const notoSansJP = Noto_Sans_JP({
  subsets: ["latin"],
  variable: "--font-noto-sans-jp",
  weight: ["400", "500", "700"],
});

const pottaOne = Potta_One({
  subsets: ["latin"],
  variable: "--font-potta-one",
  weight: "400",
});

export const metadata: Metadata = {
  title: "AutoScanlate AI",
  description: "Automated Manga Translation",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body
        className={`${notoSansJP.variable} ${pottaOne.variable} antialiased bg-background text-foreground font-sans`}
      >
        {children}
      </body>
    </html>
  );
}
