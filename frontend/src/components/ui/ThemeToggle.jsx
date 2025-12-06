import React from 'react';
import { Sun, Moon } from 'lucide-react';
import { cn } from '../../utils/helpers';

export function ThemeToggle({ isDark, onToggle }) {
  return (
    <button
      onClick={onToggle}
      className={cn(
        'relative w-14 h-7 rounded-full p-1 transition-colors duration-300',
        isDark ? 'bg-zinc-700' : 'bg-indigo-100'
      )}
      title={isDark ? 'Light mode' : 'Dark mode'}
      aria-label={isDark ? 'Switch to light mode' : 'Switch to dark mode'}
    >
      {/* Track icons */}
      <Sun className={cn(
        'absolute left-1.5 top-1/2 -translate-y-1/2 w-4 h-4 transition-opacity',
        isDark ? 'opacity-30 text-zinc-500' : 'opacity-100 text-amber-500'
      )} />
      <Moon className={cn(
        'absolute right-1.5 top-1/2 -translate-y-1/2 w-4 h-4 transition-opacity',
        isDark ? 'opacity-100 text-indigo-300' : 'opacity-30 text-zinc-400'
      )} />
      
      {/* Sliding knob */}
      <div
        className={cn(
          'w-5 h-5 rounded-full shadow-md transition-transform duration-300 flex items-center justify-center',
          isDark 
            ? 'translate-x-7 bg-zinc-900' 
            : 'translate-x-0 bg-white'
        )}
      >
        {isDark ? (
          <Moon className="w-3 h-3 text-indigo-300" />
        ) : (
          <Sun className="w-3 h-3 text-amber-500" />
        )}
      </div>
    </button>
  );
}
