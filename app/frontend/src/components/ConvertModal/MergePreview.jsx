import React from 'react';
import { Loader2, FileText, BookOpen } from 'lucide-react';

export function MergePreview({ html, loading, meta }) {
  return (
    <div className="border border-zinc-700 rounded-xl overflow-hidden">
      {/* Preview Header */}
      <div className="bg-zinc-800 px-4 py-3 flex items-center justify-between flex-shrink-0">
        <div className="flex items-center gap-2">
          <BookOpen className="w-4 h-4 text-indigo-400" />
          <span className="text-zinc-200 text-sm font-medium">
            Birleştirilmiş Önizleme
          </span>
        </div>
        <div className="flex items-center gap-4 text-xs text-zinc-400">
          <span className="flex items-center gap-1">
            <FileText className="w-3 h-3" />
            {meta.totalFiles} dosya
          </span>
          <span>~{meta.estimatedPages} sayfa</span>
        </div>
      </div>

      {/* Preview Content - Fixed height with scroll */}
      <div className="h-48 overflow-y-auto bg-zinc-950 p-4">
        {loading ? (
          <div className="flex items-center justify-center py-8">
            <Loader2 className="w-6 h-6 text-indigo-400 animate-spin" />
          </div>
        ) : html ? (
          <div
            className="markdown-preview text-sm"
            dangerouslySetInnerHTML={{ __html: html }}
          />
        ) : (
          <p className="text-zinc-500 text-center py-4">
            Önizleme yükleniyor...
          </p>
        )}
      </div>
    </div>
  );
}
