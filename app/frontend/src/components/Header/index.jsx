import React from 'react';
import { motion } from 'framer-motion';
import { FileOutput } from 'lucide-react';
import { ThemeToggle } from '../ui/ThemeToggle';

export function Header({ isDark, onThemeToggle }) {
  return (
    <motion.header
      initial={{ opacity: 0, y: -20 }}
      animate={{ opacity: 1, y: 0 }}
      className="max-w-7xl mx-auto mb-8"
    >
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-indigo-600 flex items-center justify-center shadow-lg shadow-indigo-500/25">
            <FileOutput className="w-6 h-6 text-white" />
          </div>
          <div>
            <h1 className="text-2xl font-semibold text-heading">MD → PDF Converter</h1>
            <p className="text-muted text-sm">
              Markdown dosyalarınızı profesyonel PDF'lere dönüştürün
            </p>
          </div>
        </div>
        
        <ThemeToggle isDark={isDark} onToggle={onThemeToggle} />
      </div>
    </motion.header>
  );
}
