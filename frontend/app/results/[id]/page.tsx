import { getResults } from "@/lib/api-client";
import { ResultViewer } from "@/components/features/ResultViewer";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { ArrowLeft } from "lucide-react";
import { notFound } from "next/navigation";

interface PageProps {
  params: Promise<{ id: string }>;
}

export default async function ResultPage(props: PageProps) {
  const params = await props.params;
  const id = params.id;

  let result;
  try {
    result = await getResults(id);
  } catch (error) {
    console.error("Failed to fetch results:", error);
    return notFound();
  }

  if (!result || result.pages.length === 0) {
    return notFound();
  }

  return (
    <div className="min-h-screen p-4 flex flex-col gap-4 bg-background font-sans">
      <header className="flex items-center gap-4">
        <Link href="/status">
          <Button variant="ghost" className="w-10 h-10 p-0">
            <ArrowLeft className="w-6 h-6" />
          </Button>
        </Link>
        <h1 className="text-xl font-display uppercase tracking-tight">
          Reader Mode - {id}
        </h1>
      </header>

      <main className="flex-1">
        <ResultViewer result={result} />
      </main>
    </div>
  );
}
