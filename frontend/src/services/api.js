const API_BASE = '';

/**
 * Fetches markdown preview HTML
 * @param {string} content - Markdown content
 * @returns {Promise<{html: string}>}
 */
export async function fetchPreview(content) {
  const response = await fetch(`${API_BASE}/api/preview`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ content }),
  });

  if (!response.ok) {
    throw new Error('Preview fetch failed');
  }

  return response.json();
}

/**
 * Fetches merged files preview
 * @param {Array<{id: string, name: string, content: string, order: number}>} files
 * @returns {Promise<{html: string, totalFiles: number, estimatedPages: number}>}
 */
export async function fetchMergePreview(files) {
  const response = await fetch(`${API_BASE}/api/preview/merge`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ files }),
  });

  if (!response.ok) {
    throw new Error('Merge preview fetch failed');
  }

  return response.json();
}

/**
 * Converts files to PDF
 * @param {Object} params
 * @param {Array<{id: string, name: string, content: string, order: number}>} params.files
 * @param {'separate' | 'merged'} params.mergeMode
 * @param {string} [params.outputName]
 * @returns {Promise<{success: boolean, files?: string[], error?: string}>}
 */
export async function convertToPdf({ files, mergeMode, outputName }) {
  const response = await fetch(`${API_BASE}/api/convert`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      files,
      mergeMode,
      outputName,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Conversion failed');
  }

  return response.json();
}

/**
 * Previews PDF extraction
 * @param {string} base64Content - Base64 encoded PDF
 * @param {string} fileName - Original file name
 * @returns {Promise<{preview: string, pageCount: number, fileName: string}>}
 */
export async function fetchPdfPreview(base64Content, fileName) {
  const response = await fetch(`${API_BASE}/api/pdf/preview`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      content: base64Content,
      fileName,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'PDF preview failed');
  }

  return response.json();
}

/**
 * Extracts PDF to Markdown
 * @param {string} base64Content - Base64 encoded PDF
 * @param {string} fileName - Original file name
 * @returns {Promise<{success: boolean, markdown: string, filePath?: string, fileName?: string, error?: string}>}
 */
export async function extractPdfToMarkdown(base64Content, fileName) {
  const response = await fetch(`${API_BASE}/api/pdf/extract`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      content: base64Content,
      fileName,
    }),
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'PDF extraction failed');
  }

  return response.json();
}

/**
 * Health check
 * @returns {Promise<{status: string, version: string, timestamp: number}>}
 */
export async function healthCheck() {
  const response = await fetch(`${API_BASE}/api/health`);
  return response.json();
}

/**
 * Gets the download URL for a file
 * @param {string} path - File path from API
 * @returns {string}
 */
export function getDownloadUrl(path) {
  return `${API_BASE}${path}`;
}
