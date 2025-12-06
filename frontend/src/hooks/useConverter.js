import { useState, useCallback } from 'react';
import { convertToPdf } from '../services/api';

/**
 * Custom hook for PDF conversion
 * @returns {Object} Conversion state and utilities
 */
export function useConverter() {
  const [converting, setConverting] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState(null);

  const convert = useCallback(async ({ files, mergeMode, outputName }) => {
    setConverting(true);
    setError(null);
    setResult(null);

    try {
      const data = await convertToPdf({ files, mergeMode, outputName });
      setResult(data);
      return data;
    } catch (err) {
      const errorMessage = err.message || 'Conversion failed';
      setError(errorMessage);
      setResult({ success: false, error: errorMessage });
      throw err;
    } finally {
      setConverting(false);
    }
  }, []);

  const clearResult = useCallback(() => {
    setResult(null);
    setError(null);
  }, []);

  return {
    convert,
    converting,
    result,
    error,
    clearResult,
    isSuccess: result?.success === true,
    isError: result?.success === false || error !== null,
  };
}
