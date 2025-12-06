import React from 'react';
import { cn } from '../../utils/helpers';

export function Card({ children, className = '', ...props }) {
  return (
    <div
      className={cn(
        'bg-zinc-900/70 backdrop-blur-sm rounded-2xl border border-zinc-800',
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export function CardHeader({ children, className = '', action }) {
  return (
    <div className={cn('p-4 border-b border-zinc-800 flex items-center justify-between', className)}>
      <div className="flex items-center gap-3">{children}</div>
      {action && <div>{action}</div>}
    </div>
  );
}

export function CardBody({ children, className = '' }) {
  return (
    <div className={cn('p-4', className)}>
      {children}
    </div>
  );
}
