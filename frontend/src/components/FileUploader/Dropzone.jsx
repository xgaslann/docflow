import React from 'react';
import { useDropzone } from 'react-dropzone';
import { Upload } from 'lucide-react';
import { cn } from '../../utils/helpers';

export function Dropzone({ onDrop, disabled = false }) {
  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    disabled,
    accept: {
      'text/markdown': ['.md', '.markdown'],
      'text/plain': ['.txt'],
    },
  });

  return (
    <div
      {...getRootProps()}
      className={cn(
        'border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all duration-300',
        isDragActive
          ? 'border-indigo-500 bg-indigo-500/10'
          : 'border-zinc-700 hover:border-zinc-500 bg-zinc-900/50',
        disabled && 'opacity-50 cursor-not-allowed'
      )}
    >
      <input {...getInputProps()} />
      
      <div
        className={cn(
          'w-16 h-16 mx-auto mb-4 rounded-2xl flex items-center justify-center transition-all',
          isDragActive ? 'bg-indigo-500/20' : 'bg-zinc-800'
        )}
      >
        <Upload
          className={cn(
            'w-8 h-8',
            isDragActive ? 'text-indigo-400' : 'text-zinc-400'
          )}
        />
      </div>

      <p className="text-zinc-300 mb-2">
        {isDragActive ? 'Dosyaları bırakın...' : 'Markdown dosyalarını sürükleyin'}
      </p>
      <p className="text-zinc-500 text-sm">veya tıklayarak seçin</p>
      <p className="text-zinc-600 text-xs mt-2">.md, .markdown, .txt</p>
    </div>
  );
}
