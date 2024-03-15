import { Inter } from "next/font/google";
import "./globals.css";

const inter = Inter({ subsets: ["latin"] });

export const metadata = {
  title: "Latest Apple App Reviews",
  description: "See the app reviews for the last 48 hours",
};

export default function RootLayout({ children }) {
  // Happy Pi Day!
  return (
    <html lang="en">
      <body className={inter.className}>{children}</body>
    </html>
  );
}
