export interface Request {
  id: string;
  filename: string;
  status: "queued" | "processing" | "completed" | "failed";
  progress: number;
  thumbnail?: string;
  pageCount: number;
  createdAt: string;
}

export interface Result {
  requestId: string;
  pages: {
    original: string;
    translated: string;
  }[];
}

export const mockRequests: Request[] = [
  {
    id: "req_1",
    filename: "One Piece - Chapter 1090.zip",
    status: "processing",
    progress: 45,
    pageCount: 18,
    createdAt: "2023-08-15T10:00:00Z",
  },
  {
    id: "req_2",
    filename: "Naruto - Chapter 450.zip",
    status: "completed",
    progress: 100,
    thumbnail: "https://placehold.co/200x300/e63946/white?text=Naruto",
    pageCount: 17,
    createdAt: "2023-08-14T15:30:00Z",
  },
  {
    id: "req_3",
    filename: "Bleach - Chapter 200.zip",
    status: "queued",
    progress: 0,
    pageCount: 19,
    createdAt: "2023-08-15T11:20:00Z",
  },
];

// Mock pages for the completed request
export const mockResults: Record<string, Result> = {
  req_2: {
    requestId: "req_2",
    pages: Array.from({ length: 5 }).map((_, i) => ({
      original: `https://placehold.co/600x900/111/white?text=Original+Page+${
        i + 1
      }`,
      translated: `https://placehold.co/600x900/fdfbf7/black?text=Translated+Page+${
        i + 1
      }`,
    })),
  },
};
