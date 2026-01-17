import { useState, useEffect, useCallback } from 'react';
import { fetchPreview, fetchMergePreview } from '../services/api';
import { debounce } from '../utils/helpers';

/**
 * Custom hook for markdown preview
 * @param {Object|null} file - Selected file object
 * @returns {Object} Preview state and utilities
 */
export function usePreview(file) {
  const [html, setHtml] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const loadPreview = useCallback(
    debounce(async (content) => {
      if (!content) {
        setHtml('');
        return;
      }

      setLoading(true);
      setError(null);

      try {
        const data = await fetchPreview(content);
        setHtml(data.html || '');
      } catch (err) {
        setError(err.message);
        setHtml('<p class="error">Preview yüklenemedi</p>');
      } finally {
        setLoading(false);
      }
    }, 300),
    []
  );

  useEffect(() => {
    if (file?.content) {
      loadPreview(file.content);
    } else {
      setHtml('');
    }
  }, [file?.content, loadPreview]);

  return { html, loading, error };
}

/**
 * Custom hook for merge preview
 * @returns {Object} Merge preview state and utilities
 */
export function useMergePreview() {
  const [html, setHtml] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [meta, setMeta] = useState({ totalFiles: 0, estimatedPages: 0 });

  const loadMergePreview = useCallback(async (files) => {
    if (!files || files.length === 0) {
      setHtml('');
      setMeta({ totalFiles: 0, estimatedPages: 0 });
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const data = await fetchMergePreview(files);
      setHtml(data.html || '');
      setMeta({
        totalFiles: data.totalFiles || files.length,
        estimatedPages: data.estimatedPages || 1,
      });
    } catch (err) {
      setError(err.message);
      setHtml('<p class="error">Merge preview yüklenemedi</p>');
    } finally {
      setLoading(false);
    }
  }, []);

  const clearPreview = useCallback(() => {
    setHtml('');
    setMeta({ totalFiles: 0, estimatedPages: 0 });
    setError(null);
  }, []);

  return { html, loading, error, meta, loadMergePreview, clearPreview };
}
