import { UploadZone } from "@/components/features/UploadZone";
import Link from "next/link";
import { Button } from "@/components/ui/button";

export default function Home() {
  return (
    <div className="min-h-screen p-8 pb-20 sm:p-20 font-sans relative overflow-hidden">
      {/* Background decoration */}
      <div className="absolute inset-0 z-0 opacity-5 pointer-events-none bg-[radial-gradient(#000_1px,transparent_1px)] [background-size:16px_16px]"></div>

      <main className="flex flex-col gap-12 items-center relative z-10 max-w-5xl mx-auto">
        {/* Header Section */}
        <header className="text-center space-y-4">
          <h1 className="text-6xl md:text-8xl font-display uppercase tracking-tighter text-foreground drop-shadow-[4px_4px_0px_var(--accent)]">
            Auto<span className="text-accent">Scanlate</span>
          </h1>
          <p className="text-xl text-muted-foreground font-medium max-w-xl mx-auto">
            AI-powered manga translation pipeline. Local, Private, Fast.
          </p>
          <div className="pt-4">
            <Link href="/status">
              <Button variant="secondary">View Dashboard</Button>
            </Link>
          </div>
        </header>

        {/* Upload Section */}
        <section className="w-full">
          <UploadZone />
        </section>
      </main>

      <footer className="fixed bottom-0 w-full p-4 text-center text-sm text-muted-foreground bg-background/80 backdrop-blur-sm border-t-2 border-border/10">
        <p>AutoScanlate AI - V10 Worker Compatible</p>
      </footer>
    </div>
  );
}
