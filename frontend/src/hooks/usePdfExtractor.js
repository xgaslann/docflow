import { useState, useCallback } from 'react';
import { fetchPdfPreview, extractPdfToMarkdown } from '../services/api';

/**
 * Reads file as base64
 * @param {File} file
 * @returns {Promise<string>}
 */
function readFileAsBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      // Remove data URL prefix (data:application/pdf;base64,)
      const base64 = reader.result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = () => reject(new Error('Failed to read file'));
    reader.readAsDataURL(file);
  });
}

/**
 * Custom hook for PDF extraction
 * @returns {Object} PDF extraction utilities
 */
export function usePdfExtractor() {
  const [pdfFile, setPdfFile] = useState(null);
  const [preview, setPreview] = useState(null);
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [extracting, setExtracting] = useState(false);
  const [error, setError] = useState(null);

  const loadPdf = useCallback(async (file) => {
    if (!file || !file.name.toLowerCase().endsWith('.pdf')) {
      setError('Please select a PDF file');
      return;
    }

    setLoading(true);
    setError(null);
    setPreview(null);
    setResult(null);

    try {
      const base64 = await readFileAsBase64(file);
      setPdfFile({ file, base64, name: file.name, size: file.size });

      // Fetch preview
      const previewData = await fetchPdfPreview(base64, file.name);
      setPreview(previewData);
    } catch (err) {
      setError(err.message);
      setPdfFile(null);
    } finally {
      setLoading(false);
    }
  }, []);

  const extract = useCallback(async () => {
    if (!pdfFile) {
      setError('No PDF file loaded');
      return;
    }

    setExtracting(true);
    setError(null);
    setResult(null);

    try {
      const extractResult = await extractPdfToMarkdown(pdfFile.base64, pdfFile.name);
      setResult(extractResult);
    } catch (err) {
      setError(err.message);
    } finally {
      setExtracting(false);
    }
  }, [pdfFile]);

  const clear = useCallback(() => {
    setPdfFile(null);
    setPreview(null);
    setResult(null);
    setError(null);
  }, []);

  return {
    pdfFile,
    preview,
    result,
    loading,
    extracting,
    error,
    loadPdf,
    extract,
    clear,
    hasPdf: pdfFile !== null,
  };
}
