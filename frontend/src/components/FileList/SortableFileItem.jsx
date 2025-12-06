import React from 'react';
import { useSortable } from '@dnd-kit/sortable';
import { CSS } from '@dnd-kit/utilities';
import { FileText, X, GripVertical } from 'lucide-react';
import { cn, formatFileSize } from '../../utils/helpers';

export function SortableFileItem({ file, isSelected, onSelect, onRemove }) {
  const {
    attributes,
    listeners,
    setNodeRef,
    transform,
    transition,
    isDragging,
  } = useSortable({ id: file.id });

  const style = {
    transform: CSS.Transform.toString(transform),
    transition,
  };

  return (
    <div
      ref={setNodeRef}
      style={style}
      className={cn(
        'flex items-center gap-2 p-3 transition-all border-b border-zinc-800/50 last:border-b-0',
        isSelected
          ? 'bg-indigo-500/10 border-l-2 border-l-indigo-500'
          : 'hover:bg-zinc-800/50',
        isDragging && 'opacity-50 bg-zinc-800'
      )}
    >
      {/* Drag Handle */}
      <button
        {...attributes}
        {...listeners}
        className="p-1 rounded hover:bg-zinc-700 text-zinc-500 hover:text-zinc-300 cursor-grab active:cursor-grabbing"
        title="Sürükleyerek sırala"
      >
        <GripVertical className="w-4 h-4" />
      </button>

      {/* Order Badge */}
      <span className="w-6 h-6 flex items-center justify-center rounded bg-zinc-800 text-zinc-400 text-xs font-medium">
        {file.order + 1}
      </span>

      {/* File Info */}
      <div
        className="flex-1 min-w-0 flex items-center gap-3 cursor-pointer"
        onClick={() => onSelect(file.id)}
      >
        <FileText
          className={cn(
            'w-5 h-5 flex-shrink-0',
            isSelected ? 'text-indigo-400' : 'text-zinc-500'
          )}
        />
        <div className="flex-1 min-w-0">
          <p className="text-zinc-200 truncate text-sm">{file.name}</p>
          <p className="text-zinc-500 text-xs">{formatFileSize(file.size)}</p>
        </div>
      </div>

      {/* Remove Button */}
      <button
        onClick={(e) => {
          e.stopPropagation();
          onRemove(file.id);
        }}
        className="p-1.5 rounded-lg hover:bg-zinc-700 text-zinc-500 hover:text-red-400 transition-all"
        title="Dosyayı kaldır"
      >
        <X className="w-4 h-4" />
      </button>
    </div>
  );
}
