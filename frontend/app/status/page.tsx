import { StatusList } from "@/components/features/StatusList";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { Button } from "@/components/ui/button";

export default function StatusPage() {
  return (
    <div className="min-h-screen p-8 bg-background font-sans">
      <div className="max-w-5xl mx-auto space-y-8">
        {/* Header */}
        <header className="flex items-center justify-between">
          <div className="flex items-center gap-4">
            <Link href="/">
              <Button variant="ghost" size="icon">
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
          <StatusList />
        </main>
      </div>
    </div>
  );
}
