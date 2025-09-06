import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "ToDo App",
  description: "A simple todo application built with Next.js and Go",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        {/* AuthProvider will be added here later */}
        {children}
      </body>
    </html>
  );
}