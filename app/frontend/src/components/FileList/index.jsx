import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import {
  DndContext,
  closestCenter,
  KeyboardSensor,
  PointerSensor,
  useSensor,
  useSensors,
} from '@dnd-kit/core';
import {
  SortableContext,
  sortableKeyboardCoordinates,
  verticalListSortingStrategy,
} from '@dnd-kit/sortable';
import { Files } from 'lucide-react';
import { Card, CardHeader } from '../ui';
import { SortableFileItem } from './SortableFileItem';

export function FileList({
  files,
  selectedFileId,
  onSelect,
  onRemove,
  onReorder,
  onClearAll,
}) {
  const sensors = useSensors(
    useSensor(PointerSensor, {
      activationConstraint: {
        distance: 8,
      },
    }),
    useSensor(KeyboardSensor, {
      coordinateGetter: sortableKeyboardCoordinates,
    })
  );

  const handleDragEnd = (event) => {
    const { active, over } = event;

    if (over && active.id !== over.id) {
      onReorder(active.id, over.id);
    }
  };

  if (files.length === 0) {
    return null;
  }

  return (
    <AnimatePresence>
      <motion.div
        initial={{ opacity: 0, height: 0 }}
        animate={{ opacity: 1, height: 'auto' }}
        exit={{ opacity: 0, height: 0 }}
      >
        <Card>
          <CardHeader
            action={
              <button
                onClick={onClearAll}
                className="text-zinc-500 hover:text-red-400 text-sm transition-colors"
              >
                Tümünü sil
              </button>
            }
          >
            <Files className="w-4 h-4 text-zinc-400" />
            <span className="text-zinc-200 font-medium">{files.length} dosya</span>
            <span className="text-zinc-500 text-xs ml-2">(sürükleyerek sırala)</span>
          </CardHeader>

          <div className="max-h-72 overflow-y-auto">
            <DndContext
              sensors={sensors}
              collisionDetection={closestCenter}
              onDragEnd={handleDragEnd}
            >
              <SortableContext
                items={files.map((f) => f.id)}
                strategy={verticalListSortingStrategy}
              >
                {files.map((file) => (
                  <SortableFileItem
                    key={file.id}
                    file={file}
                    isSelected={selectedFileId === file.id}
                    onSelect={onSelect}
                    onRemove={onRemove}
                  />
                ))}
              </SortableContext>
            </DndContext>
          </div>
        </Card>
      </motion.div>
    </AnimatePresence>
  );
}
