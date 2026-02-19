// API Client for Manga Translator Backend

// On the server (SSR/RSC), use API_URL (internal Docker service name).
// On the client (browser), use NEXT_PUBLIC_API_URL (public hostname).
const API_BASE_URL =
  (typeof window === "undefined"
    ? process.env.API_URL
    : process.env.NEXT_PUBLIC_API_URL) || "http://localhost:8080";

export interface Request {
  id: string;
  filename: string;
  fileType: "image" | "zip";
  status: "queued" | "processing" | "completed" | "failed";
  progress: number;
  pageCount: number;
  thumbnail?: string;
  errorMessage?: string;
  createdAt: string;
  updatedAt: string;
  completedAt?: string;
}

export interface Page {
  pageNumber: number;
  original: string;
  translated: string;
}

export interface Result {
  requestId: string;
  pages: Page[];
}

export interface ProgressUpdate {
  status: string;
  progress: number;
  message: string;
}

export interface ListRequestsResponse {
  requests: Request[];
  total: number;
  limit: number;
  offset: number;
}

// Upload files for translation
export async function uploadFiles(files: File[]): Promise<Request> {
  const formData = new FormData();

  files.forEach((file) => {
    formData.append("files", file);
  });

  const response = await fetch(`${API_BASE_URL}/api/translate`, {
    method: "POST",
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || "Upload failed");
  }

  return response.json();
}

// Get list of translation requests
export async function getRequests(
  status?: string,
  limit: number = 20,
  offset: number = 0,
): Promise<ListRequestsResponse> {
  const params = new URLSearchParams({
    limit: limit.toString(),
    offset: offset.toString(),
  });

  if (status) {
    params.append("status", status);
  }

  const response = await fetch(`${API_BASE_URL}/api/requests?${params}`);

  if (!response.ok) {
    throw new Error("Failed to fetch requests");
  }

  return response.json();
}

// Get a specific request by ID
export async function getRequest(id: string): Promise<Request> {
  const response = await fetch(`${API_BASE_URL}/api/requests/${id}`);

  if (!response.ok) {
    throw new Error("Failed to fetch request");
  }

  return response.json();
}

// Get translation results for a request
export async function getResults(requestId: string): Promise<Result> {
  const response = await fetch(`${API_BASE_URL}/api/results/${requestId}`);

  if (!response.ok) {
    throw new Error("Failed to fetch results");
  }

  return response.json();
}

// Subscribe to progress updates via Server-Sent Events
export function subscribeToProgress(
  requestId: string,
  onProgress: (update: ProgressUpdate) => void,
  onComplete: (update: ProgressUpdate) => void,
  onError: (error: string) => void,
): EventSource {
  const eventSource = new EventSource(
    `${API_BASE_URL}/api/requests/${requestId}/events`,
  );

  let hasCompleted = false;
  let hasFailed = false;

  eventSource.addEventListener("connected", (e) => {
    try {
      const data = JSON.parse(e.data) as ProgressUpdate;
      console.log("SSE connected:", data);
    } catch (err) {
      console.error("Failed to parse connected event:", err, e);
    }
  });

  eventSource.addEventListener("progress", (e) => {
    try {
      const data = JSON.parse(e.data) as ProgressUpdate;
      console.log("SSE progress:", data);
      onProgress(data);
    } catch (err) {
      console.error("Failed to parse progress event:", err, e);
    }
  });

  eventSource.addEventListener("complete", (e) => {
    try {
      const data = JSON.parse(e.data) as ProgressUpdate;
      console.log("SSE complete:", data);
      hasCompleted = true;
      onComplete(data);
      eventSource.close();
    } catch (err) {
      console.error("Failed to parse complete event:", err, e);
    }
  });

  // Handle custom SSE 'error' event from server (not the native onerror)
  eventSource.addEventListener("error", (e) => {
    const messageEvent = e as MessageEvent;
    if (messageEvent.data) {
      try {
        const data = JSON.parse(messageEvent.data) as ProgressUpdate;
        console.log("SSE error event:", data);
        hasFailed = true;
        onError(data.message);
        eventSource.close();
      } catch (parseErr) {
        console.error("Failed to parse SSE error event:", parseErr);
      }
    }
  });

  // Handle native EventSource errors (connection issues)
  eventSource.onerror = (err) => {
    console.log("SSE connection error or closed", err);
    // Only treat as error if we haven't completed successfully
    if (!hasCompleted && !hasFailed) {
      onError("Connection interrupted");
    }
    eventSource.close();
  };

  return eventSource;
}

// Build file URL
export function getFileUrl(
  requestId: string,
  type: "uploads" | "originals" | "translated",
  filename: string,
): string {
  return `${API_BASE_URL}/api/files/${requestId}/${type}/${filename}`;
}
