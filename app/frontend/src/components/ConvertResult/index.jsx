import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { Check, AlertCircle, FileText, Download, X } from 'lucide-react';
import { getDownloadUrl } from '../../services/api';

export function ConvertResult({ result, onClear }) {
  if (!result) return null;

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0, scale: 0.95 }}
        animate={{ opacity: 1, scale: 1 }}
        exit={{ opacity: 0, scale: 0.95 }}
        className={`p-4 rounded-xl relative ${
          result.success
            ? 'bg-emerald-500/10 border border-emerald-500/30'
            : 'bg-red-500/10 border border-red-500/30'
        }`}
      >
        {/* Close button */}
        <button
          onClick={onClear}
          className="absolute top-2 right-2 p-1 rounded hover:bg-white/10 text-zinc-400 hover:text-zinc-200 transition-colors"
        >
          <X className="w-4 h-4" />
        </button>

        {result.success ? (
          <SuccessContent files={result.files} />
        ) : (
          <ErrorContent error={result.error} />
        )}
      </motion.div>
    </AnimatePresence>
  );
}

function SuccessContent({ files }) {
  return (
    <div>
      <div className="flex items-center gap-2 mb-3">
        <Check className="w-5 h-5 text-emerald-400" />
        <span className="text-emerald-300 font-medium">
          Dönüştürme başarılı!
        </span>
      </div>
      <div className="space-y-2">
        {files?.map((file, i) => (
          <a
            key={i}
            href={getDownloadUrl(file)}
            download
            className="flex items-center gap-2 p-3 bg-zinc-900/50 rounded-lg hover:bg-zinc-800/50 transition-colors group"
          >
            <FileText className="w-4 h-4 text-zinc-400 group-hover:text-indigo-400" />
            <span className="text-zinc-300 text-sm flex-1 truncate">
              {file.split('/').pop()}
            </span>
            <Download className="w-4 h-4 text-zinc-500 group-hover:text-indigo-400" />
          </a>
        ))}
      </div>
    </div>
  );
}

function ErrorContent({ error }) {
  return (
    <div className="flex items-center gap-2">
      <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0" />
      <span className="text-red-300">{error}</span>
    </div>
  );
}
