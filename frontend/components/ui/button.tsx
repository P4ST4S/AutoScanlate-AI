import { ButtonHTMLAttributes, forwardRef } from "react";
import { cn } from "@/lib/utils";
import { Slot } from "@radix-ui/react-slot";

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  asChild?: boolean;
  variant?: "primary" | "secondary" | "ghost";
}

const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant = "primary", asChild = false, ...props }, ref) => {
    const Comp = asChild ? Slot : "button";
    return (
      <Comp
        className={cn(
          "inline-flex items-center justify-center whitespace-nowrap text-sm font-medium ring-offset-background transition-all hover:scale-105 active:scale-95 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50",
          "h-10 px-6 py-2 border-2 border-border shadow-[4px_4px_0px_0px_var(--border)] hover:shadow-[2px_2px_0px_0px_var(--border)] hover:translate-x-[2px] hover:translate-y-[2px]",
          variant === "primary" && "bg-accent text-white hover:bg-accent/90",
          variant === "secondary" &&
            "bg-background text-foreground hover:bg-muted",
          variant === "ghost" &&
            "border-transparent shadow-none hover:bg-muted hover:shadow-none hover:translate-x-0 hover:translate-y-0",
          className
        )}
        ref={ref}
        {...props}
      />
    );
  }
);
Button.displayName = "Button";

export { Button };
