// API Client for Manga Translator Backend

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

export interface Request {
  id: string;
  filename: string;
  fileType: 'image' | 'zip';
  status: 'queued' | 'processing' | 'completed' | 'failed';
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
    formData.append('files', file);
  });

  const response = await fetch(`${API_BASE_URL}/api/translate`, {
    method: 'POST',
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Upload failed');
  }

  return response.json();
}

// Get list of translation requests
export async function getRequests(
  status?: string,
  limit: number = 20,
  offset: number = 0
): Promise<ListRequestsResponse> {
  const params = new URLSearchParams({
    limit: limit.toString(),
    offset: offset.toString(),
  });

  if (status) {
    params.append('status', status);
  }

  const response = await fetch(`${API_BASE_URL}/api/requests?${params}`);

  if (!response.ok) {
    throw new Error('Failed to fetch requests');
  }

  return response.json();
}

// Get a specific request by ID
export async function getRequest(id: string): Promise<Request> {
  const response = await fetch(`${API_BASE_URL}/api/requests/${id}`);

  if (!response.ok) {
    throw new Error('Failed to fetch request');
  }

  return response.json();
}

// Get translation results for a request
export async function getResults(requestId: string): Promise<Result> {
  const response = await fetch(`${API_BASE_URL}/api/results/${requestId}`);

  if (!response.ok) {
    throw new Error('Failed to fetch results');
  }

  return response.json();
}

// Subscribe to progress updates via Server-Sent Events
export function subscribeToProgress(
  requestId: string,
  onProgress: (update: ProgressUpdate) => void,
  onComplete: (update: ProgressUpdate) => void,
  onError: (error: string) => void
): EventSource {
  const eventSource = new EventSource(
    `${API_BASE_URL}/api/requests/${requestId}/events`
  );

  eventSource.addEventListener('connected', (e) => {
    const data = JSON.parse(e.data) as ProgressUpdate;
    console.log('SSE connected:', data);
  });

  eventSource.addEventListener('progress', (e) => {
    const data = JSON.parse(e.data) as ProgressUpdate;
    onProgress(data);
  });

  eventSource.addEventListener('complete', (e) => {
    const data = JSON.parse(e.data) as ProgressUpdate;
    onComplete(data);
    eventSource.close();
  });

  eventSource.addEventListener('error', (e: Event) => {
    const messageEvent = e as MessageEvent;
    if (messageEvent.data) {
      const data = JSON.parse(messageEvent.data) as ProgressUpdate;
      onError(data.message);
    } else {
      onError('Connection error');
    }
    eventSource.close();
  });

  eventSource.onerror = () => {
    // Connection error or closed by server
    eventSource.close();
  };

  return eventSource;
}

// Build file URL
export function getFileUrl(requestId: string, type: 'uploads' | 'originals' | 'translated', filename: string): string {
  return `${API_BASE_URL}/api/files/${requestId}/${type}/${filename}`;
}
