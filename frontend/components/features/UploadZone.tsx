"use client";

import { useState, useCallback } from "react";
import { useDropzone } from "react-dropzone";
import { motion, AnimatePresence } from "framer-motion";
import { Upload, FileArchive, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";

export function UploadZone() {
  const [files, setFiles] = useState<File[]>([]);

  const onDrop = useCallback((acceptedFiles: File[]) => {
    // Only accept zip files, but for now we accept all for demo
    setFiles((prev) => [...prev, ...acceptedFiles]);
  }, []);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: {
      "application/zip": [".zip"],
      "application/x-zip-compressed": [".zip"],
    },
    maxFiles: 1,
  });

  const removeFile = (fileToRemove: File) => {
    setFiles((prev) => prev.filter((f) => f !== fileToRemove));
  };

  return (
    <div className="w-full max-w-2xl mx-auto space-y-8">
      <div
        {...getRootProps()}
        className={cn(
          "relative group cursor-pointer transition-all duration-300",
          "border-4 border-dashed border-border bg-background p-10 rounded-none",
          "hover:bg-muted/50 hover:border-accent",
          isDragActive && "bg-accent/10 border-accent scale-[1.02]"
        )}
      >
        <input {...getInputProps()} />
        <div className="flex flex-col items-center justify-center text-center space-y-4">
          <div className="relative">
            <div className="absolute inset-0 bg-accent/20 rounded-full blur-xl transform group-hover:scale-150 transition-transform duration-500" />
            <Upload className="relative w-16 h-16 text-foreground group-hover:text-accent transition-colors duration-300" />
          </div>

          <div className="space-y-2">
            <h3 className="text-2xl font-display uppercase tracking-widest">
              {isDragActive ? "Drop it here!" : "Upload Manga"}
            </h3>
            <p className="text-muted-foreground font-medium">
              Drag & drop your .zip file here, or click to select
            </p>
          </div>
        </div>

        {/* Screentone overlay for texture */}
        <div className="absolute inset-0 pointer-events-none screentone-bg z-0" />
      </div>

      <AnimatePresence>
        {files.length > 0 && (
          <motion.div
            initial={{ opacity: 0, y: 20 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -20 }}
          >
            <Card>
              <CardContent className="flex items-center justify-between p-4">
                <div className="flex items-center space-x-4">
                  <div className="p-2 bg-muted border-2 border-border">
                    <FileArchive className="w-6 h-6" />
                  </div>
                  <div>
                    <p className="font-bold text-lg">{files[0].name}</p>
                    <p className="text-sm text-muted-foreground">
                      {(files[0].size / 1024 / 1024).toFixed(2)} MB
                    </p>
                  </div>
                </div>
                <div className="flex space-x-2">
                  <Button
                    variant="ghost"
                    onClick={(e) => {
                      e.stopPropagation();
                      removeFile(files[0]);
                    }}
                  >
                    <X className="w-5 h-5" />
                  </Button>
                  <Button
                    onClick={(e) => {
                      e.stopPropagation();
                      alert("Upload logic to be implemented");
                    }}
                  >
                    Start Translate
                  </Button>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
