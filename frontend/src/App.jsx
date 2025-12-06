import React, { useState } from 'react';
import { motion } from 'framer-motion';
import { Download, FileText, FileOutput } from 'lucide-react';

// Components
import { Header } from './components/Header';
import { FileUploader } from './components/FileUploader';
import { FileList } from './components/FileList';
import { Preview } from './components/Preview';
import { ConvertModal } from './components/ConvertModal';
import { ConvertResult } from './components/ConvertResult';
import { PdfExtractor } from './components/PdfExtractor';
import { Button, TabSwitch } from './components/ui';

// Hooks
import { useFiles } from './hooks/useFiles';
import { usePreview } from './hooks/usePreview';
import { useConverter } from './hooks/useConverter';
import { useTheme } from './hooks/useTheme';

const TABS = [
  { id: 'md-to-pdf', label: 'MD → PDF', icon: FileOutput },
  { id: 'pdf-to-md', label: 'PDF → MD', icon: FileText },
];

function App() {
  const [activeTab, setActiveTab] = useState('md-to-pdf');
  const [showModal, setShowModal] = useState(false);

  // Theme management
  const { isDark, toggleTheme } = useTheme();

  // File management (for MD to PDF)
  const {
    files,
    selectedFile,
    selectedFileId,
    addFiles,
    removeFile,
    clearFiles,
    selectFile,
    reorderFiles,
    getOrderedFiles,
    hasFiles,
    fileCount,
  } = useFiles();

  // Preview
  const { html, loading: previewLoading } = usePreview(selectedFile);

  // Conversion
  const { convert, converting, result, clearResult } = useConverter();

  // Handle file drop
  const handleFilesAdded = async (acceptedFiles) => {
    await addFiles(acceptedFiles);
  };

  // Handle convert button click
  const handleConvertClick = () => {
    if (fileCount > 1) {
      setShowModal(true);
    } else {
      handleConvert({ mergeMode: 'separate' });
    }
  };

  // Handle conversion
  const handleConvert = async ({ mergeMode, outputName }) => {
    const orderedFiles = getOrderedFiles();
    try {
      await convert({
        files: orderedFiles,
        mergeMode,
        outputName,
      });
      setShowModal(false);
    } catch (error) {
      console.error('Conversion error:', error);
    }
  };

    return (
        <div className="min-h-screen p-6 md:p-8 bg-main transition-colors duration-300">
            <div className="container-centered">
                <Header isDark={isDark} onThemeToggle={toggleTheme} />

                {/* Tab Switcher */}
                <div className="mb-6">
                    <TabSwitch
                        tabs={TABS}
                        activeTab={activeTab}
                        onChange={setActiveTab}
                    />
                </div>

                {/* MD to PDF Tab */}
                {activeTab === 'md-to-pdf' && (
                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                        {/* Left Panel - File Management */}
                        <motion.div
                            initial={{ opacity: 0, x: -20 }}
                            animate={{ opacity: 1, x: 0 }}
                            transition={{ delay: 0.1 }}
                            className="space-y-4"
                        >
                            {/* File Uploader */}
                            <FileUploader onFilesAdded={handleFilesAdded} disabled={converting} />

                            {/* File List with Drag & Drop Sorting */}
                            <FileList
                                files={files}
                                selectedFileId={selectedFileId}
                                onSelect={selectFile}
                                onRemove={removeFile}
                                onReorder={reorderFiles}
                                onClearAll={clearFiles}
                            />

                            {/* Convert Button */}
                            {hasFiles && (
                                <motion.div
                                    initial={{ opacity: 0, y: 20 }}
                                    animate={{ opacity: 1, y: 0 }}
                                >
                                    <Button
                                        onClick={handleConvertClick}
                                        disabled={converting}
                                        loading={converting}
                                        icon={Download}
                                        className="w-full py-4"
                                    >
                                        {converting ? 'Dönüştürülüyor...' : 'PDF\'e Dönüştür'}
                                    </Button>
                                </motion.div>
                            )}

                            {/* Conversion Result */}
                            <ConvertResult result={result} onClear={clearResult} />
                        </motion.div>

                        {/* Right Panel - Preview */}
                        <div>
                            <Preview
                                selectedFile={selectedFile}
                                html={html}
                                loading={previewLoading}
                            />
                        </div>
                    </div>
                )}

                {/* PDF to MD Tab */}
                {activeTab === 'pdf-to-md' && (
                    <motion.div
                        initial={{ opacity: 0, y: 20 }}
                        animate={{ opacity: 1, y: 0 }}
                        className="max-w-2xl mx-auto"
                    >
                        <PdfExtractor />
                    </motion.div>
                )}
            </div>

            {/* Convert Modal for Multiple Files */}
            <ConvertModal
                isOpen={showModal}
                onClose={() => setShowModal(false)}
                files={getOrderedFiles()}
                onConvert={handleConvert}
                converting={converting}
            />
        </div>
    );
}

export default App;
