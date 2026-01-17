import React from 'react';
import { cn } from '../../utils/helpers';

export function TabSwitch({ tabs, activeTab, onChange }) {
  return (
    <div className="inline-flex p-1 rounded-xl bg-zinc-800/50 border border-zinc-700">
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => onChange(tab.id)}
          className={cn(
            'flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
            activeTab === tab.id
              ? 'bg-indigo-500 text-white shadow-lg shadow-indigo-500/25'
              : 'text-zinc-400 hover:text-zinc-200 hover:bg-zinc-700/50'
          )}
        >
          {tab.icon && <tab.icon className="w-4 h-4" />}
          {tab.label}
        </button>
      ))}
    </div>
  );
}
