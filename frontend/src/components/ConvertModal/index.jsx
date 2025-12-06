import React, { useState, useEffect } from 'react';
import { Files, Layers, ArrowRight } from 'lucide-react';
import { Modal, ModalFooter, Button } from '../ui';
import { MergePreview } from './MergePreview';
import { useMergePreview } from '../../hooks/usePreview';

export function ConvertModal({
  isOpen,
  onClose,
  files,
  onConvert,
  converting,
}) {
  const [outputName, setOutputName] = useState('');
  const [showMergePreview, setShowMergePreview] = useState(false);
  const { html, loading, meta, loadMergePreview, clearPreview } = useMergePreview();

  // Load merge preview when showing it
  useEffect(() => {
    if (showMergePreview && files.length > 0) {
      loadMergePreview(files);
    }
  }, [showMergePreview, files, loadMergePreview]);

  // Reset state when modal closes
  useEffect(() => {
    if (!isOpen) {
      setOutputName('');
      setShowMergePreview(false);
      clearPreview();
    }
  }, [isOpen, clearPreview]);

  const handleSeparateConvert = () => {
    onConvert({ mergeMode: 'separate' });
  };

  const handleMergedConvert = () => {
    onConvert({ mergeMode: 'merged', outputName: outputName.trim() || undefined });
  };

  // Show file order in modal
  const orderedFiles = [...files].sort((a, b) => a.order - b.order);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="Birden Fazla Dosya Algılandı"
      size="lg"
    >
      <p className="text-zinc-400 mb-6">
        {files.length} dosya yüklediniz. Nasıl dönüştürmek istersiniz?
      </p>

      <div className="space-y-4">
        {/* Option 1: Separate */}
        <button
          onClick={handleSeparateConvert}
          disabled={converting}
          className="w-full p-4 rounded-xl border border-zinc-700 hover:border-indigo-500 hover:bg-indigo-500/5 transition-all text-left group disabled:opacity-50"
        >
          <div className="flex items-start gap-4">
            <div className="w-12 h-12 rounded-xl bg-zinc-800 group-hover:bg-indigo-500/20 flex items-center justify-center transition-colors flex-shrink-0">
              <Files className="w-6 h-6 text-zinc-400 group-hover:text-indigo-400" />
            </div>
            <div className="flex-1 min-w-0">
              <h3 className="font-medium text-zinc-200 group-hover:text-indigo-300">
                Ayrı Ayrı Dönüştür
              </h3>
              <p className="text-zinc-500 text-sm mt-1">
                Her dosya için ayrı PDF oluşturulur ({files.length} PDF)
              </p>
            </div>
          </div>
        </button>

        {/* Option 2: Merged */}
        <div className="p-4 rounded-xl border border-zinc-700 hover:border-indigo-500 transition-all">
          <div className="flex items-start gap-4">
            <div className="w-12 h-12 rounded-xl bg-zinc-800 flex items-center justify-center flex-shrink-0">
              <Layers className="w-6 h-6 text-zinc-400" />
            </div>
            <div className="flex-1 min-w-0">
              <h3 className="font-medium text-zinc-200">Tek Dosyada Birleştir</h3>
              <p className="text-zinc-500 text-sm mt-1 mb-4">
                Tüm dosyalar aşağıdaki sırayla tek PDF'de birleştirilir
              </p>

              {/* File Order Display */}
              <div className="bg-zinc-800/50 rounded-lg p-3 mb-4">
                <p className="text-zinc-400 text-xs mb-2 font-medium">
                  Birleştirme sırası:
                </p>
                <div className="flex flex-wrap gap-2">
                  {orderedFiles.map((file, index) => (
                    <div
                      key={file.id}
                      className="flex items-center gap-1.5 bg-zinc-900 px-2 py-1 rounded text-xs"
                    >
                      <span className="text-indigo-400 font-medium">
                        {index + 1}.
                      </span>
                      <span className="text-zinc-300 truncate max-w-32">
                        {file.name}
                      </span>
                      {index < orderedFiles.length - 1 && (
                        <ArrowRight className="w-3 h-3 text-zinc-600 ml-1" />
                      )}
                    </div>
                  ))}
                </div>
                <p className="text-zinc-500 text-xs mt-2">
                  💡 Sıralamayı değiştirmek için dosya listesinde sürükle-bırak yapın
                </p>
              </div>

              {/* Toggle Merge Preview */}
              <button
                type="button"
                onClick={() => setShowMergePreview(!showMergePreview)}
                className="text-indigo-400 hover:text-indigo-300 text-sm mb-4 flex items-center gap-1"
              >
                {showMergePreview ? 'Önizlemeyi gizle' : 'Birleştirilmiş önizlemeyi göster'}
              </button>

              {/* Merge Preview - Fixed height with internal scroll */}
              {showMergePreview && (
                <div className="mb-4">
                  <MergePreview html={html} loading={loading} meta={meta} />
                </div>
              )}

              {/* Output Name Input */}
              <input
                type="text"
                value={outputName}
                onChange={(e) => setOutputName(e.target.value)}
                placeholder="Dosya adı (opsiyonel)"
                className="w-full px-3 py-2 bg-zinc-800 border border-zinc-700 rounded-lg text-zinc-200 text-sm placeholder:text-zinc-500 focus:outline-none focus:border-indigo-500 mb-3"
              />

              <Button
                onClick={handleMergedConvert}
                disabled={converting}
                loading={converting}
                className="w-full"
              >
                Birleştir ve Dönüştür
              </Button>
            </div>
          </div>
        </div>
      </div>

      <ModalFooter>
        <Button variant="ghost" onClick={onClose} disabled={converting}>
          İptal
        </Button>
      </ModalFooter>
    </Modal>
  );
}
