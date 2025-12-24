"use client";

import { useState } from "react";
import { motion, AnimatePresence } from "framer-motion";
import { ArrowLeft, ArrowRight, Eye, EyeOff } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import type { Result } from "@/lib/mock-data";
import Link from "next/link"; // Correct import for Link inside component if needed, usually passed or used in page

export function ResultViewer({ result }: { result: Result }) {
  const [currentPage, setCurrentPage] = useState(0);
  const [showOriginal, setShowOriginal] = useState(false);

  const totalPages = result.pages.length;
  const currentImage = showOriginal
    ? result.pages[currentPage].original
    : result.pages[currentPage].translated;

  const handlePrev = () => setCurrentPage((p) => Math.max(0, p - 1));
  const handleNext = () =>
    setCurrentPage((p) => Math.min(totalPages - 1, p + 1));

  return (
    <div className="flex flex-col h-[calc(100vh-100px)] w-full max-w-6xl mx-auto gap-4">
      {/* Toolbar */}
      <Card className="p-2 flex items-center justify-between shadow-sm">
        <div className="flex items-center gap-2">
          <span className="font-display text-lg px-2">
            Page {currentPage + 1} / {totalPages}
          </span>
        </div>

        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 bg-muted rounded-full p-1 border border-border">
            <Button
              variant={showOriginal ? "primary" : "ghost"}
              size="sm"
              onClick={() => setShowOriginal(true)}
              className="rounded-full"
            >
              Original
            </Button>
            <Button
              variant={!showOriginal ? "primary" : "ghost"}
              size="sm"
              onClick={() => setShowOriginal(false)}
              className="rounded-full"
            >
              Translated
            </Button>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <Button
            variant="secondary"
            size="icon"
            onClick={handlePrev}
            disabled={currentPage === 0}
          >
            <ArrowLeft className="w-4 h-4" />
          </Button>
          <Button
            variant="secondary"
            size="icon"
            onClick={handleNext}
            disabled={currentPage === totalPages - 1}
          >
            <ArrowRight className="w-4 h-4" />
          </Button>
        </div>
      </Card>

      {/* Image Viewer */}
      <div className="flex-1 relative bg-muted/20 border-2 border-dashed border-border rounded-lg overflow-hidden flex items-center justify-center p-4">
        <AnimatePresence mode="wait">
          <motion.img
            key={currentPage + (showOriginal ? "_orig" : "_trans")}
            src={currentImage}
            alt={`Page ${currentPage + 1}`}
            className="max-h-full max-w-full object-contain shadow-lg border border-border"
            initial={{ opacity: 0, scale: 0.98 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0 }}
            transition={{ duration: 0.2 }}
          />
        </AnimatePresence>
      </div>
    </div>
  );
}
