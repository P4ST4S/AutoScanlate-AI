# <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/panel-top.svg" width="24" height="24" /> AutoScanlate AI - Frontend

The modern, responsive web interface for the AutoScanlate AI pipeline. Built with Next.js 16 and styled with a custom "Manga/Japanese" aesthetic.

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/layers.svg" width="24" height="24" /> Tech Stack

- **Framework**: [Next.js 16](https://nextjs.org/) (App Router)
- **Language**: [TypeScript](https://www.typescriptlang.org/)
- **Styling**: [TailwindCSS v4](https://tailwindcss.com/)
- **Animations**: [Framer Motion](https://www.framer.com/motion/)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Fonts**: [Next/Font](https://nextjs.org/docs/basic-features/font-optimization) (Google Fonts: Noto Sans JP & Potta One)

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/palette.svg" width="24" height="24" /> Design System (Manga Theme)

The UI follows a strict "Ink & Paper" aesthetic to match the subject matter.

- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/droplet.svg" width="16" height="16" /> Colors**:

  - `bg-paper` (#fdfbf7): Warm white, emulating manga paper.
  - `text-ink` (#0f0f0f): Deep black for text and borders.
  - `accent` (#e63946): Stamp red for highlights and actions.

- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/type.svg" width="16" height="16" /> Typography**:

  - **Headers**: _Potta One_ (Brush style).
  - **Body**: _Noto Sans JP_ (Clean, legible).

- **<img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/component.svg" width="16" height="16" /> Components**:
  - **Buttons/Cards**: "Panel" style with thick 3px black borders and drop shadows (`box-shadow: 4px 4px 0px`).
  - **Textures**: CSS-based "Screentone" overlays (radial gradients).

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/rocket.svg" width="24" height="24" /> Getting Started

### Prerequisites

- Node.js 18+
- pnpm (recommended)

### Installation

```bash
cd frontend
pnpm install
```

### Development

Run the development server:

```bash
pnpm dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/folder-tree.svg" width="24" height="24" /> Project Structure

```bash
frontend/
├── app/                  # Next.js App Router pages & layouts
│   ├── globals.css       # Global theme variables & Tailwind config
│   ├── layout.tsx        # Root layout (Fonts provider)
│   └── page.tsx          # Landing page (Upload Zone)
├── components/
│   ├── features/         # Complex domain components (e.g., UploadZone)
│   └── ui/               # Reusable primitives (Card, Button)
├── lib/                  # Utilities (cn, mock-data)
└── public/               # Static assets
```

## <img src="https://raw.githubusercontent.com/lucide-icons/lucide/main/icons/terminal.svg" width="24" height="24" /> Scripts

| Command      | Description                                 |
| ------------ | ------------------------------------------- |
| `pnpm dev`   | Start development server at localhost:3000  |
| `pnpm build` | Build the application for production        |
| `pnpm start` | Run the built production server             |
| `pnpm lint`  | Run ESLint to check for code quality issues |
