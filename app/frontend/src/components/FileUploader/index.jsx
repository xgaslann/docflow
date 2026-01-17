import React from 'react';
import { Dropzone } from './Dropzone';

export function FileUploader({ onFilesAdded, disabled = false }) {
  const handleDrop = async (acceptedFiles) => {
    if (onFilesAdded) {
      await onFilesAdded(acceptedFiles);
    }
  };

  return <Dropzone onDrop={handleDrop} disabled={disabled} />;
}
