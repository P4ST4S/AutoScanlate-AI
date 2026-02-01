import { StatusList } from "@/components/features/StatusList";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Suspense } from "react";

export default function StatusPage() {
  return (
    <div className="min-h-screen p-8 bg-background font-sans">
      <div className="max-w-5xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link href="/">
              <Button variant="ghost" className="w-10 h-10 p-0">
                <ArrowLeft className="w-6 h-6" />
              </Button>
            </Link>
            <h1 className="text-4xl font-display uppercase tracking-tight">
              Translation Status
            </h1>
          </div>
        </header>

        {/* List */}
        <main>
          <Suspense fallback={<div>Loading...</div>}>
            <StatusList />
          </Suspense>
        </main>
      </div>
    </div>
  );
}
