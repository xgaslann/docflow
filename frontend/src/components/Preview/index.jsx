import React from 'react';
import { motion } from 'framer-motion';
import { Eye, FileText } from 'lucide-react';
import { Card, CardHeader } from '../ui';
import { MarkdownPreview } from './MarkdownPreview';

export function Preview({ selectedFile, html, loading }) {
  return (
    <motion.div
      initial={{ opacity: 0, x: 20 }}
      animate={{ opacity: 1, x: 0 }}
      transition={{ delay: 0.2 }}
      className="h-full"
    >
      <Card className="h-[calc(100vh-180px)] flex flex-col">
        <CardHeader>
          <Eye className="w-5 h-5 text-zinc-400" />
          <span className="text-zinc-200 font-medium">Önizleme</span>
          {selectedFile && (
            <span className="text-zinc-500 text-sm ml-auto">
              {selectedFile.name}
            </span>
          )}
        </CardHeader>

        <div className="flex-1 p-6 overflow-y-auto">
          {selectedFile ? (
            <MarkdownPreview html={html} loading={loading} />
          ) : (
            <EmptyPreview />
          )}
        </div>
      </Card>
    </motion.div>
  );
}

function EmptyPreview() {
  return (
    <div className="flex flex-col items-center justify-center h-64 text-zinc-500">
      <FileText className="w-16 h-16 mb-4 opacity-30" />
      <p>Önizlemek için bir dosya seçin</p>
    </div>
  );
}
