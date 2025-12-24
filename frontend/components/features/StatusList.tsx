"use client";

import { motion } from "framer-motion";
import {
  CheckCircle2,
  Clock,
  Loader2,
  AlertCircle,
  ArrowRight,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { mockRequests, type Request } from "@/lib/mock-data";
import Link from "next/link";

interface StatusConfigItem {
  icon: React.ElementType;
  color: string;
  label: string;
  bg: string;
  animate?: boolean;
}

const statusConfig: Record<string, StatusConfigItem> = {
  queued: {
    icon: Clock,
    color: "text-muted-foreground",
    label: "In Queue",
    bg: "bg-muted",
  },
  processing: {
    icon: Loader2,
    color: "text-accent",
    label: "Translating...",
    bg: "bg-accent/10 border-accent",
    animate: true,
  },
  completed: {
    icon: CheckCircle2,
    color: "text-green-600",
    label: "Completed",
    bg: "bg-green-100 border-green-600",
  },
  failed: {
    icon: AlertCircle,
    color: "text-red-600",
    label: "Failed",
    bg: "bg-red-100 border-red-600",
  },
};

export function StatusList() {
  return (
    <div className="w-full max-w-4xl mx-auto space-y-4">
      {mockRequests.map((req, index) => {
        const config = statusConfig[req.status];
        const Icon = config.icon;

        return (
          <motion.div
            key={req.id}
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ delay: index * 0.1 }}
          >
            <Card
              className={cn(
                "group hover:shadow-[6px_6px_0px_var(--border)] transition-all duration-200",
                req.status === "processing" &&
                  "border-accent shadow-[4px_4px_0px_var(--accent)]"
              )}
            >
              <div className="flex items-center p-4 gap-4">
                {/* Thumbnail / Status Icon */}
                <div
                  className={cn(
                    "w-16 h-24 flex items-center justify-center border-2 border-border shrink-0 bg-background relative overflow-hidden",
                    config.bg
                  )}
                >
                  {req.thumbnail ? (
                    <img
                      src={req.thumbnail}
                      alt={req.filename}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <Icon
                      className={cn(
                        "w-8 h-8",
                        config.color,
                        config.animate && "animate-spin"
                      )}
                    />
                  )}
                </div>

                {/* Info */}
                <div className="flex-1 min-w-0">
                  <div className="flex items-center justify-between mb-1">
                    <h3 className="font-display text-lg truncate pr-4">
                      {req.filename}
                    </h3>
                    <span
                      className={cn(
                        "text-xs font-bold uppercase tracking-wider px-2 py-1 border border-current rounded-full",
                        config.color
                      )}
                    >
                      {config.label}
                    </span>
                  </div>

                  <p className="text-sm text-muted-foreground mb-3">
                    {req.pageCount} pages •{" "}
                    {new Date(req.createdAt).toLocaleDateString()}
                  </p>

                  {/* Progress Bar for processing */}
                  {req.status === "processing" && (
                    <div className="w-full h-3 bg-muted border border-border rounded-full overflow-hidden relative">
                      <motion.div
                        className="h-full bg-accent"
                        initial={{ width: 0 }}
                        animate={{ width: `${req.progress}%` }}
                        transition={{ duration: 1, ease: "easeOut" }}
                      />
                      {/* Striped pattern overlay */}
                      <div className="absolute inset-0 opacity-20 bg-[linear-gradient(45deg,transparent_25%,#000_25%,#000_50%,transparent_50%,transparent_75%,#000_75%,#000_100%)] bg-[length:10px_10px]" />
                    </div>
                  )}
                </div>

                {/* Actions */}
                <div className="shrink-0">
                  {req.status === "completed" && (
                    <Link href={`/results/${req.id}`}>
                      <Button className="group-hover:translate-x-1 transition-transform">
                        Read <ArrowRight className="ml-2 w-4 h-4" />
                      </Button>
                    </Link>
                  )}
                </div>
              </div>
            </Card>
          </motion.div>
        );
      })}
    </div>
  );
}
