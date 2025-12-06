import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import {
  fetchPreview,
  fetchMergePreview,
  convertToPdf,
  fetchPdfPreview,
  extractPdfToMarkdown,
  healthCheck,
  getDownloadUrl,
} from '../api';

describe('API Service', () => {
  const originalFetch = global.fetch;

  beforeEach(() => {
    global.fetch = vi.fn();
  });

  afterEach(() => {
    global.fetch = originalFetch;
  });

  describe('fetchPreview', () => {
    it('sends POST request with content', async () => {
      const mockResponse = { html: '<h1>Hello</h1>' };
      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await fetchPreview('# Hello');

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/preview',
        expect.objectContaining({
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ content: '# Hello' }),
        })
      );
      expect(result).toEqual(mockResponse);
    });

    it('throws on error response', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: false,
        status: 400,
      });

      await expect(fetchPreview('')).rejects.toThrow('Preview fetch failed');
    });
  });

  describe('fetchMergePreview', () => {
    it('sends POST request with files', async () => {
      const files = [
        { id: '1', name: 'a.md', content: '# A', order: 0 },
        { id: '2', name: 'b.md', content: '# B', order: 1 },
      ];
      const mockResponse = {
        html: '<h1>A</h1><h1>B</h1>',
        totalFiles: 2,
        estimatedPages: 1,
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await fetchMergePreview(files);

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/preview/merge',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({ files }),
        })
      );
      expect(result.totalFiles).toBe(2);
    });
  });

  describe('convertToPdf', () => {
    it('sends conversion request', async () => {
      const params = {
        files: [{ id: '1', name: 'test.md', content: '# Test', order: 0 }],
        mergeMode: 'separate',
      };
      const mockResponse = {
        success: true,
        files: ['/output/test.pdf'],
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await convertToPdf(params);

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/convert',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify(params),
        })
      );
      expect(result.success).toBe(true);
    });

    it('throws with error message on failure', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: false,
        json: () => Promise.resolve({ error: 'Conversion failed' }),
      });

      await expect(
        convertToPdf({ files: [], mergeMode: 'separate' })
      ).rejects.toThrow('Conversion failed');
    });
  });

  describe('fetchPdfPreview', () => {
    it('sends base64 content and filename', async () => {
      const mockResponse = {
        preview: '# Extracted',
        pageCount: 3,
        fileName: 'doc.pdf',
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await fetchPdfPreview('base64content', 'doc.pdf');

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/pdf/preview',
        expect.objectContaining({
          method: 'POST',
          body: JSON.stringify({
            content: 'base64content',
            fileName: 'doc.pdf',
          }),
        })
      );
      expect(result.pageCount).toBe(3);
    });
  });

  describe('extractPdfToMarkdown', () => {
    it('sends extraction request', async () => {
      const mockResponse = {
        success: true,
        markdown: '# Extracted Content',
        filePath: '/output/doc.md',
        fileName: 'doc.md',
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await extractPdfToMarkdown('base64content', 'doc.pdf');

      expect(global.fetch).toHaveBeenCalledWith(
        '/api/pdf/extract',
        expect.objectContaining({
          method: 'POST',
        })
      );
      expect(result.success).toBe(true);
      expect(result.markdown).toBe('# Extracted Content');
    });
  });

  describe('healthCheck', () => {
    it('calls health endpoint', async () => {
      const mockResponse = {
        status: 'healthy',
        version: '1.0.0',
        timestamp: 1234567890,
      };

      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      });

      const result = await healthCheck();

      expect(global.fetch).toHaveBeenCalledWith('/api/health');
      expect(result.status).toBe('healthy');
    });
  });

  describe('getDownloadUrl', () => {
    it('returns correct URL', () => {
      const url = getDownloadUrl('/output/test.pdf');
      expect(url).toBe('/output/test.pdf');
    });

    it('handles paths without leading slash', () => {
      const url = getDownloadUrl('output/test.pdf');
      expect(url).toBe('output/test.pdf');
    });
  });
});
