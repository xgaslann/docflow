import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, act } from '@testing-library/react';
import { useTheme } from '../useTheme';

describe('useTheme', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.clear();
    document.documentElement.classList.remove('light', 'dark');
  });

  it('defaults to light theme', () => {
    const { result } = renderHook(() => useTheme());

    expect(result.current.theme).toBe('light');
    expect(result.current.isLight).toBe(true);
    expect(result.current.isDark).toBe(false);
  });

  it('toggles theme', () => {
    const { result } = renderHook(() => useTheme());

    expect(result.current.theme).toBe('light');

    act(() => {
      result.current.toggleTheme();
    });

    expect(result.current.theme).toBe('dark');
    expect(result.current.isDark).toBe(true);
    expect(result.current.isLight).toBe(false);

    act(() => {
      result.current.toggleTheme();
    });

    expect(result.current.theme).toBe('light');
  });

  it('sets light theme explicitly', () => {
    const { result } = renderHook(() => useTheme());

    act(() => {
      result.current.setDarkTheme();
    });

    expect(result.current.isDark).toBe(true);

    act(() => {
      result.current.setLightTheme();
    });

    expect(result.current.isLight).toBe(true);
  });

  it('sets dark theme explicitly', () => {
    const { result } = renderHook(() => useTheme());

    act(() => {
      result.current.setDarkTheme();
    });

    expect(result.current.theme).toBe('dark');
    expect(result.current.isDark).toBe(true);
  });

  it('persists theme to localStorage', () => {
    const { result } = renderHook(() => useTheme());

    act(() => {
      result.current.setDarkTheme();
    });

    expect(localStorage.setItem).toHaveBeenCalledWith('md-to-pdf-theme', 'dark');
  });

  it('applies class to document element', () => {
    const { result } = renderHook(() => useTheme());

    // Initial light theme
    expect(document.documentElement.classList.contains('light')).toBe(true);

    act(() => {
      result.current.setDarkTheme();
    });

    expect(document.documentElement.classList.contains('dark')).toBe(true);
    expect(document.documentElement.classList.contains('light')).toBe(false);
  });

  it('reads saved theme from localStorage', () => {
    localStorage.getItem.mockReturnValue('dark');

    const { result } = renderHook(() => useTheme());

    expect(result.current.theme).toBe('dark');
  });

  it('ignores invalid localStorage value', () => {
    localStorage.getItem.mockReturnValue('invalid');

    const { result } = renderHook(() => useTheme());

    expect(result.current.theme).toBe('light');
  });
});
