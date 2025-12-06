import React, { useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { motion, AnimatePresence } from 'framer-motion';
import { 
  FileText, 
  Upload, 
  Download, 
  Trash2, 
  Loader2, 
  Check, 
  AlertCircle,
  FileOutput,
  Copy
} from 'lucide-react';
import { Button, Card, CardHeader, CardBody } from '../ui';
import { usePdfExtractor } from '../../hooks/usePdfExtractor';
import { getDownloadUrl } from '../../services/api';
import { formatFileSize } from '../../utils/helpers';

export function PdfExtractor() {
  const {
    pdfFile,
    preview,
    result,
    loading,
    extracting,
    error,
    loadPdf,
    extract,
    clear,
    hasPdf,
  } = usePdfExtractor();

  const onDrop = useCallback((acceptedFiles) => {
    if (acceptedFiles.length > 0) {
      loadPdf(acceptedFiles[0]);
    }
  }, [loadPdf]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'application/pdf': ['.pdf'] },
    multiple: false,
  });

  const copyToClipboard = () => {
    if (result?.markdown) {
      navigator.clipboard.writeText(result.markdown);
    }
  };

  return (
    <div className="space-y-4">
      {/* Dropzone */}
      {!hasPdf && (
        <div
          {...getRootProps()}
          className={`
            border-2 border-dashed rounded-2xl p-8 text-center cursor-pointer transition-all duration-300
            ${isDragActive 
              ? 'border-indigo-500 bg-indigo-500/10' 
              : 'border-zinc-700 hover:border-zinc-500 bg-zinc-900/50'
            }
          `}
        >
          <input {...getInputProps()} />
          <div className={`w-16 h-16 mx-auto mb-4 rounded-2xl flex items-center justify-center transition-all ${isDragActive ? 'bg-indigo-500/20' : 'bg-zinc-800'}`}>
            <Upload className={`w-8 h-8 ${isDragActive ? 'text-indigo-400' : 'text-zinc-400'}`} />
          </div>
          <p className="text-zinc-300 mb-2">
            {isDragActive ? 'PDF dosyasını bırakın...' : 'PDF dosyasını sürükleyin'}
          </p>
          <p className="text-zinc-500 text-sm">veya tıklayarak seçin</p>
          <p className="text-zinc-600 text-xs mt-2">.pdf</p>
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <Card>
          <CardBody>
            <div className="flex items-center justify-center gap-3 py-8">
              <Loader2 className="w-6 h-6 text-indigo-400 animate-spin" />
              <span className="text-zinc-300">PDF analiz ediliyor...</span>
            </div>
          </CardBody>
        </Card>
      )}

      {/* Error State */}
      <AnimatePresence>
        {error && (
          <motion.div
            initial={{ opacity: 0, y: -10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="p-4 rounded-xl bg-red-500/10 border border-red-500/30 flex items-center gap-3"
          >
            <AlertCircle className="w-5 h-5 text-red-400 flex-shrink-0" />
            <span className="text-red-300">{error}</span>
          </motion.div>
        )}
      </AnimatePresence>

      {/* PDF Info & Preview */}
      {pdfFile && preview && !loading && (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="space-y-4"
        >
          {/* File Info */}
          <Card>
            <CardHeader
              action={
                <button
                  onClick={clear}
                  className="p-2 rounded-lg hover:bg-zinc-700 text-zinc-500 hover:text-red-400 transition-all"
                  title="Dosyayı kaldır"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              }
            >
              <FileText className="w-5 h-5 text-indigo-400" />
              <div>
                <p className="text-zinc-200 font-medium">{pdfFile.name}</p>
                <p className="text-zinc-500 text-xs">
                  {formatFileSize(pdfFile.size)} • {preview.pageCount} sayfa
                </p>
              </div>
            </CardHeader>
          </Card>

          {/* Preview */}
          <Card>
            <CardHeader>
              <FileOutput className="w-5 h-5 text-zinc-400" />
              <span className="text-zinc-200 font-medium">Markdown Önizleme</span>
            </CardHeader>
            <div className="p-4 max-h-80 overflow-y-auto">
              <pre className="text-zinc-300 text-sm whitespace-pre-wrap font-mono">
                {preview.preview}
              </pre>
            </div>
          </Card>

          {/* Extract Button */}
          {!result && (
            <Button
              onClick={extract}
              disabled={extracting}
              loading={extracting}
              icon={Download}
              className="w-full py-4"
            >
              {extracting ? 'Dönüştürülüyor...' : 'Markdown\'a Dönüştür'}
            </Button>
          )}
        </motion.div>
      )}

      {/* Result */}
      <AnimatePresence>
        {result && result.success && (
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            className="space-y-4"
          >
            <div className="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30">
              <div className="flex items-center gap-2 mb-3">
                <Check className="w-5 h-5 text-emerald-400" />
                <span className="text-emerald-300 font-medium">Dönüştürme başarılı!</span>
              </div>
              
              <div className="flex gap-2">
                <a
                  href={getDownloadUrl(result.filePath)}
                  download
                  className="flex-1 flex items-center justify-center gap-2 p-3 bg-zinc-900/50 rounded-lg hover:bg-zinc-800/50 transition-colors group"
                >
                  <Download className="w-4 h-4 text-zinc-400 group-hover:text-indigo-400" />
                  <span className="text-zinc-300 text-sm">{result.fileName}</span>
                </a>
                <button
                  onClick={copyToClipboard}
                  className="p-3 bg-zinc-900/50 rounded-lg hover:bg-zinc-800/50 transition-colors group"
                  title="Markdown'ı kopyala"
                >
                  <Copy className="w-4 h-4 text-zinc-400 group-hover:text-indigo-400" />
                </button>
              </div>
            </div>

            {/* Markdown Result */}
            <Card>
              <CardHeader>
                <FileText className="w-5 h-5 text-zinc-400" />
                <span className="text-zinc-200 font-medium">Çıkarılan Markdown</span>
              </CardHeader>
              <div className="p-4 max-h-96 overflow-y-auto">
                <pre className="text-zinc-300 text-sm whitespace-pre-wrap font-mono">
                  {result.markdown}
                </pre>
              </div>
            </Card>

            {/* Convert Another */}
            <Button
              onClick={clear}
              variant="secondary"
              icon={Upload}
              className="w-full"
            >
              Başka PDF Dönüştür
            </Button>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
