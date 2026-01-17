import { useState, useCallback } from 'react';
import { generateId, readFileAsText, isMarkdownFile } from '../utils/helpers';

/**
 * Custom hook for managing file state
 * @returns {Object} File management utilities
 */
export function useFiles() {
  const [files, setFiles] = useState([]);
  const [selectedFileId, setSelectedFileId] = useState(null);

  const selectedFile = files.find(f => f.id === selectedFileId) || null;

  const addFiles = useCallback(async (acceptedFiles) => {
    const mdFiles = acceptedFiles.filter(isMarkdownFile);
    
    const newFiles = await Promise.all(
      mdFiles.map(async (file, index) => {
        const content = await readFileAsText(file);
        return {
          id: generateId(),
          name: file.name,
          content,
          size: file.size,
          order: files.length + index,
        };
      })
    );

    setFiles(prev => {
      const updated = [...prev, ...newFiles];
      // Auto-select first file if none selected
      if (!selectedFileId && updated.length > 0) {
        setSelectedFileId(updated[0].id);
      }
      return updated;
    });

    return newFiles;
  }, [files.length, selectedFileId]);

  const removeFile = useCallback((id) => {
    setFiles(prev => {
      const updated = prev.filter(f => f.id !== id);
      // Update orders
      return updated.map((f, i) => ({ ...f, order: i }));
    });

    if (selectedFileId === id) {
      setSelectedFileId(files.find(f => f.id !== id)?.id || null);
    }
  }, [selectedFileId, files]);

  const clearFiles = useCallback(() => {
    setFiles([]);
    setSelectedFileId(null);
  }, []);

  const selectFile = useCallback((id) => {
    setSelectedFileId(id);
  }, []);

  const reorderFiles = useCallback((activeId, overId) => {
    setFiles(prev => {
      const oldIndex = prev.findIndex(f => f.id === activeId);
      const newIndex = prev.findIndex(f => f.id === overId);

      if (oldIndex === -1 || newIndex === -1) return prev;

      const updated = [...prev];
      const [removed] = updated.splice(oldIndex, 1);
      updated.splice(newIndex, 0, removed);

      // Update orders
      return updated.map((f, i) => ({ ...f, order: i }));
    });
  }, []);

  const getOrderedFiles = useCallback(() => {
    return [...files].sort((a, b) => a.order - b.order);
  }, [files]);

  return {
    files,
    selectedFile,
    selectedFileId,
    addFiles,
    removeFile,
    clearFiles,
    selectFile,
    reorderFiles,
    getOrderedFiles,
    hasFiles: files.length > 0,
    fileCount: files.length,
  };
}
